package rest

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gnames/gnfinder"
	"github.com/gnames/gnfinder/config"
	"github.com/gnames/gnfinder/ent/api"
	"github.com/gnames/gnfinder/ent/nlp"
	"github.com/gnames/gnfinder/ent/output"
	"github.com/gnames/gnfinder/ent/verifier"
	"github.com/gnames/gnfinder/io/dict"
	"github.com/gnames/gnfmt"
	"github.com/labstack/echo/v4"
	"github.com/tj/assert"
)

var (
	dictionary = dict.LoadDictionary()
	weights    = nlp.BayesWeights()
)

func TestGet(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	cfg := config.New()
	gnf := gnfinder.New(cfg, dictionary, weights)

	t.Run("test ping", func(t *testing.T) {
		c.SetPath("/api/v1/ping")
		assert.Nil(t, ping()(c))
		assert.Equal(t, rec.Code, http.StatusOK)
		assert.Equal(t, rec.Body.String(), "pong")
	})

	t.Run("test version", func(t *testing.T) {
		c.SetPath("/api/v1/version")
		assert.Nil(t, ver(gnf)(c))
		assert.Equal(t, rec.Code, http.StatusOK)
		assert.Contains(t, rec.Body.String(), "version")
	})
}

func TestPost(t *testing.T) {
	cfg := config.New()
	gnf := gnfinder.New(cfg, dictionary, weights)
	gnv := verifier.New(cfg.VerifierURL, []int{})
	text := `
	Thalictroides, 18s per doz.
	vitifoiia, Is. 6d. each
	Calopogon, or Cymbidium pul-

	cheilum, 1 5s. per doz.
	Conostylis americana, 2i. 6d.
	`
	tests := []struct {
		params       api.FinderParams
		bayes        bool
		verif        bool
		bytes        bool
		format       gnfmt.Format
		lang         string
		detectedLang string
		cardinality  []int
	}{
		{api.FinderParams{Text: text},
			true, false, false, gnfmt.CompactJSON, "eng", "", []int{1, 1, 2, 2}},

		{api.FinderParams{Text: text, NoBayes: true},
			false, false, false, gnfmt.CompactJSON, "eng", "", []int{1, 1, 1, 2}},

		{api.FinderParams{Text: text, Verification: true},
			true, true, false, gnfmt.CompactJSON, "eng", "", []int{1, 1, 2, 2}},

		{api.FinderParams{Text: text, BytesOffset: true},
			true, false, true, gnfmt.CompactJSON, "eng", "", []int{1, 1, 2, 2}},

		{api.FinderParams{Text: text, Format: "tsv"},
			true, false, true, gnfmt.TSV, "eng", "", []int{1, 1, 2, 2}},

		{api.FinderParams{Text: text, BytesOffset: true, Language: "deu"},
			true, false, true, gnfmt.CompactJSON, "deu", "", []int{1, 1, 1, 2}},

		{api.FinderParams{Text: text, Language: "detect"},
			true, false, true, gnfmt.CompactJSON, "eng", "eng", []int{1, 1, 2, 2}},
	}

	for i, v := range tests {
		msg := fmt.Sprintf("params %d", i)
		if v.verif && !gnv.IsConnected() {
			log.Print("WARNING: no internet connection, skipping some tests")
		}

		t.Run(msg, func(t *testing.T) {
			reqBody, err := gnfmt.GNjson{}.Encode(v.params)
			assert.Nil(t, err)
			c, rec := initPOST(t, reqBody)
			err = find(gnf)(c)
			assert.Nil(t, err)
			if v.format != gnfmt.TSV {
				var out output.Output
				err = gnfmt.GNjson{}.Decode(rec.Body.Bytes(), &out)
				assert.Nil(t, err)
				cardinalities := make([]int, len(out.Names))
				for i := range out.Names {
					cardinalities[i] = out.Names[i].Cardinality
				}
				assert.Equal(t, cardinalities, v.cardinality)
				assert.Equal(t, out.WithVerification, v.verif)
				assert.Equal(t, out.WithBayes, v.bayes)
				assert.Equal(t, out.Language, v.lang)
			} else if v.format == gnfmt.TSV {
				body := rec.Body.String()
				assert.Contains(t, body, "\tCardinality")
			}
		})
	}
}

func TestPostURL(t *testing.T) {
	cfg := config.New()
	gnf := gnfinder.New(cfg, dictionary, weights)
	gnv := verifier.New(cfg.VerifierURL, []int{})
	url := `https://en.wikipedia.org/wiki/Monochamus_galloprovincialis`
	tests := []struct {
		params api.FinderParams
	}{
		{api.FinderParams{URL: url}},
	}

	for i, v := range tests {
		msg := fmt.Sprintf("params %d", i)
		if !gnv.IsConnected() {
			log.Print("WARNING: no internet connection, skipping URL test")
			return
		}

		t.Run(msg, func(t *testing.T) {
			reqBody, err := gnfmt.GNjson{}.Encode(v.params)
			assert.Nil(t, err)
			c, rec := initPOST(t, reqBody)
			err = find(gnf)(c)
			assert.Nil(t, err)
			var out output.Output
			err = gnfmt.GNjson{}.Decode(rec.Body.Bytes(), &out)
			assert.Nil(t, err)
			assert.Greater(t, len(out.Names), 1)
		})
	}
}

func initPOST(
	t *testing.T,
	recBody []byte,
) (echo.Context, *httptest.ResponseRecorder) {
	r := bytes.NewReader(recBody)
	req := httptest.NewRequest(http.MethodPost, "/", r)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/find")
	return c, rec
}
