package tool

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
)

func Rand(min,max int) int {
	if max < min {
		return 0
	}
	c := max - min
	r := rand.Intn(c+1)
	return r + min
}

func Response(w http.ResponseWriter,msg interface{}){
	w.Header().Set("Access-Control-Allow-Origin", "*")             	//允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") 	//header的类型
	w.Header().Set("content-type", "application/json")

	data,err := json.Marshal(msg)
	if err != nil {
		log.Info("Response err %s",err.Error())
		return
	}

	w.Write(data)
}

func Slice2Map(s []int) map[int]int {
	m := make(map[int]int)
	for _,v := range s {
		if _,ok := m[v];ok {
			m[v] = m[v] + 1
		}
	}
	return m
}

func SliceInsert(s *[]interface{}, index int, value interface{}) {
	if (index >= 0 && index < len(*s)) {
		return
	}
	*s = append(append((*s)[:index], value), (*s)[index:]...)
}


func SliceRemove(s *[]int,v int,args...int) {
	cnt := 1
	if len(args) == 1 {
		cnt = args[0]
	}

	for i:=len(*s)-1;i>=0;i-- {
		if (*s)[i] == v {
			(*s) = append((*s)[:i], (*s)[i+1:]...)
			cnt--

			if cnt == 0 {
				return
			}
		}
	}
}


func SliceRemoveByValue(s []int,v int,args...int) []int {
	cnt := 1
	if len(args) == 1 {
		cnt = args[0]
	}

	for i:=len(s)-1;i>=0;i-- {
		if s[i] == v {
			s = append(s[:i], s[i+1:]...)
			cnt--

			if cnt == 0 {
				return s
			}
		}
	}
	return s
}

func MD5(text string) string {
	ctx := md5.New()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}

func HttpGet(url string,path string,values url.Values) ([]byte,error){
	url = "http://"+url
	url = url + path
	if len(values) > 0 {
		url = url + "?" + values.Encode()
	}

	var res,err = http.Get(url)
	//log.Info("HttpGet:"+url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var body,err1 = ioutil.ReadAll(res.Body)
	if err1 != nil {
		return nil,err1
	}
	return body,nil
}