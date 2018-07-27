package boot

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	//log "github.com/sirupsen/logrus"
)

func HmacSha256(secret, message []byte) string {
	hash := hmac.New(sha256.New, secret)
	hash.Write(message)
	return hex.EncodeToString(hash.Sum(nil))
}

func getSign(secret, verb, path string, expires int64, data string) string {
	return HmacSha256([]byte(secret), []byte(fmt.Sprintf("%s%s%d%s", verb, path, expires, data)))
}

func mustMarshal(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}
