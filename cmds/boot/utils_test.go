package boot

import (
	//"fmt"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetSign(t *testing.T) {
	var sign string

	verb := "GET"
	//key := "LAqUlngMIQkIUjXMUreyu3qn"
	secret := "chNOOS4KvNXR_Xq4k4c9qsfoKWvnDecLATCRlcBwyKDYnWgO"
	path := "/api/v1/instrument"
	expires := int64(1518064236)
	data := ""
	sign = getSign(secret, verb, path, expires, data)

	assert.Equal(t, "c7682d435d0cfe87c16098df34ef2eb5a549d4c5a3c2b1f0f77b8af73423bf00", sign)

	//s, _ := json.Marshal("")
	//log.Info(string(s))

}
