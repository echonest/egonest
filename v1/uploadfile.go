package egonest

// The UploadFile interface allows callers to upload "files" that are not *os.Files, as well as allowing an *os.File.
type UploadFile interface {
	io.Reader
	Name() string
}

// A ReaderWrapper is a convenient wrapper for providing an UploadFile to PostCall when all you have is an io.Reader.
type ReaderWrapper struct {
	FileName string
	io.Reader
}

func (r ReaderWrapper) Name() string {
	return r.FileName
}
