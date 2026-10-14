package iologger

import (
	"fmt"

	"github.com/gnames/gn"
)

// CreateLogFileError returns error for log file creation failure.
func CreateLogFileError(path string, err error) *gn.Error {
	return &gn.Error{
		Msg: fmt.Sprintf("Failed to create log file: %s", path),
		Err: err,
	}
}
