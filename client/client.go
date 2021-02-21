package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"

	"github.com/ijsong/grpctest/pb"
)

var (
	address     string
	connTimeout time.Duration
	callTimeout time.Duration
	count       int
)

func init() {
	flag.StringVar(&address, "a", "127.0.0.1:9997", "server address")
	flag.DurationVar(&connTimeout, "conn-timeout", 1*time.Second, "connection timeout")
	flag.DurationVar(&callTimeout, "call-timeout", 1*time.Second, "call timeout")
	flag.IntVar(&count, "c", 1, "ping-pong count")
	flag.Parse()
}

func call(i int) error {
	connCtx, connCancel := context.WithTimeout(context.Background(), connTimeout)
	defer connCancel()

	conn, err := grpc.DialContext(connCtx, address, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewPingPongClient(conn)

	callCtx, callCancel := context.WithTimeout(context.Background(), callTimeout)
	defer callCancel()
	rsp, err := client.Call(callCtx, &pb.Ping{
		Msg: fmt.Sprintf("ping-%s", time.Now().String()),
	})
	if err != nil {
		return err
	}
	log.Printf("pingpong: %+v", rsp)
	return nil
}

func main() {
	for i := 0; i < count; i++ {
		if err := call(i + 1); err != nil {
			log.Fatal(err)
		}
	}
}
