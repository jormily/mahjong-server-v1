package service

import (
	"github.com/jormily/mahjong-server/manager"
	"github.com/jormily/mahjong-server/msg"
)


func OnLogin(room *manager.Room,id int,data interface{}){
	var playerInfo msg.PlayerInfo
	result := msg.LoginResultResponse{}
	result.Errcode = 0
	result.Errmsg = "ok"
	result.Result.Seats = make([]msg.PlayerInfo,0)
	for _,v := range room.GetAllPlayer() {
		if v.UserId != 0 {
			playerSeat := msg.PlayerInfo{
				UserId:     v.UserId,
				Ip:         "" ,
				Score:      v.Score,
				Name:       v.Name,
				Online:     v.State > 0,
				Ready:      v.Ready,
				SeateIndex: v.Seate,
			}
			if player := manager.PlayerMgr.GetPlayerByUserId(v.UserId);player != nil {
				playerSeat.Ip = player.Ip()
			}
			result.Result.Seats = append(result.Result.Seats, playerSeat)

			if room.GetUserId(id) == v.UserId {
				playerInfo = playerSeat
			}
		}
	}
	room.Broadcast("login_result",result)
	room.BroadcastEx("new_user_comes_push",playerInfo,playerInfo.SeateIndex)
	room.Send(id,"login_finished",nil)

	room.GetLogic().OnReady(id)
}

func OnReady(room *manager.Room,id int,data interface{}){
	//playerData := room.GetPlayerData(room.GetUserId(id))
	//playerData.Ready = true
	//room.BroadcastEx("user_ready_push",msg.ReadyResponse{
	//	UserId: room.GetUserId(id),
	//	Ready: true,
	//},id);
	room.GetLogic().OnReady(id)
}

func OnMakeLack(room *manager.Room,id int,data interface{}){
	room.GetLogic().OnMakeLack(id,data.(int))
}

func OnPlayCard(room *manager.Room,id int,data interface{}){
	room.GetLogic().OnPlayCard(id,data.(int))
}

func OnPeng(room *manager.Room,id int,data interface{}){
	room.GetLogic().OnPeng(id)
}

func OnGang(room *manager.Room,id int,data interface{}){
	room.GetLogic().OnGang(id,data.(int))
}

func OnPass(room *manager.Room,id int,data interface{}){
	room.GetLogic().OnPass(id)
}

func OnHu(room *manager.Room,id int,data interface{}){
	room.GetLogic().OnHu(id)
}

func OnGamePing(room *manager.Room,id int,data interface{}){
	room.Send(id,"game_pong",nil)
}
