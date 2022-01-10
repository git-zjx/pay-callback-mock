package handlers

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"net/http"
	"payment-mocker/pkg/param"
	"payment-mocker/pkg/xhttp"
	"strings"

	"github.com/gin-gonic/gin"
)

type WechatRequest struct {
	SignType    string `json:"sign_type" form:"sign_type"`
	PrivateKey  string `json:"private_key" form:"private_key"`
	CallbackUrl string `json:"callback_url" form:"callback_url"`
	ParamsStr   string `json:"params" form:"params"`
	Params      param.Params
}

func WechatHandler(c *gin.Context) {
	var (
		request WechatRequest
		signTmp string
		sign    string
		err     error
	)
	err = c.Bind(&request)
	if err != nil {
		c.String(http.StatusOK, "参数错误")
		return
	}
	err = json.Unmarshal([]byte(request.ParamsStr), &request.Params)
	if err != nil {
		c.String(http.StatusOK, "params 格式为 JSON")
		return
	}
	signData := param.GetQueryString(request.Params) + "&key=" + request.PrivateKey
	switch request.SignType {
	case "MD5":
		md5Init := md5.New()
		_, err = md5Init.Write([]byte(signData))
		if err != nil {
			c.String(http.StatusOK, err.Error())
			return
		}
		signTmp = hex.EncodeToString(md5Init.Sum(nil))
	case "HMAC-SHA256":
		key := []byte(request.PrivateKey)
		hash := hmac.New(sha256.New, key)
		_, err := hash.Write([]byte(signData))
		if err != nil {
			c.String(http.StatusOK, err.Error())
			return
		}
		bytes := hash.Sum(nil)
		signTmp = hex.EncodeToString(bytes)
	default:
		md5Init := md5.New()
		_, err = md5Init.Write([]byte(signData))
		if err != nil {
			c.String(http.StatusOK, err.Error())
			return
		}
		signTmp = hex.EncodeToString(md5Init.Sum(nil))
	}
	sign = strings.ToUpper(signTmp)
	request.Params["sign"] = sign
	xmlStr, err := xml.Marshal(request.Params)
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	res, err := xhttp.Post(request.CallbackUrl, string(xmlStr), "application/xml")
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	c.String(http.StatusOK, res)
}
