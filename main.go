package main

/********************************************************
参数:
	param1: 扫描目录，如果有多个，则使用,分割，后面目录的同名文件将覆盖前面的
	param2: 存储类型
	param3: 存储参数
	param4: token，可选参数
示例：
    ./main . disk /tmp/gocache/
	./main . redis '{"conn":"192.168.1.101:7478","password":"4rT35Ker4m"}'
*********************************************************/

import (
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
)

const (
	Sheet_Least_Row  = 4
	Sheet_Max_Column = 1000
)

var (
	help bool
)

func init() {
	flag.BoolVar(&help, "h", false, "help")
	flag.BoolVar(&help, "help", false, "help")
}
func usage() {
	flag.PrintDefaults()
}

func main() {
	flag.Parse()
	if help {
		fmt.Fprintf(os.Stderr, ``+"\n")
		return
	}

	go func() {
		http.ListenAndServe("0.0.0.0:9999", nil)
	}()

	fmt.Println("reading old directory")
	oldContents := fetchDirectoryContents(flag.Args()[0])

	fmt.Println("reading new directory")
	newContents := fetchDirectoryContents(flag.Args()[1])

	fmt.Println("diffing")
	diffContents(oldContents, newContents)

	fmt.Println("end, please check result.cvs in the running path")
}
