package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"

	"github.com/spf13/viper"

	"github.com/jormily/mahjong-server/db"
	"github.com/jormily/mahjong-server/handler"
	"github.com/jormily/mahjong-server/service"
)

func init() {
	log.Println("init")
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln("read config failed: %v", err)
	}
}

func dbStartup() func() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.psw"),
		viper.GetString("mysql.host"),
		viper.GetString("mysql.port"),
		viper.GetString("mysql.db"),
		viper.GetString("mysql.args"))

	return db.MustStartup(
		dsn,
		db.MaxIdleConns(viper.GetInt("database.max_idle_conns")),
		db.MaxIdleConns(viper.GetInt("database.max_open_conns")),
		db.ShowSQL(viper.GetBool("database.show_sql")))
}

func AccountStartup(){
	log.Info("AccountStart")
	dbStartup()
	handler.AccountServe(viper.GetString("act_svr.host"))
}

func HallStartup(){
	log.Info("HallStart")
	dbStartup()
	handler.HallServe(viper.GetString("hall_svr.host"))
}

func GameStartup()  {
	dbStartup()
	service.SocketServStart()
	handler.GameServe(viper.GetString("game_svr.ip"),viper.GetString("game_svr.http_port"))
}

func main()  {
	switch os.Args[1] {
	case "accout":
		AccountStartup()
	case "hall":
		HallStartup()
	case "game":
		GameStartup()
	}

}