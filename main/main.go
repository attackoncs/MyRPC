package main

import (
	"fmt"
	"log"
	"myrpc"
	"net"
	"sync"
	"time"
)

func startServer(addr chan string) {
	// pick a free port
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("network error:", err)
	}
	log.Println("start rpc server on", l.Addr())
	addr <- l.Addr().String()
	myrpc.Accept(l)
}

func main() {
	log.SetFlags(0)
	addr := make(chan string)
	//服务器部分代码，下面是客户端代码
	go startServer(addr)

	//客户端代码
	client, _ := myrpc.Dial("tcp", <-addr)
	defer func() { _ = client.Close() }()

	time.Sleep(time.Second)
	// send options
	var wg sync.WaitGroup

	// send request & receive response
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			args := fmt.Sprintf("myrpc req %d", i)
			var reply string
			if err := client.Call("Foo.Sum", args, &reply); err != nil {
				log.Fatal("call Foo.Sum error:", err)
			}
			log.Println("reply:", reply)
		}(i)
	}
	wg.Wait()
}
