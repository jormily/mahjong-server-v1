package service

import (
	"github.com/jormily/mahjong-server/manager"
	"github.com/jormily/mahjong-server/msg"
	"github.com/jormily/mahjong-server/tool"
	"github.com/spf13/viper"
	"strconv"

	sio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	log "github.com/sirupsen/logrus"
	"net/http"
)



func socketServRun()  {

	//server.On("huanpai", func(c *sio.Channel, msg BaseResponse) string {
	//	return "OK"
	//})
	//




	//server.On("chat", func(c *sio.Channel, msg BaseResponse) string {
	//	return "OK"
	//})
	//
	//server.On("quick_chat", func(c *sio.Channel, msg BaseResponse) string {
	//	return "OK"
	//})
	//
	//server.On("voice_msg", func(c *sio.Channel, msg BaseResponse) string {
	//	return "OK"
	//})
	//
	//server.On("emoji", func(c *sio.Channel, msg BaseResponse) string {
	//	return "OK"
	//})
	//
	//server.On("exit", func(c *sio.Channel, msg BaseResponse) string {
	//	return "OK"
	//})
	//
	//server.On("dispress", func(c *sio.Channel, msg BaseResponse) string {
	//	return "OK"
	//})
	//
	//server.On("dissolve_request", func(c *sio.Channel, msg BaseResponse) string {
	//	return "OK"
	//})
	//
	//server.On("dissolve_agree", func(c *sio.Channel, msg BaseResponse) string {
	//	return "OK"
	//})
	//
	//server.On("dissolve_reject", func(c *sio.Channel, msg BaseResponse) string {
	//	return "OK"
	//})
	//


	server := sio.NewServer(transport.GetDefaultWebsocketTransport())

	//handle connected
	server.On(sio.OnConnection, func(c *sio.Channel) {
		log.Info("client connected")
	})

	server.On(sio.OnDisconnection, func(c *sio.Channel) {
		log.Info("client disconnected")
	})

	//handle custom event
	server.On("login", func(c *sio.Channel, data msg.LoginMsg) string {
		log.Info("login")
		if player :=  manager.PlayerMgr.GetPlayerByChannel(c);player != nil {
			return "OK"
		}

		if data.Token == "null" || data.RoomId == "null" || data.Token == "" || data.RoomId == ""{
			c.Emit("login_result",msg.BaseResponse{1,"invalid parameters"})
			return "OK"
		}
		md5 := tool.MD5(data.RoomId + data.Token + strconv.FormatInt(data.Time,10) + viper.GetString("core.room_pri_key"))
		if md5 != data.Sign {
			c.Emit("login_result",msg.BaseResponse{2,"check sign error"})
			return "OK"
		}

		if !manager.TokenMgr.IsTokenValid(data.Token) {
			c.Emit("login_result",msg.BaseResponse{3,"token out of time"})
			return "OK";
		}
		userId := manager.TokenMgr.GetUserId(data.Token)
		room := manager.RoomMgr.GetRoomByUserId(userId)
		playerData := room.GetPlayerData(userId)
		if playerData == nil {
			c.Emit("login_result",msg.BaseResponse{3,"token out of time"})
			return "OK"
		}

		playerData.Ready = false
		playerData.State = 1
		manager.PlayerMgr.AddPlayer(c,room,playerData.UserId,playerData.Name,0)

		OnLogin(room,room.GetPlayerId(userId),nil)

		return "OK"
	})

	server.On("ready", func(c *sio.Channel) string {
		player := manager.PlayerMgr.GetPlayerByChannel(c)
		if player == nil {
			return "OK"
		}
		room := player.GetRoom()
		if room == nil {
			return "OK"
		}
		OnReady(room,room.GetPlayerId(player.GetUserId()),nil)
		return "OK"
	})

	server.On("dingque", func(c *sio.Channel, lack int) string {
		player := manager.PlayerMgr.GetPlayerByChannel(c)
		if player == nil {
			return "OK"
		}
		room := player.GetRoom()
		if room == nil {
			return "OK"
		}
		OnMakeLack(room,room.GetPlayerId(player.GetUserId()),lack)
		return "OK"
	})


	server.On("chupai", func(c *sio.Channel, card int) string {
		player := manager.PlayerMgr.GetPlayerByChannel(c)
		if player == nil {
			return "OK"
		}
		room := player.GetRoom()
		if room == nil {
			return "OK"
		}
		OnPlayCard(room,room.GetPlayerId(player.GetUserId()),card)
		return "OK"
	})


	server.On("peng", func(c *sio.Channel, msg string) string {
		player := manager.PlayerMgr.GetPlayerByChannel(c)
		if player == nil {
			return "OK"
		}
		room := player.GetRoom()
		if room == nil {
			return "OK"
		}
		OnPeng(room,room.GetPlayerId(player.GetUserId()),nil)
		return "OK"
	})


	server.On("gang", func(c *sio.Channel, card int) string {
		player := manager.PlayerMgr.GetPlayerByChannel(c)
		if player == nil {
			return "OK"
		}
		room := player.GetRoom()
		if room == nil {
			return "OK"
		}
		OnGang(room,room.GetPlayerId(player.GetUserId()),card)
		return "OK"
	})

	server.On("hu", func(c *sio.Channel, msg string) string {
		player := manager.PlayerMgr.GetPlayerByChannel(c)
		if player == nil {
			return "OK"
		}
		room := player.GetRoom()
		if room == nil {
			return "OK"
		}
		OnHu(room,room.GetPlayerId(player.GetUserId()),nil)
		return "OK"
	})

	server.On("guo", func(c *sio.Channel, msg string) string {
		player := manager.PlayerMgr.GetPlayerByChannel(c)
		if player == nil {
			return "OK"
		}
		room := player.GetRoom()
		if room == nil {
			return "OK"
		}
		OnPass(room,room.GetPlayerId(player.GetUserId()),nil)
		return "OK"
	})


	server.On("game_ping", func(c *sio.Channel) string {
		player := manager.PlayerMgr.GetPlayerByChannel(c)
		if player == nil {
			return "OK"
		}
		room := player.GetRoom()
		if room == nil {
			return "OK"
		}
		OnGamePing(room,room.GetPlayerId(player.GetUserId()),nil)
		return "OK"
	})

	http.HandleFunc("/socket.io/",func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		server.ServeHTTP(w, r)
	})
	http.ListenAndServe("127.0.0.1:1000", nil)
}

func SocketServStart() {
	manager.RoomMgr.CreateAllRoom()
	go socketServRun()
}
