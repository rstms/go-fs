package fs

import "io"

const Version = "0.0.2"

// File is a single file within a filesystem.
type File interface {
	io.Reader
	io.Writer
}
