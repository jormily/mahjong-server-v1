package main

import (
	"fmt"
	"github.com/gohouse/converter"
)

func main() {
	T2S := converter.NewTable2Struct()
	T2S.Config(&converter.T2tConfig{
		StructNameToHump:true,
		// 如果字段首字母本来就是大写, 就不添加tag, 默认false添加, true不添加
		RmTagIfUcFirsted: false,
		// tag的字段名字是否转换为小写, 如果本身有大写字母的话, 默认false不转
		TagToLower: true,
		// 字段首字母大写的同时, 是否要把其他字母转换为小写,默认false不转换
		UcFirstOnly: false,
		//// 每个struct放入单独的文件,默认false,放入同一个文件(暂未提供)
		//SeperatFile: false,
	})
	err := T2S.
		SavePath("./db/struct.go").
		PackageName("db").
		RealNameMethod("TableName").
		//Prefix("t_").
		TagKey("xorm").
		Dsn("root:pwd#2019@tcp(localhost:3306)/db_babykylin?charset=utf8").
		Run()
	fmt.Println(err)
}
