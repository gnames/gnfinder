package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/gnames/gnfinder/protob"
	"google.golang.org/grpc"
)

type gnfinderServer struct{}

func (gnfinderServer) Ping(ctx context.Context, void *protob.Void) (*protob.Pong, error) {
	pong := protob.Pong{Value: "pong"}
	return &pong, nil
}

func Run(port string) {
	srv := grpc.NewServer()
	var gnf gnfinderServer
	protob.RegisterGNFinderServer(srv, gnf)
	portVal := fmt.Sprintf(":%s", port)
	l, err := net.Listen("tcp", portVal)
	if err != nil {
		log.Fatalf("could not listen to %s: %v", portVal, err)
	}
	log.Fatal(srv.Serve(l))
}
