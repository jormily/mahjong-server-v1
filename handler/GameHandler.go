package handler

import (
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"mahjong/tool"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/jormily/mahjong-server/db"
	"github.com/jormily/mahjong-server/manager"
	"github.com/jormily/mahjong-server/msg"
)

var (
	DI_FEN = []int{1,2,5}
	MAX_FAN = []int{3,4,5}
	JU_SHU = []int{4,8}
	JU_SHU_COST = []int{2,3}

	Error_Check_Conf = errors.New("check conf error")
)

func checkConf(conf string) (*msg.RoomConf,error) {
	roomConf := &msg.RoomConf{}
	if err := json.Unmarshal([]byte(conf),roomConf);err != nil {
		return nil,err
	}

	if roomConf.BaseScore < 0 || roomConf.BaseScore > len(DI_FEN) {
		return nil,Error_Check_Conf
	}

	if roomConf.ZiMo < 0 || roomConf.ZiMo > 2 {
		return nil,Error_Check_Conf
	}

	if roomConf.MaxFan < 0 || roomConf.MaxFan > len(MAX_FAN) {
		return nil,Error_Check_Conf
	}

	if roomConf.MaxGames < 0 || roomConf.MaxGames > len(JU_SHU) {
		return nil,Error_Check_Conf
	}

	roomConf.BaseScore = DI_FEN[roomConf.BaseScore]
	roomConf.MaxGames = JU_SHU[roomConf.MaxGames]
	roomConf.MaxFan = MAX_FAN[roomConf.MaxFan]
	roomConf.DianGanH,_ = strconv.ParseInt(roomConf.DianGangHua,10,64)
	return roomConf,nil
}

func generateRoomId() string {
	getRoomId := func() string {
		roomId := 0
		for i := 0;i < 6;i++ {
			rand.Seed(time.Now().UnixNano())
			roomId = roomId*10 + (rand.Intn(9) + 1)
		}
		return strconv.FormatInt(int64(roomId),10)
	}

	var roomid string
	for {
		roomid = getRoomId()
		if db.RoomExist(roomid) {
			continue
		}else{
			break
		}
	}
	return roomid
}

//func GetServerInfoGs(w http.ResponseWriter, r *http.Request) {
//	query := r.URL.Query()
//	serverid := query.Get("serverid")
//	sign := query.Get("sign")
//	if serverid != viper.GetString("game_svr.id") || sign == "" {
//		Response(w,msg.BaseResponse{1,"invalid parameters"})
//		return
//	}
//
//	md5 := MD5(serverid+viper.GetString("core.room_pri_key"))
//	if md5 != sign {
//		Response(w,msg.BaseResponse{1,"sign check failed"})
//		return
//	}
//
//	location := manager.RoomMgr.GetUserLocation(userid)
//}

func CreateRoom(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	userid,err1 := strconv.ParseInt(query.Get("userid"), 10, 64)
	gems,err2 := strconv.ParseInt(query.Get("gems"), 10, 64)
	sign := query.Get("sign")
	conf := query.Get("conf")
	if err1 != nil || sign == "null" || err2 != nil || conf == "null" {
		tool.Response(w,msg.BaseResponse{1,"invalid parameters"})
		return
	}

	md5 := tool.MD5(strconv.FormatInt(userid,10)+ conf+ strconv.FormatInt(gems,10) + viper.GetString("core.room_pri_key"))
	if md5 != sign {
		tool.Response(w,msg.BaseResponse{1,"sign check failed."})
		return
	}

	roomConf,err2 := checkConf(conf)
	if err2 != nil {
		tool.Response(w,msg.BaseResponse{1,err2.Error()})
		return
	}

	cost := JU_SHU_COST[roomConf.MaxGames]
	if cost > int(gems) {
		tool.Response(w,msg.BaseResponse{1,"gems is not enough"})
		return
	}

	roomId := generateRoomId()
	config,_ := json.Marshal(roomConf)
	if _,b := db.InsertRoom(roomId,string(config),viper.GetString("game_svr.ip"),viper.GetInt("game_svr.port"),time.Now().Unix());b {
		tool.Response(w,msg.BaseResponse{1,"create failed"})
		return
	}else{
		if err := manager.RoomMgr.CreateRoom(roomId,userid);err != nil {
			tool.Response(w,msg.BaseResponse{1,err.Error()})
			return
		}
		tool.Response(w,msg.CreateRoomResponse{msg.BaseResponse{0,"ok"},roomId})
		return
	}
}

func EnterRoom(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	userid,_ := strconv.ParseInt(query.Get("userid"),10,64)
	name := query.Get("name")
	roomid := query.Get("roomid")
	sign := query.Get("sign")
	if name == "null" || name == "" ||
		roomid == "null" || roomid == "" ||
		sign == "null" || sign == "" {
		tool.Response(w,msg.BaseResponse{1,"invalid parameters"})
		return
	}

	md5 := tool.MD5(strconv.FormatInt(userid,10)+name+roomid+viper.GetString("core.room_pri_key"))
	if md5 != sign {
		tool.Response(w,msg.BaseResponse{2,"sign check error"})
		return
	}

	if err := manager.RoomMgr.EnterRoom(roomid,userid,name,1000);err != nil {
		tool.Response(w,msg.BaseResponse{3,err.Error()})
		return
	}

	tool.Response(w,msg.EnterRoomResponse{
		BaseResponse:msg.BaseResponse{0,"ok"},
		Token: manager.TokenMgr.CreateToken(userid,5000),
	})
	return
}

func Ping(w http.ResponseWriter, r *http.Request) {

}

func IsRoomRuning(w http.ResponseWriter, r *http.Request) {

}


func update(){
	//log.Info("update")
	//log.Info(viper.GetString("game_svr.id"))
	values := url.Values{
		"id":{viper.GetString("game_svr.id")},
		"clientip":{viper.GetString("game_svr.ip")},
		"clientport":{viper.GetString("game_svr.port")},
		"httpPort":{viper.GetString("game_svr.http_port")},
		"load":{"0"},
	}
	_,err := tool.HttpGet(viper.GetString("hall_svr.host"),"/register_gs",values)
	if err != nil {
		log.Info("update err:"+err.Error())
		return
	}
}

func Ticker(sec int){
	tiker := time.NewTicker(5*time.Second)
	for {
		select {
			case<-tiker.C:
				update()
		}
	}
}

func GameServe(host string,port string){
	go Ticker(5)
	//http.HandleFunc("/get_server_info", GetServerInfo)
	http.HandleFunc("/create_room", CreateRoom)
	http.HandleFunc("/enter_room", EnterRoom)
	http.HandleFunc("/is_room_runing", IsRoomRuning)

	http.ListenAndServe(host+":"+port,nil)
}