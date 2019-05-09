package main

import (
	"flag"
	"log"
	"net"
	"os"
	"runtime"

	"github.com/Arktische/redirector"
)

func main() {
	// remoteMap := make(map[int16]*net.TCPAddr)
	// var err error
	runtime.GOMAXPROCS(runtime.NumCPU())
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	// remoteMap[1080], err = net.ResolveTCPAddr("tcp", "127.0.0.1:1080")
	// if err != nil {
	// 	logger.Output(2, "parse addr error")
	// }
	verbose := flag.Bool("v", false, "should every proxy request be logged to stdout")
	addr := flag.String("addr", ":8888", "proxy listen address")
	flag.Parse()

	redirector.Verbose = *verbose

	ln, err := net.Listen("tcp", *addr)
	if err != nil {
		panic(err)
	}

	logger.Printf("Listening %s \n", *addr)
	pool, _ := redirector.NewPool(2000)
	handler := func(conn interface{}) {
		redirector.HandleConnection(conn.(net.Conn))
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			logger.Println(err)
			return
		}
		if redirector.Verbose {
			logger.Println("new client:", conn.RemoteAddr())
		}

		err = pool.Submit(handler, conn)
		if err != nil {
			logger.Println(err)
		}
	}
}
