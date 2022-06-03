package main

/********************************************************
参数:
	param1: 旧的excel目录
	param2: 新的excel目录
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
