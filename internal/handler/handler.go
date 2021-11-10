package handler

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/l12u/userm/internal/errcode"
	"github.com/l12u/userm/internal/model"
	"net/http"
	"os"
	"time"
)

var RsaPrivateKey *rsa.PrivateKey

func Setup() {
	privKey, err := parseRsaPrivateKeyFromPemStr(os.Getenv("JWT_PRIVATE_KEY"))
	if err != nil {
		panic(err)
	}
	RsaPrivateKey = privKey
}

func Login(c *gin.Context) {
	b, err := c.GetRawData()
	if err != nil {
		_ = c.Error(err)
		return
	}
	if b == nil {
		errcode.S(c, http.StatusBadRequest, "need to contain json body with username and password")
		return
	}

	var data model.LoginData
	err = json.Unmarshal(b, &data)
	if err != nil || !data.IsValid() {
		errcode.S(c, http.StatusBadRequest, "need to contain json body with username and password")
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"user": data.Username,
		"iss":  "testing@l12u.party",
		"sub":  "testing@l12u.party",
		"exp":  time.Now().Add(3 * 24 * time.Hour).Unix(),
	})

	tokenStr, err := token.SignedString(RsaPrivateKey)
	c.JSON(200, gin.H{"token": tokenStr})
}

func parseRsaPrivateKeyFromPemStr(privPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return priv, nil
}
