package trxgw_api

import (
	"github.com/aburavi/snaputils/proto/inquiry"

	jwt "github.com/dgrijalva/jwt-go"
)

type ClientIdClaims struct {
	ClientId string `json:"clientId, omitempty"`
	jwt.StandardClaims
}

type InterGrpcRequest struct {
	Resourcesetid []string `json:"resourceset_id", omitempty`
	Srcrekening   []string `json:"src_rekening", omitempty`
	Signature     string   `json:"signature", omitempty`
}

type GrpcHeader struct {
	XTimestamp  string `json:"xTimestamp", omitempty`
	STimestamp  string `json:"sTimestamp", omitempty`
	Token       string `json:"token", omitempty`
	UrlMethod   string `json:"urlMethod", omitempty`
	UrlPath     string `json:"urlPath", omitempty`
	Signature   string `json:"signature", omitempty`
	ClientKey   string `json:"clientKey", omitempty`
	Origin      string `json:"origin", omitempty`
	XPartnerid  string `json:"xPartnerid", omitempty`
	XExternalId string `json:"xExternalId", omitempty`
	XChannelid  string `json:"xChannelid", omitempty`
}

type ExternalInquiryRequest struct {
	AccountNo            string        `json:"accountNo", omitempty`
	BeneficiaryBankCode  string        `json:"beneficiaryBankCode", omitempty`
	BeneficiaryAccountNo string        `json:"beneficiaryAccountNo", omitempty`
	PartnerReferenceNo   string        `json:"partnerReferenceNo", omitempty`
	AdditionalInfo       *inquiry.Info `json:"additionalInfo", omitempty`
}
