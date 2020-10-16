package manager

import (
	sio "github.com/graarh/golang-socketio"
	log "github.com/sirupsen/logrus"
)

var (
	PlayerMgr = NewPlayerManager()
)
/**
 玩家数据
 */
type Player struct {
	*sio.Channel
	room 		*Room
	// 基础数据
	userid 		int64
	name 		string
	score 		int
}

func NewPlayer(c *sio.Channel,room *Room,userId int64,name string,score int) *Player {
	this := new(Player)
	this.Channel = c
	this.room = room
	this.userid = userId
	this.score = score
	return this
}

func (this *Player)Send(event string,msg interface{}){
	this.Emit(event,msg)
}

func (this *Player)GetUserId() int64 {
	return this.userid
}

func (this *Player)GetRoom() *Room {
	return this.room
}


type PlayerManager struct {
	uPlayers 	map[int64]*Player
	cPlayers 	map[*sio.Channel]*Player
	playerOnline	int
}

func NewPlayerManager() *PlayerManager {
	this := new(PlayerManager)
	this.uPlayers = map[int64]*Player{}
	this.cPlayers = map[*sio.Channel]*Player{}
	this.playerOnline = 0
	return this
}

func (this *PlayerManager) GetPlayerByChannel(c *sio.Channel) *Player {
	if player,ok := this.cPlayers[c];ok {
		return player
	}
	return nil
}

func (this *PlayerManager)GetPlayerByUserId(userid int64) *Player {
	if player,ok := this.uPlayers[userid];ok {
		return player
	}
	return nil
}



func (this *PlayerManager)AddPlayer(c *sio.Channel,room *Room,userId int64,name string,score int) *Player {
	if p,ok := this.cPlayers[c];ok {
		return p
	}

	player := NewPlayer(c,room,userId,name,score)
	this.cPlayers[c] = player
	this.uPlayers[userId] = player
	this.playerOnline++
	return player
}

func (this *PlayerManager)DeletePlayer(userid  int64){
	if player,ok := this.uPlayers[userid];ok {
		delete(this.cPlayers,player.Channel)
		delete(this.uPlayers,player.userid)
		this.playerOnline--
	}
}

func (this *PlayerManager)IsOnline(userid int64) bool {
	if _,ok := this.uPlayers[userid];ok {
		return true
	}
	return false
}

func (this *PlayerManager)GetOnlineCount() int {
	return this.playerOnline
}

func (this *PlayerManager) Send(userid int64,event string,msg interface{}){
	if p,ok := this.uPlayers[userid];ok {
		if event != "game_pong" {
			log.Infof("SendMessage-%d:{%s:%v}", userid, event, msg)
		}
		p.Emit(event,msg)
	}
}

func (this *PlayerManager) Broadcast(event string,msg interface{}){
	for _,v := range this.uPlayers {
		v.Emit(event,msg)
	}
}
//func (this *UserManager) Bind(userid int64,c *sio.Channel){
//	this.u2c[userid] = c
//	this.c2u[c] = userid
//	this.userOnline ++
//}

//func (this *UserManager) Delete(userid int64){
//	if c,ok := this.u2c[userid];ok {
//		delete(this.u2c,userid)
//		delete(this.c2u,c)
//	}
//}
//
//func (this *UserManager) IsOnline(userid int64) bool {
//	if _,ok := this.u2c[userid];ok {
//		return true
//	}
//	return false
//}
//
//func (this *UserManager) GetOnlineCount() int {
//	return this.userOnline
//}
//
//func (this *UserManager) SendMessage(userid int64,event string,msg interface{}){
//	if c,ok := this.u2c[userid];ok {
//		c.Emit(event,msg)
//	}
//}
//
//func (this *UserManager) KickAllInRoom(roomid string){
//
//}
//
//func (this *UserManager) BroacastInRoom(event string,msg interface{},user int64,include bool){
//
//}