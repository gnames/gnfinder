package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/gnames/gnfinder"
	"github.com/gnames/gnfinder/dict"
	"github.com/gnames/gnfinder/lang"
	"github.com/gnames/gnfinder/protob"
	"github.com/gnames/gnfinder/resolver"
	"github.com/gnames/gnfinder/util"
	"google.golang.org/grpc"
)

type gnfinderServer struct{}

var dictionary dict.Dictionary

func (gnfinderServer) Ping(ctx context.Context,
	void *protob.Void) (*protob.Pong, error) {
	pong := protob.Pong{Value: "pong"}
	return &pong, nil
}

func (gnfinderServer) FindNames(ctx context.Context,
	params *protob.Params) (*protob.NameStrings, error) {
	text := params.Text
	var opts []util.Opt

	if params.WithBayes {
		opts = append(opts, util.WithBayes(true))
	}
	if params.WithVerification {
		opts = append(opts, util.WithVerification(true))
	}

	if len(params.Language) > 0 {
		l, err := lang.NewLanguage(params.Language)
		if err == nil {
			opts = append(opts, util.WithLanguage(l))
		}
	}

	if len(params.Sources) > 0 {
		sources := make([]int, len(params.Sources))
		for i, v := range params.Sources {
			sources[i] = int(v)
		}
		opts = append(opts, util.WithSources(sources))
	}

	m := util.NewModel(opts...)
	output := gnfinder.FindNames([]rune(string(text)), &dictionary, m)

	if m.Resolver.Verify {
		names := gnfinder.UniqueNameStrings(output.Names)
		namesResolved := resolver.Verify(names, m)
		for i, n := range output.Names {
			if v, ok := namesResolved[n.Name]; ok {
				output.Names[i].Verification = v
			}
		}
	}

	names := protobNameStrings(&output)

	return &names, nil
}

func protobNameStrings(output *gnfinder.Output) protob.NameStrings {
	var names []*protob.NameString
	for _, n := range output.Names {
		name := protob.NameString{
			Value:    n.Name,
			Verbatim: n.Verbatim,
			Odds:     float32(n.Odds),
		}
		name.Path = n.Verification.ClassificationPath
		name.Curated = n.Verification.DataSourceQuality == "HasCuratedSources"
		name.Match = getMatchType(n.Verification.MatchType)
		name.SourceId = int32(n.Verification.DataSourceID)
		names = append(names, &name)
	}
	return protob.NameStrings{Names: names}
}

func getMatchType(match string) protob.MatchType {
	switch match {
	case "ExactMatch":
		return protob.MatchType_EXACT
	case "ExactCanonicalMatch":
		return protob.MatchType_CANONICAL_EXACT
	case "FuzzyCanonicalMatch":
		return protob.MatchType_CANONICAL_FUZZY
	case "ExactPartialMatch":
		return protob.MatchType_PARTIAL_EXACT
	case "FuzzyPartialMatch":
		return protob.MatchType_PARTIAL_FUZZY
	}
	return protob.MatchType_NONE
}

func Run(port int) {
	var gnf gnfinderServer
	srv := grpc.NewServer()
	dictionary = dict.LoadDictionary()
	protob.RegisterGNFinderServer(srv, gnf)
	portVal := fmt.Sprintf(":%d", port)
	l, err := net.Listen("tcp", portVal)
	if err != nil {
		log.Fatalf("could not listen to %s: %v", portVal, err)
	}
	log.Fatal(srv.Serve(l))
}
