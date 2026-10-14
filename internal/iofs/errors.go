package iofs

import (
	"fmt"

	"github.com/gnames/gn"
)

// CreateDirError returns error for directory creation failure.
func CreateDirError(dir string, err error) *gn.Error {
	return &gn.Error{
		Msg: fmt.Sprintf("Failed to create directory: %s", dir),
		Err: err,
	}
}

// CopyFileError returns error for file copy failure.
func CopyFileError(path string, err error) *gn.Error {
	return &gn.Error{
		Msg: fmt.Sprintf("Failed to write file: %s", path),
		Err: err,
	}
}

// ReadFileError returns error for file read failure.
func ReadFileError(path string, err error) *gn.Error {
	return &gn.Error{
		Msg: fmt.Sprintf("Failed to read file: %s", path),
		Err: err,
	}
}
