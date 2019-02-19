package grpc

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/gnames/gnfinder"
	"github.com/gnames/gnfinder/dict"
	"github.com/gnames/gnfinder/lang"
	"github.com/gnames/gnfinder/output"
	"github.com/gnames/gnfinder/protob"
	"github.com/gnames/gnfinder/verifier"
	"google.golang.org/grpc"
)

type gnfinderServer struct{}

var dictionary *dict.Dictionary

func Run(port int) {
	var gnfs gnfinderServer
	srv := grpc.NewServer()
	dictionary = dict.LoadDictionary()
	protob.RegisterGNFinderServer(srv, gnfs)
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
	gnf := gnfinder.NewGNfinder(opts...)
	res := gnf.FindNames([]byte(text))

	if gnf.Verifier != nil {
		verifiedNames := gnf.Verifier.Run(res.UniqueNameStrings())
		res.MergeVerification(verifiedNames)
	}

	names := protobNameStrings(res)

	return &names, nil
}

func setOpts(params *protob.Params) []gnfinder.Option {
	opts := []gnfinder.Option{gnfinder.OptDict(dictionary)}

	if params.WithBayes {
		opts = append(opts, gnfinder.OptBayes(true))
	}

	if params.WithVerification {
		var verOpts []verifier.Option
		var sources []int
		for _, v := range params.Sources {
			sources = append(sources, int(v))
		}
		verOpts = append(verOpts, verifier.OptSources(sources))

		opts = append(opts, gnfinder.OptVerify(verOpts...))
	}

	if len(params.Language) > 0 {
		l, err := lang.NewLanguage(params.Language)
		if err == nil {
			opts = append(opts, gnfinder.OptLanguage(l))
		}
	}

	return opts
}

func protobNameStrings(out *output.Output) protob.NameStrings {
	var names []*protob.NameString
	for _, n := range out.Names {
		name := &protob.NameString{
			Type:         n.Type,
			Verbatim:     n.Verbatim,
			Name:         n.Name,
			Odds:         float32(n.Odds),
			OffsetStart:  int32(n.OffsetStart),
			OffsetEnd:    int32(n.OffsetEnd),
			Verification: verification(n.Verification),
		}
		names = append(names, name)
	}
	return protob.NameStrings{Names: names}
}

func verification(ver *verifier.Verification) *protob.Verification {
	if ver == nil {
		var protoVer *protob.Verification
		return protoVer
	}
	protoVer := &protob.Verification{
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
	return protoVer
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
	case "ExactMatch":
		return protob.MatchType_EXACT
	case "ExactCanonicalMatch":
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
