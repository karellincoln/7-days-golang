package main

import (
	"fmt"
	geerpc "github.com/karellincoln/7-day-golang/gee-rpc"
	"github.com/karellincoln/7-day-golang/gee-rpc/log"
	"sync"

	"net"
	"time"
)

func startServer(addr chan string) {
	// pick a free port
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Error("network error:", err)
	}
	log.Info("start rpc server on", l.Addr())
	addr <- l.Addr().String()
	geerpc.Accept(l)
}

func main() {
	addr := make(chan string)
	go startServer(addr)

	// in fact, following code is like a simple geerpc client
	client, _ := geerpc.Dial("tcp", <-addr)
	defer func() { _ = client.Close() }()

	time.Sleep(time.Second)
	// send request & receive response
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			args := fmt.Sprintf("geerpc req %d", i)
			var reply string
			if err := client.Call("Foo.Sum", args, &reply); err != nil {
				log.Error("call Foo.Sum error:", err)
			}
			log.Info("reply:", reply)
		}(i)
	}
	wg.Wait()
}
