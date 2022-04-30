package file

import "os"

func FileNotExists(path string) bool {
	_, err := os.Stat(path)
	return os.IsNotExist(err)
}

type ReaderError struct {
	Message string
}

func (cre *ReaderError) Error() string {
	return cre.Message
}

