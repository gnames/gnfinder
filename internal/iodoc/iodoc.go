// Package iodoc creates GNdoc instances that convert files, web pages and
// uploads into UTF8-encoded texts according to gnfinder settings.
package iodoc

import (
	"log/slog"

	"github.com/gnames/gndoc"
	"github.com/gnames/gnfinder/pkg/config"
)

// New creates GNdoc according to the configuration. By default files are
// converted locally, and Apache Tika is a fallback for formats that cannot
// be converted locally. If WithTika is false, only local conversion is
// used. If WithTikaFirst is true, Tika is used first, and local conversion
// is a fallback. GNdoc is safe for concurrent use.
func New(cfg config.Config, version string) gndoc.GNdoc {
	ua := "gnfinder/" + version + " (https://github.com/gnames/gnfinder)"
	opts := []gndoc.Option{
		gndoc.OptUserAgent(ua),
		gndoc.OptPreferTika(cfg.WithTikaFirst),
	}

	// The cache makes PDF conversion start faster on every run after
	// the first one.
	if dir, err := gndoc.DefaultPDFCacheDir(); err == nil {
		opts = append(opts, gndoc.OptPDFCacheDir(dir))
	} else {
		slog.Warn("Cannot find directory for PDF cache", "error", err)
	}

	var tikaURL string
	if cfg.WithTika {
		tikaURL = cfg.TikaURL
	}
	return gndoc.New(tikaURL, opts...)
}
