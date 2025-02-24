package trxgw_api

import (
	"context"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	"github.com/aburavi/snaputils/apisignature"
	"github.com/aburavi/snaputils/keycloakapi"
	"github.com/aburavi/snaputils/ratelimiter"
	"github.com/aburavi/snaputils/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func GetTokenUserId(daccesstoken string) (string, error) {
	var kcpublickey = os.Getenv("KEYCLOAK_PUBLICKEY")
	pubKey := []byte(kcpublickey)
	block, _ := pem.Decode(pubKey)
	if block == nil || block.Type != "PUBLIC KEY" {
		return "", errors.New("ssh: no public key found")
	}
	parsedKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}
	var ok bool
	var pubkey *rsa.PublicKey
	if pubkey, ok = parsedKey.(*rsa.PublicKey); !ok {
		return "", errors.New("unable to parse public key")
	}

	if daccesstoken == "" {
		fmt.Println("could not find bearer token")
		return "", errors.New("Failed to parsing token")
	}
	accesstoken := strings.Split(daccesstoken, " ")
	if len(accesstoken) == 1 {
		fmt.Println("could not find bearer token")
		return "", errors.New("Failed to parsing token")
	}

	token, err := jwt.ParseWithClaims(
		accesstoken[1],
		&ClientIdClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return pubkey, nil
		})
	if err != nil {
		fmt.Printf("parse error: %v\n", err)
		//return "", errors.New("Unknown Token")
	}

	c, ok := token.Claims.(*ClientIdClaims)
	if !ok {
		fmt.Println("fail parsing claim")
		return "", errors.New("Failed to parsing claim")
	}
	fmt.Printf("claims: %+v\n", c.ClientId)

	return c.ClientId, nil
}

func hashBody(body string) string {
	//data := `{"partnerReferenceNo":"2020102900000000000001","bankCardToken":"6d7963617264746f6b656e","accountNo":"7382382957893840","balanceType":["Cash","Coins"],"additionalInfo":{"deviceId":"12345679237", "channel":"mobilephone"}}`
	h := sha256.New()
	h.Write([]byte(body))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func GrpcValidasi(path, token, dsig, clientid, extid, refid string) (*InterGrpcRequest, error) {
	intergrpcrequest := InterGrpcRequest{}
	//1. Storage, set ExternalId
	datenow := time.Now().Format(time.RFC3339)
	dextid, err6 := storage.GrpcExternalId(clientid, extid, datenow)
	if err6 != nil {
		fmt.Println(err6.Error())
		return nil, err6
	}
	if !dextid.Status {
		fmt.Println("ExternalId Failed")
		return nil, status.Errorf(4091300, "Conflict ExternalId")
	}
	//2. Keycloak Check Allowed Path uri
	resourcesetid, err2 := keycloakapi.GrpcCheckUriAccess(path, token)
	if err2 != nil {
		fmt.Println(err2.Error())
		return nil, err2
	}
	if len(resourcesetid) == 0 {
		return nil, errors.New("Sorry, URL is Not Allowed.")
	}

	//3. Keycloak Get Source Rekening
	srcrekening, max, err3 := keycloakapi.GrpcKeycloakResourceAttributes(resourcesetid[0], token)
	if err3 != nil {
		fmt.Println(err3.Error())
		return nil, err3
	}
	if len(srcrekening) == 0 {
		return nil, errors.New("Sorry, Empty Rekening.")
	}

	//4. Backend, get PrivateKey
	signature, err4 := apisignature.RSA_OAEP_Decrypt(dsig, clientid)
	if err4 != nil {
		fmt.Println(err4.Error())
		return nil, err4
	}

	//5. Ratelimiter, push data
	dpush, err5 := ratelimiter.GrpcPushSlidingWindow(clientid, max[0])
	if err5 != nil {
		fmt.Println(err5.Error())
		return nil, err5
	}
	if dpush.Status {
		fmt.Println("API Overlimit")
		return nil, errors.New("Sorry, API sequence is out of limit.")
	}
	//6. Storage, set RefferenceNo
	if refid != "None" {
		dnow := time.Now().Format(time.RFC3339)
		dreffno, err7 := storage.GrpcRefferenceNo(clientid, refid, dnow)
		if err7 != nil {
			fmt.Println(err7.Error())
			return nil, err7
		}
		if !dreffno.Status {
			fmt.Println("RefferenceNo is Failed")
			return nil, status.Errorf(4091301, "Conflict RefferenceNo")
		}
	}

	intergrpcrequest.Resourcesetid = resourcesetid
	intergrpcrequest.Srcrekening = srcrekening
	intergrpcrequest.Signature = signature

	return &intergrpcrequest, nil
}

func GrpcSymetricValidasi(path, token, body, clientid, extid, refid string) (*InterGrpcRequest, error) {
	intergrpcrequest := InterGrpcRequest{}
	//1. Storage, set ExternalId
	datenow := time.Now().Format(time.RFC3339)
	dextid, err6 := storage.GrpcExternalId(clientid, extid, datenow)
	if err6 != nil {
		fmt.Println(err6.Error())
		return nil, err6
	}
	if !dextid.Status {
		fmt.Println("ExternalId Failed")
		return nil, status.Errorf(4091300, "Conflict ExternalId")
	}
	//2. Keycloak Check Allowed Path uri
	resourcesetid, err2 := keycloakapi.GrpcCheckUriAccess(path, token)
	if err2 != nil {
		fmt.Println(err2.Error())
		return nil, err2
	}
	if len(resourcesetid) == 0 {
		return nil, errors.New("Sorry, URL is Not Allowed.")
	}

	//3. Keycloak Get Source Rekening
	srcrekening, max, err3 := keycloakapi.GrpcKeycloakResourceAttributes(resourcesetid[0], token)
	if err3 != nil {
		fmt.Println(err3.Error())
		return nil, err3
	}
	if len(srcrekening) == 0 {
		return nil, errors.New("Sorry, Empty Rekening.")
	}

	//4. Backend, get PrivateKey
	bmsg := []byte(body)
	hasher := sha256.New()
	hasher.Write(bmsg)
	dmsg := strings.ToLower(hex.EncodeToString(hasher.Sum(nil)))
	msg := "POST:" + path + ":" + token + ":" + dmsg
	signature, err4 := apisignature.HMAC512_Encrypt(msg, clientid)
	if err4 != nil {
		fmt.Println(err4.Error())
		return nil, err4
	}

	//5. Ratelimiter, push data
	dpush, err5 := ratelimiter.GrpcPushSlidingWindow(clientid, max[0])
	if err5 != nil {
		fmt.Println(err5.Error())
		return nil, err5
	}
	if dpush.Status {
		fmt.Println("API Overlimit")
		return nil, errors.New("Sorry, API sequence is out of limit.")
	}
	//6. Storage, set RefferenceNo
	if refid != "None" {
		dnow := time.Now().Format(time.RFC3339)
		dreffno, err7 := storage.GrpcRefferenceNo(clientid, refid, dnow)
		if err7 != nil {
			fmt.Println(err7.Error())
			return nil, err7
		}
		if !dreffno.Status {
			fmt.Println("RefferenceNo is Failed")
			return nil, status.Errorf(4091301, "Conflict RefferenceNo")
		}
	}

	intergrpcrequest.Resourcesetid = resourcesetid
	intergrpcrequest.Srcrekening = srcrekening
	intergrpcrequest.Signature = signature

	return &intergrpcrequest, nil
}

// Read metadata from client Header.
func getGrpcHeader(ctx context.Context) (*GrpcHeader, error) {
	data := GrpcHeader{}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.DataLoss, "UnaryEcho: failed to get metadata")
	}
	if t1, ok := md["x-timestamp"]; ok {
		for _, e := range t1 {
			data.XTimestamp = e
		}
	}
	if t2, ok := md["authorization"]; ok {
		for _, e := range t2 {
			data.Token = e
		}
	}
	if t3, ok := md["x-url-method"]; ok {
		for _, e := range t3 {
			data.UrlMethod = e
		}
	}
	if t4, ok := md["x-url-path"]; ok {
		for _, e := range t4 {
			data.UrlPath = e
		}
	}
	if t5, ok := md["x-signature"]; ok {
		for _, e := range t5 {
			data.Signature = e
		}
	}
	if t6, ok := md["x-client-key"]; ok {
		for _, e := range t6 {
			data.ClientKey = e
		}
	}
	if t7, ok := md["origin"]; ok {
		for _, e := range t7 {
			data.Origin = e
		}
	}
	if t8, ok := md["x-partner-id"]; ok {
		for _, e := range t8 {
			data.XPartnerid = e
		}
	}
	if t9, ok := md["x-external-id"]; ok {
		for _, e := range t9 {
			data.XExternalId = e
		}
	}
	if t10, ok := md["channel-id"]; ok {
		for _, e := range t10 {
			data.XChannelid = e
		}
	}

	headerout, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	fmt.Println("header: " + string(headerout))

	return &data, nil
}
