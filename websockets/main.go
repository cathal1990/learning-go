package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/cathal1990/learning-go/websockets/internal/hardware"
	"nhooyr.io/websocket"
)

type server struct {
	subscriberMessageBuffer int
	mux                     http.ServeMux
	subscribersMutex        sync.Mutex
	subscribers             map[*subscriber]struct{}
}

type subscriber struct {
	msgs chan []byte
}

func NewServer() *server {
	s := &server{
		subscriberMessageBuffer: 10,
		subscribers:             make(map[*subscriber]struct{}),
	}

	s.mux.Handle("/", http.FileServer(http.Dir("./htmx")))
	s.mux.HandleFunc("/ws", s.subscribeHandler)
	return s
}

func (s *server) subscribeHandler(writer http.ResponseWriter, req *http.Request) {
	err := s.subscribe(req.Context(), writer, req)

	if err != nil {
		fmt.Println(err)
		return
	}
}

func (s *server) subscribe(ctx context.Context, writer http.ResponseWriter, req *http.Request) error {
	var conn *websocket.Conn
	subscriber := &subscriber{
		msgs: make(chan []byte, s.subscriberMessageBuffer),
	}

	s.addSubscriber(subscriber)

	conn, err := websocket.Accept(writer, req, nil)

	if err != nil {
		return err
	}

	defer conn.CloseNow()

	ctx = conn.CloseRead(ctx)
	for {
		select {
		case msg := <-subscriber.msgs:
			ctx, cancel := context.WithTimeout(ctx, time.Second*5)
			defer cancel()
			err := conn.Write(ctx, websocket.MessageText, msg)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (s *server) addSubscriber(subscriber *subscriber) {
	s.subscribersMutex.Lock()
	s.subscribers[subscriber] = struct{}{}
	s.subscribersMutex.Unlock()
	fmt.Println("Added subscriber", subscriber)
}

func (s *server) broadcast(msg []byte) {
	s.subscribersMutex.Lock()
	for subscriber := range s.subscribers {
		subscriber.msgs <- msg
	}
	s.subscribersMutex.Unlock()
}

func main() {
	srv := NewServer()

	go func(s *server) {
		for {

			memStats, err := hardware.GetMemory()

			if err != nil {
				fmt.Println(err)
			}

			cpuStats, err := hardware.GetCpu()

			if err != nil {
				fmt.Println(err)
			}

			diskStats, err := hardware.GetDisk()

			if err != nil {
				fmt.Println(err)
			}

			timestamp := time.Now().Format("2006-01-02 15:04:05")

			html := `<div hx-swap-oob="innerHTML:#update-timestamp"> ` + timestamp + `</div>
			<div hx-swap-oob="innerHTML:#system-data"> ` + memStats + `</div
			<div hx-swap-oob="innerHTML:#cpu-data"> ` + cpuStats + `</div>
			<div hx-swap-oob="innerHTML:#disk-data"> ` + diskStats + `</div>`

			s.broadcast([]byte(html))

			time.Sleep(3 * time.Second)
		}
	}(srv)

	err := http.ListenAndServe(":8080", &srv.mux)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
