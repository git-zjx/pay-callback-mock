package handlers

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash"
	"net/http"
	"payment-mocker/pkg/key"
	"payment-mocker/pkg/param"
	"payment-mocker/pkg/xhttp"

	"github.com/gin-gonic/gin"
)

type AlipayRequest struct {
	SignType    string `json:"sign_type" form:"sign_type"`
	PrivateKey  string `json:"private_key" form:"private_key"`
	CallbackUrl string `json:"callback_url" form:"callback_url"`
	ParamsStr   string `json:"params" form:"params"`
	Params      param.Params
}

func AlipayHandler(c *gin.Context) {
	var (
		request        AlipayRequest
		h              hash.Hash
		hashs          crypto.Hash
		encryptedBytes []byte
		priKey         *rsa.PrivateKey
		err            error
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
	signData := param.GetQueryString(request.Params)
	switch request.SignType {
	case "RSA":
		h = sha1.New()
		hashs = crypto.SHA1
	case "RSA2":
		h = sha256.New()
		hashs = crypto.SHA256
	default:
		h = sha256.New()
		hashs = crypto.SHA256
	}
	if _, err := h.Write([]byte(signData)); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	if priKey, err = key.DecodePrivateKey([]byte(key.FormatPrivateKey(request.PrivateKey))); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	encryptedBytes, err = rsa.SignPKCS1v15(rand.Reader, priKey, hashs, h.Sum(nil))
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	sign := base64.StdEncoding.EncodeToString(encryptedBytes)
	request.Params["sign"] = sign
	data := request.Params.FormatURLParam()
	fmt.Println(data)
	res, err := xhttp.Post(request.CallbackUrl, data, "application/x-www-form-urlencoded")
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	c.String(http.StatusOK, res)
}
