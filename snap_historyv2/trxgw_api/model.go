package trxgw_api

import (
        jwt "github.com/dgrijalva/jwt-go"
)

type ClientIdClaims struct {
	ClientId string `json:"clientId, omitempty"`
	jwt.StandardClaims
}

type InterGrpcRequest struct {
        Resourcesetid   []string `json:"resourceset_id", omitempty`
        Srcrekening     []string `json:"src_rekening", omitempty`
        Signature       string `json:"signature", omitempty`
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

type Info struct {
        DeviceId string `json:"deviceId, omitempty"`
        Channel string `json:"channel, omitempty"`
}

type bodyHinListRequest struct {
	AccountNo string `json:"accountNo, omitempty"`
        PartnerReferenceNo string `json:"partnerReferenceNo, omitempty"`
	FromDateTime string `json:"fromDateTime, omitempty"`
        ToDateTime string `json:"toDateTime, omitempty"`
        PageSize string `json:"pageSize, omitempty"`
        PageNumber string `json:"pageNumber, omitempty"`
        AdditionalInfo Info `json:"additionalInfo, omitempty"`
}

type bodyHinDetailRequest struct {
	AccNbr string `json:"accNbr, omitempty"`
	TrxNbr string `json:"trxNbr, omitempty"`
}