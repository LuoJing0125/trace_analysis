package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	Lfiles, _ := ioutil.ReadDir("./resource_limit_pod") //读取资源限制log日期列表
	rl, err := os.Create("./Limit_date.log")            //创建文件
	if err != nil {
		fmt.Println("data File creating error", err)
		return
	}
	for _, f := range Lfiles {
		// fmt.Println(f.Name())
		rl.WriteString(f.Name() + "\n") //写入文件(字节数组)
	}
	rl.Close()
}
