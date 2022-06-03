package main

/********************************************************
参数:
	param1: 旧的excel目录
	param2: 新的excel目录
*********************************************************/

import (
	"flag"
	"fmt"
	_ "net/http/pprof"
	"os"
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
		fmt.Fprintf(os.Stderr, ` 用于对比游戏服务器版本迭代中Excel目录的差别.
param1: 旧的excel目录
param2: 新的excel目录`+"\n")
		return
	}

	fmt.Println("reading old directory")
	oldContents := fetchDirectoryContents(flag.Args()[0])

	fmt.Println("reading new directory")
	newContents := fetchDirectoryContents(flag.Args()[1])

	fmt.Println("diffing")
	diffContents(oldContents, newContents)

	fmt.Println("end, please check result.cvs in the running path")
}
