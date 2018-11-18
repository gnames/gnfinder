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
	"github.com/gnames/gnfinder/util"
	"github.com/gnames/gnfinder/verifier"
	"google.golang.org/grpc"
)

type gnfinderServer struct{}

var dictionary dict.Dictionary

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

func (gnfinderServer) Ping(ctx context.Context,
	void *protob.Void) (*protob.Pong, error) {
	pong := protob.Pong{Value: "pong"}
	return &pong, nil
}

func (gnfinderServer) FindNames(ctx context.Context,
	params *protob.Params) (*protob.NameStrings, error) {
	text := params.Text
	opts := setOpts(params)
	m := util.NewModel(opts...)
	output := gnfinder.FindNames([]rune(string(text)), &dictionary, m)

	if m.Verifier.Verify {
		names := gnfinder.UniqueNameStrings(output.Names)
		namesVerified := verifier.Verify(names, m)
		for i, n := range output.Names {
			if v, ok := namesVerified[n.Name]; ok {
				output.Names[i].Verification = v
			}
		}
	}

	names := protobNameStrings(&output)

	return &names, nil
}

func setOpts(params *protob.Params) []util.Opt {
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
	return opts
}

func protobNameStrings(output *gnfinder.Output) protob.NameStrings {
	var names []*protob.NameString
	for _, n := range output.Names {
		name := protob.NameString{
			Type:         n.Type,
			Verbatim:     n.Verbatim,
			Name:         n.Name,
			Odds:         float32(n.Odds),
			OffsetStart:  int32(n.OffsetStart),
			OffsetEnd:    int32(n.OffsetEnd),
			Verification: verification(&n.Verification),
		}
		names = append(names, &name)
	}
	return protob.NameStrings{Names: names}
}

func verification(ver *verifier.Verification) *protob.Verification {
	return &protob.Verification{
		DataSourceId:       int32(ver.DataSourceID),
		DataSourceTitle:    ver.DataSourceTitle,
		MatchedName:        ver.MatchedName,
		CurrentName:        ver.CurrentName,
		ClassificationPath: ver.ClassificationPath,
		DataSourcesNum:     int32(ver.DataSourcesNum),
		DataSourceQuality:  ver.DataSourceQuality,
		EditDistance:       int32(ver.EditDistance),
		StemEditDistance:   int32(ver.StemEditDistance),
		MatchType:          getMatchType(ver.MatchType),
		Error:              ver.Error,
		PreferredResults:   sourcesResult(ver),
	}
}

func sourcesResult(ver *verifier.Verification) []*protob.PreferredResult {
	l := len(ver.PreferredResults)
	res := make([]*protob.PreferredResult, l)
	for i, v := range ver.PreferredResults {
		res[i] = &protob.PreferredResult{
			DataSourceId:    int32(v.DataSourceID),
			DataSourceTitle: v.DataSourceTitle,
			NameId:          v.NameID,
			Name:            v.Name,
			TaxonId:         v.TaxonID,
		}
	}
	return res
}

func getMatchType(match string) protob.MatchType {
	switch match {
	case "Exact":
		return protob.MatchType_EXACT
	case "FuzzyCanonicalMatch":
		return protob.MatchType_FUZZY
	case "ExactPartialMatch":
		return protob.MatchType_PARTIAL_EXACT
	case "FuzzyPartialMatch":
		return protob.MatchType_PARTIAL_FUZZY
	}
	return protob.MatchType_NONE
}
