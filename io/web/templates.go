package web

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"path"

	"github.com/gnames/gnlib/ent/verifier"
	"github.com/labstack/echo/v4"
)

//go:embed templates
var tmpls embed.FS

// echoTempl implements echo.Renderer interface.
type echoTempl struct {
	templates *template.Template
}

// Render implements echo.Renderer interface.
func (t *echoTempl) Render(
	w io.Writer,
	name string,
	data interface{},
	c echo.Context,
) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewTemplate() (*echoTempl, error) {
	t, err := parseFiles()
	if err != nil {
		return nil, fmt.Errorf("cannot parse file %w", err)
	}
	return &echoTempl{t}, nil
}

func parseFiles() (*template.Template, error) {
	var err error
	var t *template.Template

	var filenames []string
	dir := "templates"
	entries, _ := tmpls.ReadDir(dir)
	for i := range entries {
		if entries[i].Type().IsRegular() {
			filenames = append(
				filenames,
				fmt.Sprintf("%s/%s", dir, entries[i].Name()),
			)
		}
	}

	for _, filename := range filenames {
		name := path.Base(filename)
		var tmpl *template.Template
		if t == nil {
			t = template.New(name)
		}
		if name == t.Name() {
			tmpl = t
		} else {
			tmpl = t.New(name)
		}
		addFuncs(tmpl)
		_, err = tmpl.ParseFS(tmpls, filename)
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}

func addFuncs(tmpl *template.Template) {
	tmpl.Funcs(template.FuncMap{
		"isVerified": func(ver *verifier.Verification) template.HTML {
			if ver == nil || ver.MatchType == verifier.NoMatch {
				return template.HTML("")
			}
			switch ver.MatchType {
			case verifier.Exact:
				return template.HTML("<span class='exact-match'>✔ (" + ver.BestResult.MatchedCanonicalFull + ")</span>")
			default:
				return template.HTML("<span class='some-match'>✔ (" + ver.BestResult.MatchedCanonicalFull + ")</span>")
			}
		},
	})
}
