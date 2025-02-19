package ParamTools

import (
	"github.com/gin-gonic/gin"
	"reflect"
	"strings"
)

type TokenRes struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

func GetParam(c *gin.Context, key string) (value string, exist bool) {
	value = c.Query(key)
	if value == "" {
		return "", false
	}
	return value, true
}

func ParseTokenBody(buf []byte) TokenRes {
	token := TokenRes{
		AccessToken: "",
		Scope:       "",
		TokenType:   "",
	}
	fType := reflect.TypeOf(token)
	fVal := reflect.New(fType)
	str := string(buf)
	results := strings.Split(str, "&")
	for _, tStr := range results {
		kv := strings.Split(tStr, "=")
		ele := reflect.TypeOf(&token).Elem()
		for i := 0; i < ele.NumField(); i++ {
			if ele.Field(i).Tag.Get("json") == kv[0] {
				fVal.Elem().Field(i).SetString(kv[1])
				break
			}
		}
	}
	return fVal.Elem().Interface().(TokenRes)
}
