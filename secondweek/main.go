package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	g, ctx := errgroup.WithContext(context.Background())
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("pong"))
	})
	// 模拟单个服务错误退出
	serverOut := make(chan struct{})
	mux.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		serverOut <- struct{}{}
	})
	server := http.Server{
		Handler: mux,
		Addr:    ":8099",
	}

	//g1	// g1 退出了所有的协程都能退出么？
	//g1 退出后, context 将不再阻塞，g2, g3 都会随之退出
	//然后 main 函数中的 g.Wait() 退出，所有协程都会退出
	g.Go(func() error {
		err := server.ListenAndServe() // 服务启动后会阻塞, 虽然使用的是 go 启动，但是由于 g.WaitGroup 试得其是个阻塞的 协程
		if err != nil {
			log.Println("g1 error,will exit.", err.Error())
		}
		return err
	})
	// g2	// g2 退出了所有的协程都能退出么？
	//到调用 `/shutdown`接口时, serverOut 无缓冲管道写入数据， case接收到数据后执行server.shutdown, 此时 g1 httpServer会退出
	//g1退出后，会返回error,将error加到g中，同时会调用 cancel()
	//g3 中会 select case ctx.Done, context 将不再阻塞，g3 会随之退出
	//然后 main 函数中的 g.Wait() 退出，所有协程都会退出
	g.Go(func() error {
		select {
		case <-ctx.Done():
			log.Println("g2 errgroup exit...")
		case <-serverOut:
			log.Println("g2, request `/shutdown`, server will out...")
		}
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		//这里不是必须的，但是如果使用 _ 的话静态扫描工具会报错，加上也无伤大雅
		defer cancel()
		err := server.Shutdown(timeoutCtx)
		log.Println("shutting down server...")
		return err
	})
	//g3	// g3 捕获到 os 退出信号将会退出
	//g3 退出了所有的协程都能退出么？
	//g3 退出后, context 将不再阻塞，g2 会随之退出
	//g2 退出时，调用了 shutdown，g1 会退出
	//然后 main 函数中的 g.Wait() 退出，所有协程都会退出
	g.Go(
		func() error {
			quit := make(chan os.Signal, 0)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
			select {
			case <-ctx.Done():
				log.Println("g3, ctx execute cancel...")
				log.Println("g3 error,", ctx.Err().Error())
				// 当g2退出时，已经有错误了，此时的error 并不会覆盖到g中
				return ctx.Err()
			case sig := <-quit:
				return fmt.Errorf("g3 get os signal: %v", sig)
			}
		})
	//g.Wait 等待所有 go执行完毕后执行
	fmt.Printf("end, errgroup exiting, %+v\n", g.Wait())

}
