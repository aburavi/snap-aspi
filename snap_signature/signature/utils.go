package signature

import (
	"context"
        "crypto/sha256"
	"encoding/json"
	"fmt"
        "strings"
	"os"
        "errors"
        //"time"
        //"crypto/rsa"
        //"crypto/x509"
	//"encoding/pem"

        jwt "github.com/dgrijalva/jwt-go"
        
        "google.golang.org/grpc/metadata"
        "google.golang.org/grpc/status"
        "google.golang.org/grpc/codes"
)

func GetTokenUserId(daccesstoken string) (string, error) {
        var kcpublickey = os.Getenv("KEYCLOAK_PUBLICKEY")
        //var cert *x509.Certificate
        accesstoken := strings.Split(daccesstoken, " ")
        if len(accesstoken) == 1 {
                return "", errors.New("Empty token")
        }

        //block, _ := pem.Decode([]byte(kcpublickey))
        //cert, _ = x509.ParseCertificate(block.Bytes)
        //pub := cert.PublicKey.(*rsa.PublicKey)

        token, err := jwt.ParseWithClaims(
                accesstoken[1],
                &ClientIdClaims{},
                func(token *jwt.Token) (interface{}, error) {
                        return []byte(kcpublickey), nil
        })
        if err != nil {
                fmt.Printf("parse error: %v\n", err)
                //return "", err
        }
    
        c, ok := token.Claims.(*ClientIdClaims)
        if !ok {
                fmt.Println("fail parsing claim")
                return "", errors.New("Can't get Claims")
        }
        
        //now := time.Now()
        //sec := now.Unix()
        //fmt.Println("timestamp exp: " +  string(c.ExpiresAt))
        //fmt.Println("timestamp now: " +  string(sec))
        //if sec > c.ExpiresAt {
                //return "", errors.New("Token is Expired")
        //}

        fmt.Printf("claims: %+v\n", c.ClientId)
        
        return c.ClientId, nil
    }

func hashBody(body string) string {
	//data := `{"partnerReferenceNo":"2020102900000000000001","bankCardToken":"6d7963617264746f6b656e","accountNo":"7382382957893840","balanceType":["Cash","Coins"],"additionalInfo":{"deviceId":"12345679237", "channel":"mobilephone"}}`
	h := sha256.New()
        h.Write([]byte(body))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// Read metadata from client Header.
func getGrpcHeader(ctx context.Context) (*GrpcHeader, error){
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
        if t10, ok := md["x-channel-id"]; ok {
                for _, e := range t10 {
                        data.XChannelid = e
                }
        }

        headerout, err := json.Marshal(data)
        if err != nil {
                panic (err)
        }
        fmt.Println("header: " +  string(headerout))

        return &data, nil
}

