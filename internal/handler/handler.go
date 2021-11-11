package handler

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
	"github.com/golang-jwt/jwt"
	"github.com/l12u/userm/internal/errcode"
	"github.com/l12u/userm/internal/model"
	"github.com/l12u/userm/internal/verify"
	"github.com/l12u/userm/pkg/env"
	"k8s.io/klog"
	"net/http"
	"os"
	"time"
)

type RequestHandler struct {
	RsaPrivateKey *rsa.PrivateKey
	Verifier      verify.Verifier
}

func NewRequestHandler() *RequestHandler {
	privKey, err := parseRsaPrivateKeyFromPemStr(os.Getenv("JWT_PRIVATE_KEY"))
	if err != nil {
		panic(err)
	}

	pgVerifier := verify.NewPostgresVerifier(&pg.Options{
		Addr:     env.StringOrDefault("POSTGRES_ADDRESS", "localhost:5432"),
		Database: env.StringOrDefault("POSTGRES_DB", ""),
		User:     env.StringOrDefault("POSTGRES_USER", ""),
		Password: env.StringOrDefault("POSTGRES_PASSWORD", ""),
	})

	return &RequestHandler{
		RsaPrivateKey: privKey,
		Verifier:      pgVerifier,
	}
}

func (h *RequestHandler) Login(c *gin.Context) {
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

	klog.Infof("Try to verify login for '%s'", data.Username)
	ok, err := h.Verifier.Verify(data.Username, data.Password)
	if err != verify.ErrNotFound && err != verify.ErrHashDoesntMatch {
		klog.Errorf("Error during verification process: %v", err)
		errcode.S(c, http.StatusInternalServerError, "error during verification process")
		return
	}
	if !ok {
		errcode.S(c, http.StatusForbidden, "wrong username or password")
		return
	}

	// generate JWT for the user
	issuer := env.StringOrDefault("JWT_ISSUER", "testing@l12u.party")
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"user": data.Username,
		"iss":  issuer,
		"sub":  issuer,
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(3 * 24 * time.Hour).Unix(),
	})

	tokenStr, err := token.SignedString(h.RsaPrivateKey)
	if err != nil {
		_ = c.Error(err)
		return
	}
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
