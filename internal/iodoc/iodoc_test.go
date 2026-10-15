package iodoc_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gnames/gndoc"
	"github.com/gnames/gnfinder/internal/iodoc"
	"github.com/gnames/gnfinder/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLocalOnly(t *testing.T) {
	dir := t.TempDir()
	html := filepath.Join(dir, "names.html")
	err := os.WriteFile(html, []byte("<p>Pardosa moesta</p>"), 0644)
	require.NoError(t, err)
	img := filepath.Join(dir, "image.png")
	err = os.WriteFile(img, []byte("\x89PNG\r\n\x1a\n\x00\x00\x00\x00"), 0644)
	require.NoError(t, err)

	cfg := config.New(config.OptWithTika(false))
	doc := iodoc.New(cfg, "test")

	res, err := doc.TextFromFile(html, false)
	require.NoError(t, err)
	assert.Contains(t, res.Text, "Pardosa moesta")
	assert.Equal(t, gndoc.ConverterLocal, res.Converter)

	_, err = doc.TextFromFile(img, false)
	assert.ErrorIs(t, err, gndoc.ErrUnsupported, "images need Tika")
}
