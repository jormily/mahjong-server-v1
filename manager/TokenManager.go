package manager

import (
	"github.com/jormily/mahjong-server/tool"
	"strconv"
	"time"
)

var (
	TokenMgr = NewTokenManager()
)

type Token struct {
	UserId 		int64
	Time 		int64
	LifeTime	int64
}

type TokenManager struct {
	users 	map[int64]string
	tokens 	map[string]*Token
}

func NewTokenManager() *TokenManager {
	this := new(TokenManager)
	this.users = map[int64]string{}
	this.tokens = map[string]*Token{}
	return this
}

func (this *TokenManager) CreateToken(userId int64,lifeTime int64) string {
	if t,ok := this.users[userId];ok {
		this.DeleteToken(t)
	}

	now := time.Now().Unix()
	tokenStr := tool.MD5(strconv.FormatInt(userId,10)+"!@#$%^&"+strconv.FormatInt(now,10))
	token := &Token{
		UserId: userId,
		Time: now,
		LifeTime: lifeTime,
	}
	this.tokens[tokenStr] = token
	this.users[userId] = tokenStr

	return tokenStr
}

func (this *TokenManager) GetToken(userId int64) (string,bool) {
	if t,ok := this.users[userId];ok {
		return t,true
	}else{
		return "",false
	}
}

func (this *TokenManager) GetUserId(token string) int64 {
	if t,ok := this.tokens[token];ok {
		return t.UserId
	}else{
		return 0
	}
}

func (this *TokenManager) IsTokenValid(token string) bool {
	if t,ok := this.tokens[token];ok {
		if t == nil {
			return false
		}

		if t.Time + t.LifeTime < time.Now().Unix() {
			return false
		}
		return true
	}
	return false
}

func (this *TokenManager) DeleteToken(token string){
	if t,ok := this.tokens[token];ok {
		delete(this.tokens,token)
		delete(this.users,t.UserId)
	}
}