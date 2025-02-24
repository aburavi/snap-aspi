package backend_api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"

	"github.com/aburavi/snaputils/proto/backend"

	"github.com/valyala/fasthttp"
)

func ClientKey(ctx context.Context, req *backend.ClientKeyRequest) (*backend.ClientKeyResponse, error) {
	var base = os.Getenv("URL_BACKEND_BASE")
	var cryptoToken = os.Getenv("CRYPTO_TOKEN")
	var uri = base + "/api/v1/basecrypto/apps/key/" + req.ClientId
	var pp = backend.ClientKeyResponse{}

	type Key struct {
		PrivateKey   string `json:"private_key, omitempty"`
		PublicKey    string `json:"public_key, omitempty"`
		ClientSecret string `json:"client_secret, omitempty"`
	}

	type TokenV1Response struct {
		Data Key `json:"data, omitempty"`
	}
	var p TokenV1Response

	dreq := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(dreq)
	dreq.SetRequestURI(uri)
	dreq.Header.SetMethod("GET")
	dreq.Header.Set("Content-Type", "application/json")
	dreq.Header.Set("Authorization", "Bearer "+cryptoToken)
	dresp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(dresp)
	if err := fasthttp.Do(dreq, dresp); err != nil {
		fmt.Printf("Client get failed: %s\n", err)
		return nil, err
	}
	b := dresp.Body()
	if dresp.StatusCode() != fasthttp.StatusOK {
		fmt.Printf("load list failed code=%d. [err=%v]\n", dresp.StatusCode(), string(b))
		return nil, nil
	}

	err1 := json.Unmarshal(b, &p)
	if err1 != nil {
		fmt.Printf("Decode data payload failed: %s\n", err1)
		return nil, err1
	}
	pp.ResponseMessage = "Success"
	pp.ResponseCode = "00"
	pp.PrivateKey = p.Data.PrivateKey
	pp.PublicKey = p.Data.PublicKey
	pp.ClientSecret = p.Data.ClientSecret
	return &pp, nil
}

func UserKey(ctx context.Context, req *backend.UserKeyRequest) (*backend.UserKeyResponse, error) {
	var base = os.Getenv("URL_BACKEND_BASE")
	var cryptoToken = os.Getenv("CRYPTO_TOKEN")
	var uri = base + "/api/v1/basecrypto/apps/key/" + req.UserId

	type Key struct {
		PrivateKey   string `json:"private_key, omitempty"`
		PublicKey    string `json:"public_key, omitempty"`
		ClientSecret string `json:"client_secret, omitempty"`
	}

	type SecretV1Response struct {
		Data Key `json:"data, omitempty"`
	}

	rspdata := SecretV1Response{}
	protosdata := backend.UserKeyResponse{}

	dreq := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(dreq)
	dreq.SetRequestURI(uri)
	dreq.Header.SetMethod("GET")
	dreq.Header.Set("Content-Type", "application/json")
	dreq.Header.Set("Authorization", "Bearer "+cryptoToken)
	dresp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(dresp)

	reqdata := url.Values{}
	reqdata.Set("user_id", req.UserId)

	if err := fasthttp.Do(dreq, dresp); err != nil {
		fmt.Printf("Client API auth get failed: %s\n", err)
		protosdata.ResponseCode = "99"
		protosdata.ResponseMessage = "Client API auth get failed"
		return &protosdata, err
	}

	b := dresp.Body()
	if dresp.StatusCode() != fasthttp.StatusOK {
		fmt.Printf("load list auth failed code=%d. [err=%v]\n", dresp.StatusCode(), string(b))
		protosdata.ResponseCode = "99"
		protosdata.ResponseMessage = "Load list auh failed"
		return &protosdata, errors.New("code:" + strconv.Itoa(dresp.StatusCode()) + ", " + string(b))
	}

	err1 := json.Unmarshal(b, &rspdata)
	if err1 != nil {
		fmt.Printf("Decode auth data payload failed: %s\n", err1)
		protosdata.ResponseCode = "99"
		protosdata.ResponseMessage = "Decode auth data payload failed"
		return &protosdata, err1
	}

	protosdata.ResponseCode = "00"
	protosdata.ResponseMessage = "Success"
	protosdata.PrivateKey = rspdata.Data.PrivateKey
	protosdata.PublicKey = rspdata.Data.PublicKey
	protosdata.ClientSecret = rspdata.Data.ClientSecret
	return &protosdata, nil
}
