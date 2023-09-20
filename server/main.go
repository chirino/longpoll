package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	defer cancel()
	StartServer()
	select {
	case <-ctx.Done():
	}
}
func StartServer() {
	// starts a http server on port 8000 that has a long poll endpoint
	routes := http.NewServeMux()
	routes.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		flusher, ok := writer.(http.Flusher)
		if !ok {
			panic("can't cast to a flusher")
		}

		interval := 30 * time.Second
		intervalStr := request.URL.Query().Get("interval")
		if intervalStr != "" {
			i, err := strconv.Atoi(intervalStr)
			if err != nil {
				writer.WriteHeader(400)
				writer.Write([]byte(`invalid interval query value`))
				return
			}
			interval = time.Duration(i) * time.Second
		}

		writer.WriteHeader(200)
		flusher.Flush()

		for {
			writer.Write([]byte(fmt.Sprintf("{\"kind\":\"keep-alive\", \"interval\":%d}\n", int(interval.Seconds()))))
			flusher.Flush()
			time.Sleep(interval)
		}
	})
	l, err := net.Listen("tcp", ":8000")
	if err != nil {
		panic(err)
	}
	go func() {
		fmt.Println("http server listening on", l.Addr().String())
		_ = http.Serve(l, routes)
	}()
}
