package main

import (
	"flag"
	"fmt"
	"github.com/hinego/gfen/genx"
	_ "github.com/hinego/gfen/logic"
	"github.com/hinego/gfen/ssr"
	"log"
)

func main() {
	// 定义命令行参数并设置默认值
	action := flag.String("act", "type", "Act路径")
	daoPath := flag.String("dao", "", "DAO路径")
	modelPath := flag.String("model", "", "Model路径")
	typePath := flag.String("type", "", "Type路径")
	// 解析命令行参数
	flag.Parse()
	log.Println(*action)
	if *daoPath == "" || *modelPath == "" || *typePath == "" {
		fmt.Println("错误：请提供dao、model和type的路径参数")
		return
	}
	var err error
	err = ssr.Dao().GenModel(&genx.DaoInput{
		DaoPath:   *daoPath,
		ModelPath: *modelPath,
		TypePath:  *typePath,
	})
	log.Println(err)
}
