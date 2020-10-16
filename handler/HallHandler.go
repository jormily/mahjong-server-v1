package handler

import (
	"encoding/json"
	"errors"
	"github.com/jormily/mahjong-server/db"
	"github.com/jormily/mahjong-server/msg"
	"github.com/jormily/mahjong-server/tool"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type GameServerInfo struct {
	Ip 			string
	Id 			string
	ClientIp 	string
	ClientPort 	string
	HttpPort 	string
	Load 		string
}

var (
	GameServerMap = map[string]*GameServerInfo{}
)

func chooseServer() *GameServerInfo {
	var gs *GameServerInfo = nil
	for _,info := range GameServerMap {
		if gs == nil {
			gs = info
		}else {
			if info.Load < gs.Load {
				gs = info
			}
		}
	}
	return gs
}

func checkAccount(w http.ResponseWriter, r *http.Request) bool {
	var query = r.URL.Query()
	var account = query.Get("account")
	var sign = query.Get("sign")
	if account == "null" || sign == "null" || account == "" || sign == "" {
		tool.Response(w,msg.BaseResponse{1,"unknown error"})
		return false
	}
	return true
}

func enterRoom(userid int64,name string,roomid string) (*msg.EnterRoomResponse,error) {
	var roomModel,err1 = db.GetRoomModel(roomid)
	if err1 != nil {
		return nil,err1
	}

	var serverHost = roomModel.Ip + ":" + strconv.FormatInt(roomModel.Port,10)
	var serverInfo = GameServerMap[serverHost]
	if serverInfo == nil {
		return nil,errors.New("not found server")
	}

	var values = url.Values{
		"userid":{strconv.FormatInt(userid,10)},
		"name":{name},
		"roomid":{roomid},
		"sign":{tool.MD5(strconv.FormatInt(userid,10)+name+roomid+viper.GetString("core.room_pri_key"))},
	}

	var body,err2 = tool.HttpGet(serverInfo.Ip+":"+serverInfo.HttpPort,"/enter_room",values)
	if err2 != nil {
		return nil,err2
	}
	var res = msg.EnterRoomResponse{}
	if err3 := json.Unmarshal(body,&res);err3 != nil {
		return nil,err3
	}

	res.Ip = serverInfo.Ip
	res.Port = serverInfo.ClientPort
	db.SetUserRoomId(userid,roomid)

	return &res,nil
}

func createRoom(account string,userid int64,roomConf string) (*msg.CreateRoomResponse,error){
	log.Info("createRoom")
	var serverInfo = chooseServer()
	if serverInfo == nil {
		return nil,errors.New("chooseServer error")
	}

	var userModel,err1 = db.GetUser(account)
	if err1 != nil {
		return nil,err1
	}

	if userModel.Gems <= 0 {
		return nil,errors.New("has no gems")
	}

	var values = url.Values{
		"userid":{strconv.FormatInt(userid,10)},
		"gems":{strconv.FormatInt(userModel.Gems,10)},
		"conf":{roomConf},
		"sign":{tool.MD5(strconv.FormatInt(userid,10)+roomConf+strconv.FormatInt(userModel.Gems,10)+viper.GetString("core.room_pri_key"))},
	}
	var body,err2 = tool.HttpGet(serverInfo.Ip+":"+serverInfo.HttpPort,"/create_room",values)
	if err2 != nil {
		return nil,err2
	}

	log.Info(string(body))
	var res = msg.CreateRoomResponse{}
	if err3 := json.Unmarshal(body,&res);err3 != nil {
		return nil,err3
	}

	return &res,nil
}

func isServerOnline(ip string,port string) bool {
	var id = ip + ":" + port
	var gs,ok = GameServerMap[id]
	if !ok {
		return false
	}

	var values = url.Values{}
	values.Set("sign",tool.MD5(viper.GetString("core.room_pri_key")))
	var _,err = tool.HttpGet(gs.Ip+":"+gs.HttpPort,"/ping",values)
	if err == nil {
		return false
	}
	return true
}

func Login(w http.ResponseWriter, r *http.Request) {
	log.Info("Login")
	if !checkAccount(w,r) {
		return
	}

	var query = r.URL.Query()
	var account = query.Get("account")
	var model,err = db.GetUser(account)
	if err != nil {
		tool.Response(w,msg.BaseResponse{0,"ok"})
		return
	}

	var res = msg.LoginResponse{
		BaseResponse:msg.BaseResponse{
			Errcode:0,
			Errmsg: "ok",
		},
		Account: model.Account,
		UserId: model.Userid,
		Name: model.Name,
		Level: model.Lv,
		Exp: model.Exp,
		Coins: model.Coins,
		Gems: model.Gems,
		Ip: r.RemoteAddr,
		Sex: model.Sex,
	}

	if model.Roomid != "" {
		if db.RoomExist(model.Roomid) {
			res.RoomId = model.Roomid
			tool.Response(w,res)
			return
		}else {
			db.SetUserRoomId(model.Userid,"")
		}
	}
	tool.Response(w,res)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	log.Info("CreateUser")
	if !checkAccount(w,r) {
		return
	}

	var query = r.URL.Query()
	var account = query.Get("account")
	var name = query.Get("name")
	if db.UserExist(account) {
		tool.Response(w,msg.BaseResponse{1,"account have already exist."})
		return
	}

	err := db.CreateUser(account,name,viper.GetInt64("user.coins"),viper.GetInt64("user.gems"),0,"")
	if err != nil {
		tool.Response(w,msg.BaseResponse{2,"system error."})
		return
	}

	tool.Response(w,msg.BaseResponse{0,"ok"})
}


func GetUserStatus(w http.ResponseWriter,r *http.Request){
	log.Info("GetUserStatus")
	if !checkAccount(w,r) {
		return
	}

	var query = r.URL.Query()
	var account = query.Get("account")
	var model,err = db.GetUser(account)
	if err != nil {
		tool.Response(w,msg.BaseResponse{1,"get gams failed"})
	}else{
		tool.Response(w,msg.UserStateResponse{msg.BaseResponse{0,"ok",},int(model.Gems)})
	}
}

func CreatePrivateRoom(w http.ResponseWriter,r *http.Request)  {
	log.Info("CreatePrivateRoom")
	if !checkAccount(w,r) {
		return
	}
	var query = r.URL.Query()
	var account = query.Get("account")
	var conf = query.Get("conf")
	var userModel,err1 = db.GetUser(account)
	if err1 != nil {
		tool.Response(w,msg.BaseResponse{1,"system error"})
		return
	}

	if userModel.Roomid != "" {
		tool.Response(w,msg.BaseResponse{-1,"user is playing in room now."})
		return
	}

	//创建房间
	var res,err2 = createRoom(account,userModel.Userid,conf)
	if err2 != nil || res.Errcode != 0 || res.RoomId == "null" || res.RoomId == "" {
		tool.Response(w,msg.BaseResponse{-1,"create room error."})
		return
	}

	tool.Response(w,msg.CreateRoomResponse{msg.BaseResponse{0,"ok"},res.RoomId})
}

func RequestEnterRoom(w http.ResponseWriter, r *http.Request){
	log.Info("RequestEnterRoom")
	var query = r.URL.Query()
	var roomId = query.Get("roomid")
	if roomId == "null" || roomId == ""  {
		tool.Response(w,msg.BaseResponse{1,"parameters don't match api requirements."})
		return
	}
	if !checkAccount(w,r) {
		return
	}
	var account = query.Get("account")
	var user,err1 = db.GetUser(account)
	if err1 != nil {
		tool.Response(w,msg.BaseResponse{1,"system error"})
		return
	}

	if user.Roomid != "" && user.Roomid != roomId {
		tool.Response(w,msg.BaseResponse{1,"system error"})
		return
	}

	var res,err2 = enterRoom(user.Userid,user.Name,roomId)
	if err2 != nil || res.Errcode != 0 {
		tool.Response(w,msg.BaseResponse{1,"enter room failed."})
		return
	}

	now := time.Now().Unix()
	tool.Response(w,msg.EnterRoomResponse{
		BaseResponse:msg.BaseResponse{Errcode:0,Errmsg:"ok"},
		RoomId: roomId,
		Ip:res.Ip,
		Port: res.Port,
		Token: res.Token,
		Time: now,
		Sign: tool.MD5(roomId + res.Token + strconv.FormatInt(now,10) + viper.GetString("core.room_pri_key")),
	})
}



func GetMessage(w http.ResponseWriter, r *http.Request){
	log.Info("GetMessage")
	if !checkAccount(w,r) {
		return
	}

	var query = r.URL.Query()
	var typ = query.Get("type")
	if typ == "null" {
		tool.Response(w,msg.BaseResponse{-1,"parameters don't match api requirements"})
		return
	}

	var version = query.Get("version")
	var model,err = db.GetMessage(typ,version)
	if err != nil {
		tool.Response(w,msg.BaseResponse{1,"get message failed"})
	}else{
		tool.Response(w,msg.MessageResponse{msg.BaseResponse{0,"ok",},model.Msg,version})
	}
}

func RegisterGameServer(w http.ResponseWriter, r *http.Request){
	log.Info("RegisterGameServer")
	var query = r.URL.Query()
	var ip = query.Get("clientip")
	var clientip = query.Get("clientip")
	var clientport = query.Get("clientport")
	var httpPort = query.Get("httpPort")
	var load = query.Get("load")
	var id = clientip + ":" + clientport

	var gs,ok = GameServerMap[id]
	if ok {
		if gs.ClientIp == clientip && gs.ClientPort == clientport && gs.Ip == ip {
			gs.Load = load
			tool.Response(w,msg.GameServerResponse{msg.BaseResponse{0,"ok"},ip})
		}else {
			tool.Response(w,msg.BaseResponse{1,"duplicate gsid"+id})
		}
		return
	}
	GameServerMap[id] = &GameServerInfo{
		Ip: ip,
		Id: id,
		ClientPort: clientport,
		ClientIp: clientip,
		HttpPort: httpPort,
		Load: load,
	}
	tool.Response(w,msg.GameServerResponse{msg.BaseResponse{0,"ok"},ip})
}

func HallServe(host string){
	http.HandleFunc("/login", Login)
	http.HandleFunc("/create_user", CreateUser)
	http.HandleFunc("/get_user_status", GetUserStatus)
	http.HandleFunc("/get_message", GetMessage)
	http.HandleFunc("/enter_private_room", RequestEnterRoom)
	http.HandleFunc("/register_gs", RegisterGameServer)
	http.HandleFunc("/create_private_room",CreatePrivateRoom)

	http.ListenAndServe(host,nil)

}