package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/hashicorp/memberlist"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	nodeName := flag.String("name", "hoge", "node name")
	masterNode := flag.String("masterAddr", "127.0.0.1", "master addr")
	flag.Parse()

	// coluster settings
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

	// 終了ログ
	defer func() {
		log.Printf("bye")
	}()

	// http
	var server http.Server
	go func() {
		defer wg.Done()

		server = http.Server{
			Addr:    ":8000",
			Handler: http.HandlerFunc(requestHander),
		}

		if err := server.ListenAndServe(); err != nil {
			log.Print(err)
		}
	}()

	// 終了条件
	sigCh := make(chan os.Signal, 2)
	signal.Notify(sigCh, syscall.SIGTERM)
	signal.Notify(sigCh, syscall.SIGINT)

	go func() {
		log.Printf("wait for signal: pid=%d", os.Getpid())
		select {
		case s := <-sigCh:
			switch s {
			case syscall.SIGINT:
				log.Printf("SIGINT happen. cluter leaving")
				timeout := 10 * time.Second
				if err := list.Leave(timeout); err != nil {
					log.Fatal(err)
				}
				log.Printf("cluter left.")
				ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
				if err := server.Shutdown(ctx); err != nil {
					log.Print(err)
				}
			case syscall.SIGTERM:
				log.Printf("SIGTERM happen. bye.")
				ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
				if err := server.Shutdown(ctx); err != nil {
					log.Print(err)
				}
			}
		}
	}()

	// メインスレッドを待機
	wg.Wait()
}

func requestHander(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Fprintln(w, "GET called!!")
	} else if r.Method == "POST" {
		fmt.Fprintln(w, "POST called!!")
	} else if r.Method == "PUT" {
		fmt.Fprintln(w, "PUT called!!")
	} else if r.Method == "DELETE" {
		fmt.Fprintln(w, "DELETE called!!")
	}
}
