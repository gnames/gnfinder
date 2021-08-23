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
	"github.com/labstack/echo"
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
	text := `
	Thalictroides, 18s per doz.
	vitifoiia, Is. 6d. each
	Calopogon, or Cymbidium pul-

	cheilum, 1 5s. per doz.
	Conostylis americana, 2i. 6d.
	`
	params := []struct {
		params      api.FinderParams
		bayes       bool
		verif       bool
		bytes       bool
		cardinality []int
	}{
		{api.FinderParams{Text: text},
			true, false, false, []int{1, 1, 2, 2}},

		{api.FinderParams{Text: text, NoBayes: true},
			false, false, false, []int{1, 1, 1, 2}},

		{api.FinderParams{Text: text, Verification: true},
			true, true, false, []int{1, 1, 2, 2}},

		{api.FinderParams{Text: text, BytesOffset: true},
			true, false, true, []int{1, 1, 2, 2}},
	}

	for i, v := range params {
		msg := fmt.Sprintf("params %d", i)
		if v.verif && !verifier.HasRemote() {
			log.Print("WARNING: no internet connection, skipping some tests")
		}

		t.Run(msg, func(t *testing.T) {
			reqBody, err := gnfmt.GNjson{}.Encode(v.params)
			assert.Nil(t, err)
			c, rec := initPOST(t, reqBody)
			_ = rec
			assert.Nil(t, find(gnf)(c))
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
