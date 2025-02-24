package signature

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aburavi/snaputils/apisignature"
	"github.com/aburavi/snaputils/proto/signature"

	"google.golang.org/grpc/status"
)

func AuthSignature(ctx context.Context, req *signature.AuthSignatureRequest) (*signature.SignatureResponse, error) {
	data := signature.SignatureResponse{}
	//contoh Payload
	//Auth Always asymetric
	//message = "clientid|2022-08-25T00:01:02.000TZD"

	tstamp := req.XTimestamp
	msg := req.ClientId + "|" + tstamp
	signature, err := apisignature.AUTH_RSA_OAEP_Encrypt(msg, req.ClientId)
	if err != nil {
		er2 := fmt.Sprintf("error encrypt signature %s\n", err)
		derr := status.Errorf(4001002, er2)
		return nil, derr
	}

	data.ResponseCode = "2001000"
	data.ResponseMessage = "Success"
	data.Signature = signature

	return &data, nil
}

func TrxSignature(ctx context.Context, req *signature.TrxSignatureRequest) (*signature.SignatureResponse, error) {
	data := signature.SignatureResponse{}
	grpcheader, _ := getGrpcHeader(ctx)
	// contoh Payload
	// asymetric/v1.0 without token
	// HTTPMethod +”:“+ EndpointUrl +":"+ Lowercase(HexEncode(SHA- 256(minify(RequestBody))))+ ":“ + TimeStamp
	// symetric/v2.0 with token
	// HTTPMethod +”:“+ EndpointUrl +":"+ AccessToken +":“+ Lowercase(HexEncode(SHA- 256(minify(RequestBody))))+ ":“ + TimeStamp

	clientid, _ := GetTokenUserId(grpcheader.Token)
	secret, err := apisignature.CekTokenExist(clientid)
	if err != nil {
		er2 := fmt.Sprintf("error validation clientid %s\n", err)
		derr := status.Errorf(4001002, er2)
		return nil, derr
	}
	fmt.Println("request client with secret: " + req.ClientSecret)
	fmt.Println("validate client with secret: " + secret)
	if secret != req.ClientSecret {
		er2 := fmt.Sprintf("error invalid client %s\n", secret)
		derr := status.Errorf(4001002, er2)
		return nil, derr
	}
	method := req.Method
	dpth := req.UrlPath
	tstamp := req.XTimestamp
	accessToken := grpcheader.Token
	bbmsg := &bytes.Buffer{}
	if berr := json.Compact(bbmsg, []byte(req.Body)); berr != nil {
		ber2 := fmt.Sprintf("Unknown json body: %s\n", berr)
		derr := status.Errorf(4001002, ber2)
		return nil, derr
	}
	//bmsg := []byte(bbmsg.Bytes())
	fmt.Println("body: " + string(bbmsg.Bytes()))
	hasher := sha256.New()
	hasher.Write(bbmsg.Bytes())
	msg := strings.ToLower(hex.EncodeToString(hasher.Sum(nil)))
	trxsignature := "" //asymetric: 1.0, symetric:2.0
	if req.Version == "1.0" {
		dmsg := method + ":" + dpth + ":" + msg + ":" + tstamp
		signature, err := apisignature.TRX_RSA_OAEP_Encrypt(dmsg, req.ClientId)
		if err != nil {
			er2 := fmt.Sprintf("error encrypt signature %s\n", err)
			derr := status.Errorf(4001002, er2)
			return nil, derr
		}
		trxsignature = signature
	} else {
		atoken := strings.Split(accessToken, " ")
		dmsg := method + ":" + dpth + ":" + atoken[1] + ":" + msg + ":" + tstamp
		fmt.Println("dmsg ---> " + dmsg)
		signature, err1 := apisignature.HMAC512_Encrypt(dmsg, req.ClientSecret)
		if err1 != nil {
			er2 := fmt.Sprintf("error encrypt signature %s\n", err1)
			derr := status.Errorf(4001002, er2)
			return nil, derr
		}
		trxsignature = signature
	}
	data.ResponseCode = "2001000"
	data.ResponseMessage = "Success"
	data.Signature = trxsignature

	return &data, nil
}
