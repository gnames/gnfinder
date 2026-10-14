package cmd

import (
	"fmt"

	"github.com/gnames/gn"
)

func configValueError(path string, err error) *gn.Error {
	return &gn.Error{
		Msg: fmt.Sprintf("Invalid value in config file or environment: %s", path),
		Err: err,
	}
}

func readInputError(input string, err error) *gn.Error {
	return &gn.Error{
		Msg: fmt.Sprintf("Cannot get text from <em>%s</em>", input),
		Err: err,
	}
}

func loadDataError(name string, err error) *gn.Error {
	return &gn.Error{
		Msg: fmt.Sprintf("Cannot load %s", name),
		Err: err,
	}
}

func webServiceError(port int, err error) *gn.Error {
	return &gn.Error{
		Msg: fmt.Sprintf("Web service on port %d stopped", port),
		Err: err,
	}
}
