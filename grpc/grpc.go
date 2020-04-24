package grpc

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/gnames/bayes"
	"github.com/gnames/gnfinder"
	"github.com/gnames/gnfinder/dict"
	"github.com/gnames/gnfinder/lang"
	"github.com/gnames/gnfinder/nlp"
	"github.com/gnames/gnfinder/output"
	"github.com/gnames/gnfinder/protob"
	"github.com/gnames/gnfinder/verifier"
	"google.golang.org/grpc"
)

type gnfinderServer struct{}

var dictionary *dict.Dictionary
var weights map[lang.Language]*bayes.NaiveBayes

func Run(port int) {
	var gnfs gnfinderServer
	srv := grpc.NewServer()
	dictionary = dict.LoadDictionary()
	weights = nlp.BayesWeights()
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

func (gnfinderServer) Ver(ctx context.Context,
	void *protob.Void) (*protob.Version, error) {
	ver := protob.Version{Version: gnfinder.Version}
	return &ver, nil
}

func (gnfinderServer) FindNames(ctx context.Context,
	params *protob.Params) (*protob.Output, error) {
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
	opts := []gnfinder.Option{
		gnfinder.OptDict(dictionary),
		gnfinder.OptBayesWeights(weights),
	}

	opts = append(opts, gnfinder.OptBayes(!params.NoBayes))

	if params.Verification {
		var verOpts []verifier.Option
		var sources []int
		for _, v := range params.Sources {
			sources = append(sources, int(v))
		}
		verOpts = append(verOpts, verifier.OptSources(sources))

		opts = append(opts, gnfinder.OptVerify(verOpts...))
	}

	if params.DetectLanguage {
		opts = append(opts, gnfinder.OptDetectLanguage(true))
	} else if len(params.Language) > 0 {
		l, err := lang.NewLanguage(params.Language)
		if err == nil {
			opts = append(opts, gnfinder.OptLanguage(l))
		}
	}

	if params.TokensAround > 0 {
		opts = append(opts, gnfinder.OptTokensAround(int(params.TokensAround)))
	}

	return opts
}

func protobNameStrings(out *output.Output) protob.Output {
	var names []*protob.NameString
	for _, n := range out.Names {
		name := &protob.NameString{
			Type:           n.Type,
			Verbatim:       n.Verbatim,
			Name:           n.Name,
			Odds:           float32(n.Odds),
			AnnotNomen:     n.AnnotNomen,
			AnnotNomenType: getAnnotNomenType(n.AnnotNomenType),
			OffsetStart:    int32(n.OffsetStart),
			OffsetEnd:      int32(n.OffsetEnd),
			WordsBefore:    n.WordsBefore,
			WordsAfter:     n.WordsAfter,

			Verification: verification(n.Verification),
		}
		names = append(names, name)
	}
	ns := protob.Output{
		Date:             out.Date.String(),
		FinderVersion:    out.FinderVersion,
		Language:         out.Language,
		LanguageDetected: out.LanguageDetected,
		DetectLanguage:   out.DetectLanguage,
		TotalTokens:      int32(out.TotalTokens),
		TotalCandidates:  int32(out.TotalNameCandidates),
		TotalNames:       int32(out.TotalNames),
		Names:            names,
	}

	return ns
}

func verification(ver *verifier.Verification) *protob.Verification {
	if ver == nil {
		var protoVer *protob.Verification
		return protoVer
	}
	match := ver.BestResult
	protoVer := &protob.Verification{
		BestResult:        buildResult(match),
		PreferredResults:  preferredResults(ver),
		DataSourcesNum:    int32(ver.DataSourcesNum),
		DataSourceQuality: ver.DataSourceQuality,
		Retries:           int32(ver.Retries),
		Error:             ver.Error,
	}
	return protoVer
}

func buildResult(res *verifier.ResultData) *protob.ResultData {
	rd := &protob.ResultData{
		DataSourceId:       int32(res.DataSourceID),
		DataSourceTitle:    res.DataSourceTitle,
		TaxonId:            res.TaxonID,
		MatchedName:        res.MatchedName,
		MatchedCanonical:   res.MatchedCanonical,
		CurrentName:        res.CurrentName,
		ClassificationPath: res.ClassificationPath,
		ClassificationRank: res.ClassificationRank,
		ClassificationIds:  res.ClassificationIDs,
		EditDistance:       int32(res.EditDistance),
		StemEditDistance:   int32(res.StemEditDistance),
		MatchType:          getMatchType(res.MatchType),
	}

	return rd
}

func preferredResults(ver *verifier.Verification) []*protob.ResultData {
	res := make([]*protob.ResultData, len(ver.PreferredResults))
	for i, v := range ver.PreferredResults {
		res[i] = buildResult(v)
	}
	return res
}

func getAnnotNomenType(match string) protob.AnnotNomenType {
	switch match {
	case "SP_NOV":
		return protob.AnnotNomenType_SP_NOV
	case "COMB_NOV":
		return protob.AnnotNomenType_COMB_NOV
	case "SUBSP_NOV":
		return protob.AnnotNomenType_SUBSP_NOV
	}
	return protob.AnnotNomenType_NO_ANNOT
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
