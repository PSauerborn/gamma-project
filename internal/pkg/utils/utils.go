package utils

import (
	"bytes"
	b64 "encoding/base64"
	"io"
	"mime/multipart"
)

func Base64ToBytes(data string) ([]byte, error) {
	decoded, err := b64.StdEncoding.DecodeString(data)
	return decoded, err
}

func BytesToBase64(data []byte) string {
	return b64.StdEncoding.EncodeToString(data)
}

func FileformToBytes(file multipart.File) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
