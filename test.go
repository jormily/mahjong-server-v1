package main

import (
	"encoding/json"
	"fmt"
	socketio "github.com/googollee/go-socket.io"
	sio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	log "github.com/sirupsen/logrus"
	"mahjong/msg"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"
)

type BaseMsg struct {
	Errcode 	int		`json:"errcode"`
	Errmsg 		string	`json:"errmsg"`
}

type PingResponse struct {
	BaseMsg
	Runing		string	`json:"runing"`
}

type Token struct {
	UserId 		int64	`json:"userId"`
	Time 		int64	`json:"time"`
	LifeTime	int64	`json:"lifeTime"`
}

type BaseResEx struct {
	Errcode 	string	`json:"errcode"`
	Errmsg 		string	`json:"errmsg"`
}

type BaseResponse struct {
	Errcode 	int		`json:"errcode"`
	Errmsg 		string	`json:"errmsg"`
}

type EnterTokenResponse struct {
	BaseResponse
	Tok 	Token	`json:"token"`
}

func Testing_1(){
	array,err := json.Marshal(BaseMsg{1,"error"})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(string(array))

	ping := &PingResponse{}
	json.Unmarshal(array,ping)

	fmt.Printf("%v \n",ping)
	fmt.Printf("%v \n",ping.Errcode)
	fmt.Printf("%v \n",ping.Runing == "")
}

func Testing_2(){
	for i := 0; i < 100; i++ {
		fmt.Println(rand.Intn(2))
	}
}

func Testing_3(){
	token := EnterTokenResponse{
		BaseResponse{0,"ok"},
		Token{1,1,1},
	}
	array,err := json.Marshal(token)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(string(array))
}

func Testing_4(){
	array,err := json.Marshal(BaseResEx{"1","error"})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(string(array))

	base := &BaseResponse{}
	json.Unmarshal(array,base)
	fmt.Printf("%v",base)
}

func Testing_5()  {
	//opt := new(eio.Options)
	//opt.Transports = []ts.Transport{
	//	websocket.Default,
	//	polling.Default,
	//}
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})
	server.OnEvent("/", "login", func(s socketio.Conn, msg msg.LoginMsg) {
		fmt.Println("%v", msg)
		//s.Emit("reply", "have "+msg)
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})
	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})
	go server.Serve()
	defer server.Close()

	http.HandleFunc("/socket.io/",func(w http.ResponseWriter, r *http.Request) {
		//origin := r.Header.Get("Origin")
		//log.Println("origin",origin)
		//w.Header().Set("Access-Control-Allow-Origin", origin)
		//w.Header().Set("Access-Control-Allow-Credentials", "true")

		w.Header().Set("Access-Control-Allow-Origin", "*");
		w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With");
		w.Header().Set("Access-Control-Allow-Methods","PUT,POST,GET,DELETE,OPTIONS");
		w.Header().Set("X-Powered-By"," 3.2.1")
		w.Header().Set("Content-Type", "application/json;charset=utf-8");

		server.ServeHTTP(w, r)
	})

	//http.Handle("/socket.io/", server)
	//http.Handle("/", http.FileServer(http.Dir("../asset")))
	http.ListenAndServe("127.0.0.1:1000",nil)
}

func Testing_6(){
	server := sio.NewServer(transport.GetDefaultWebsocketTransport())

	//handle connected
	server.On(sio.OnConnection, func(c *sio.Channel) {
		log.Println("New client connected")
	})

	server.On(sio.OnDisconnection, func(c *sio.Channel) {

	})

	//handle custom event
	server.On("login", func(c *sio.Channel, msg msg.LoginMsg) string {
		log.Println("login")
		c.Emit("login",msg)
		return ""
	})

	//serveMux := http.NewServeMux()
	//serveMux.Handle("/socket.io/", server)
	http.HandleFunc("/socket.io/",func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		server.ServeHTTP(w, r)
	})
	http.ListenAndServe("127.0.0.1:1000", nil)
}

func Testing_8(){
	//array,err := json.Marshal(BaseResEx{"1","error"})
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return
	//}
	fmt.Println([]byte("'"))
	fmt.Println([]byte("\""))

	t := reflect.TypeOf(msg.LoginMsg{})
	d := reflect.New(t).Interface()
	err1 := json.Unmarshal( []byte("{\"token\":\"sdfs\",\"roomid\":\"112342\",\"time\":121233,\"sign\":\"sdfsdf\"}"), &d)
	if err1 != nil {
		log.Info(err1.Error())
		return
	}
}

func getMessage(msg string) string {
	if len(msg) >= 4 {
		if msg[:2] == "\"{" && msg[len(msg)-2:] == "}\"" {
			msg = msg[1:len(msg)]
			msg = msg[:len(msg)-1]
		}
	}

	var flag = false
	for i := len(msg)-1; i >= 0; i-- {
		if flag {
			if byte(msg[i]) == '\\' {
				if i == 0 {
					msg = msg[1:]
				}else{
					msg = msg[:i] + msg[i+1:]
				}
			}
			flag = false
		}else{
			if msg[i] == '"' {
				flag = true
			}
		}
	}
	return msg
}

func Testing_9(){
	str := `"{\"token\":\"sdfs\",\"roomid\":\"112342\",\"time\":121233,\"sign\":\"sdfsdf\"}"`
	str = getMessage(str)
	fmt.Println(str)
}

func Testing_10(){
	s := []int{1}
	//s = s[:0]
	//s = s[1:]
	s = s[:len(s)-1]
	fmt.Println(s)
}

func Testing_11(){
	m := map[int]int{}
	//if m[1] == 0 {
	//	fmt.Println("sss")
	//}
	fmt.Println(m[1])
	fmt.Println(m[2])
	m[3] = m[3] + 1
	for k,v := range m {
		fmt.Println(k,v)
	}
}

func Testing_12()  {
	m := make([]int,4)
	fmt.Println(m)
	fmt.Println(cap(m))
	t1 := func(s []int){
		s = append(s,1)
		s[0] = 1
		fmt.Println(s)
	}
	t1(m)
	fmt.Println(m)
	t4 := func(s []int){
		s = append(s,1)
		s[0] = 1
		fmt.Println(s)
	}
	t4(m[:])
	fmt.Println(m)
	t2 := func(s *[]int){
		*s = append(*s,1)
		fmt.Println(*s)
	}
	t2(&m)
	fmt.Println(m)


	mp := make(map[int]int,4)
	fmt.Println(mp)
	t3 := func(mp map[int]int){
		//s = append(s,1)
		mp[0]++
		fmt.Println(mp)
	}
	t3(mp)
	fmt.Println(mp)
}

func Testing_13(){
	str := "sss"
	str1 := "\"sss\""
	fmt.Println(str1)
	array,err := json.Marshal(str)
	if err != nil {
		fmt.Println("11")
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(array))

	err = json.Unmarshal([]byte("\"sssss\""),&str)
	if err != nil {
		fmt.Println("22")
		fmt.Println(err.Error())
		return
	}else{
		fmt.Println("33")
		fmt.Println(str)
	}

}

func Testing_14(){
	array := []int{0,0,0}
	t := func(list []int){
		list[0] = 1
	}
	t(array)
	fmt.Println(array)

	array2 := []int{0,0,0}
	ttt := func(list []int){
		list = append(list,1)
		list[0] = 1
		fmt.Println(cap(list))
	}
	ttt(array2)
	fmt.Println(array2)

	array4 := make([]int,3,6)
	ttt(array4)
	fmt.Println(array4)

	array1 := [3]int{0,0,0}
	tt := func(list [3]int){
		list[0] = 1
	}
	tt(array1)
	fmt.Println(array1)
}

func Testing_15(){
	list := []int{1,2,3,4,5}
	list2 := list[:]
	list2[0] = 5
	list2 = append(list2,1)
	fmt.Println(list2)
	fmt.Println(list)

	list2 = append([]int{},list[:]...)
	list2[0] = 10
	fmt.Println(list2)
	fmt.Println(list)
}

type Stu struct {
	Name 	string
	Score 	int
}

func say(str string,i int){
	time.Sleep(time.Millisecond * 1000)
	fmt.Printf("%s %d \n",str,i)
}

func run(str string) {
	for i := 0; i<5; i++ {
		//time.Sleep(time.Millisecond * 1000)
		//fmt.Println(fmt.Printf("%s %d",str,i))
		go say(str,i)
	}
	fmt.Println("run over")
}

func Testing_16(){
	signalChan := make(chan os.Signal, 1)
	// 启动一个goroutine线程
	go run("Hello")
	time.Sleep(time.Millisecond * 5000)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
}

func main() {

	//Testing_12()
	//Testing_9()
	//Testing_14()
}





