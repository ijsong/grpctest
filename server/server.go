package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"

	"google.golang.org/grpc"

	"github.com/ijsong/grpctest/pb"
)

var listenAddr string

type server struct {
	s *grpc.Server
	n uint64
}

func newServer() (*server, error) {
	s := &server{
		s: grpc.NewServer(),
	}
	return s, nil
}

func (s *server) run() error {
	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}
	log.Printf("server runs (%s)", listenAddr)
	pb.RegisterPingPongServer(s.s, s)
	return s.s.Serve(lis)
}

func (s *server) close() {
	s.s.Stop()
}

func (s *server) Call(ctx context.Context, ping *pb.Ping) (*pb.Pong, error) {
	seq := atomic.AddUint64(&s.n, 1)
	pong := &pb.Pong{
		Msg: fmt.Sprintf("pong-%d", seq),
	}
	if seq%1000 == 0 {
		log.Printf("request: %s, response: %s", ping.String(), pong.String())
	}
	return pong, nil
}

func init() {
	flag.StringVar(&listenAddr, "l", "0.0.0.0:9997", "listen address")
	flag.Parse()
}

func main() {
	s, err := newServer()
	if err != nil {
		log.Fatal(err)
	}

	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, os.Interrupt, syscall.SIGTERM)
	go func() {
		select {
		case <-sigC:
			s.close()
		}
	}()

	if err := s.run(); err != nil {
		log.Fatal(err)
	}
}
