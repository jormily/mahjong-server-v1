package handler

import (
	"github.com/jormily/mahjong-server/tool"
	"github.com/spf13/viper"
	"net/http"
	"strconv"

	"github.com/jormily/mahjong-server/db"
	"github.com/jormily/mahjong-server/msg"
	log "github.com/sirupsen/logrus"
)


func GetServerInfo(w http.ResponseWriter, r *http.Request) {
	log.Info("GetServerInfo")

	data := msg.ServerInfo{
		Version: viper.GetInt("core.version"),
		Hall: viper.GetString("hall_svr.host"),
		Appweb: viper.GetString("act_svr.app_web"),
	}
	tool.Response(w,data)
}

func GetVersion(w http.ResponseWriter, r *http.Request) {
	log.Info("GetVersion")
	data := msg.ServerVersion{
		Version: viper.GetInt("core.version"),
	}
	tool.Response(w,data)
}

func Guest(w http.ResponseWriter, r *http.Request) {
	log.Info("Guest")
	query := r.URL.Query()
	accout := "guest_" + query.Get("account")
	data := msg.GuestResponse{
		BaseResponse:msg.BaseResponse{Errcode: 0, Errmsg: "ok"},
		Account: accout,
		Halladdr: "",
		Sign: tool.MD5(accout + r.Host + viper.GetString("core.accout_pri_key")),
	}
	tool.Response(w,data)
}

func GetBaseInfo(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	userId,_ := strconv.ParseInt(query.Get("userid"),10,64)
	userModel,err := db.GetUserById(userId)
	if err != nil {
		tool.Response(w,msg.BaseResponse{1,"not found user"})
		return
	}
	tool.Response(w,msg.BaseInfoResponse{
		BaseResponse:msg.BaseResponse{0,"ok"},
		Name:userModel.Name,
		Sex: int(userModel.Sex),
		HeadImgUrl: userModel.Headimg,
	})



}

func Auth(w http.ResponseWriter, r *http.Request) {

}

func WechatAuth(w http.ResponseWriter, r *http.Request) {

}

func AccountServe(host string){
	http.HandleFunc("/base_info",GetBaseInfo)
	http.HandleFunc("/get_version", GetVersion)
	http.HandleFunc("/get_serverinfo", GetServerInfo)
	http.HandleFunc("/guest", Guest)

	http.ListenAndServe(host,nil)
}

