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

func (gnfinderServer) FindNames(ctx context.Context,
	params *protob.Params) (*protob.NameStrings, error) {
	var names protob.NameStrings
	// text := params.Text
	return &names, nil
}

func Run(port int) {
	var gnf gnfinderServer
	srv := grpc.NewServer()
	protob.RegisterGNFinderServer(srv, gnf)
	portVal := fmt.Sprintf(":%d", port)
	l, err := net.Listen("tcp", portVal)
	if err != nil {
		log.Fatalf("could not listen to %s: %v", portVal, err)
	}
	log.Fatal(srv.Serve(l))
}
