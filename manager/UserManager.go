package manager

import (
	"errors"
	sio "github.com/graarh/golang-socketio"
)

var (
	UserMgr = NewUserManager()
)

type UserManager struct {
	u2c			map[int64]*sio.Channel
	c2u			map[*sio.Channel]int64
	userOnline	int
}

func NewUserManager() *UserManager {
	this := new(UserManager)
	this.u2c = map[int64]*sio.Channel{}
	this.c2u = map[*sio.Channel]int64{}
	this.userOnline = 0
	return this
}

func (this *UserManager) GetChannelUser(c *sio.Channel) (int64,error) {
	if u,ok := this.c2u[c];ok {
		return u,nil
	}
	return 0,errors.New("channel not bind user")
}

func (this *UserManager) GetUserChannel(userid int64) (*sio.Channel,error) {
	if c,ok := this.u2c[userid];ok {
		return c,nil
	}
	return nil,errors.New("channel not bind user")
}

func (this *UserManager) Bind(userid int64,c *sio.Channel){
	this.u2c[userid] = c
	this.c2u[c] = userid
	this.userOnline ++
}

func (this *UserManager) Delete(userid int64){
	if c,ok := this.u2c[userid];ok {
		delete(this.u2c,userid)
		delete(this.c2u,c)
	}
}

func (this *UserManager) IsOnline(userid int64) bool {
	if _,ok := this.u2c[userid];ok {
		return true
	}
	return false
}

func (this *UserManager) GetOnlineCount() int {
	return this.userOnline
}

func (this *UserManager) SendMessage(userid int64,event string,msg interface{}){
	if c,ok := this.u2c[userid];ok {
		c.Emit(event,msg)
	}
}

func (this *UserManager) KickAllInRoom(roomid string){

}

func (this *UserManager) BroacastInRoom(event string,msg interface{},user int64,include bool){

}