package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hashicorp/memberlist"
)

func main() {
	// https://qiita.com/octu0/items/808299d232bc003d5e99
	nodeName := flag.String("name", "hoge", "node name")
	masterNode := flag.String("masterAddr", "127.0.0.1", "master addr")
	flag.Parse()

	conf := memberlist.DefaultLocalConfig()
	conf.Name = *nodeName

	list, err := memberlist.Create(conf)
	if err != nil {
		log.Fatal(err)
	}

	local := list.LocalNode()
	log.Printf("%s at %s:%d", local.Name, local.Addr.To4().String(), local.Port)

	log.Printf("wait for other member connections")
	if *masterNode != local.Addr.To4().String() {
		if _, err := list.Join([]string{*masterNode}); err != nil {
			log.Fatal(err)
		}
	}

	for _, member := range list.Members() {
		log.Printf("Member: %s(%s:%d)", member.Name, member.Addr.To4().String(), member.Port)
	}

	defer func() {
		log.Printf("bye")
	}()
	signal_chan := make(chan os.Signal, 2)
	signal.Notify(signal_chan, syscall.SIGTERM)
	signal.Notify(signal_chan, syscall.SIGINT)

	log.Printf("wait for signal: pid=%d", os.Getpid())
	for {
		select {
		case s := <-signal_chan:
			switch s {
			case syscall.SIGINT:
				log.Printf("SIGINT happen. cluter leaving")
				timeout := 10 * time.Second
				if err := list.Leave(timeout); err != nil {
					log.Fatal(err)
				}
				log.Printf("cluter left.")
			case syscall.SIGTERM:
				log.Printf("SIGTERM happen. bye.")
				return
			}
		}
	}
}
