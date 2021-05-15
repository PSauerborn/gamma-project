package utils

import b64 "encoding/base64"

func Base64ToBytes(data string) ([]byte, error) {
	decoded, err := b64.StdEncoding.DecodeString(data)
	return decoded, err
}
