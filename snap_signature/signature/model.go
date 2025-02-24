package signature

import (
        jwt "github.com/dgrijalva/jwt-go"
)

type ClientIdClaims struct {
	ClientId string `json:"clientId, omitempty"`
	jwt.StandardClaims
}

type GrpcHeader struct {
        XTimestamp string `json:"xTimestamp", omitempty`
        STimestamp string `json:"sTimestamp", omitempty`
	Token string `json:"token", omitempty`
	UrlMethod string `json:"urlMethod", omitempty`
	UrlPath string `json:"urlPath", omitempty`
	Signature string `json:"signature", omitempty`
        ClientKey string `json:"clientKey", omitempty`
        Origin string `json:"origin", omitempty`
        XPartnerid string `json:"xPartnerid", omitempty`
        XExternalId string `json:"xExternalId", omitempty`
        XChannelid string `json:"xChannelid", omitempty`
}
