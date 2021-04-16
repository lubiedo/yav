package utils

import (
	"crypto/sha256"
	"mime"
	"os"
	"strings"
)

const defaultmime = "text/plain"

func FileMimeType(path string) (ret string) {
	ret = defaultmime
	pos := strings.LastIndex(path, ".")

	if pos == -1 || len(path[pos:]) == 1 {
		return
	}
	ret = mime.TypeByExtension(path[pos:])
	if ret == "" {
		ret = defaultmime
	}
	return
}

func FileExist(f string) bool {
	_, err := os.Stat(f)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func FileDataChecksum(s []byte) [32]byte { return sha256.Sum256(s) }
func FileChecksum(p string) [32]byte {
	data, err := os.ReadFile(p)
	if err != nil {
		return [32]byte{0}
	}
	sum2 := sha256.Sum256(data)

	return sum2
}

func FileIsDir(p string) bool {
	info, err := os.Stat(p)
	if err != nil {
		return false
	}
	return info.IsDir()
}
