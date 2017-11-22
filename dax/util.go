package dax

import (
	"crypto/md5"
	"encoding/base64"
)

func hashMD5(content []byte) string {
	h := md5.New()
	h.Write(content)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
