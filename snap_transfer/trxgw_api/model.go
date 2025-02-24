package trxgw_api

import (
	"github.com/aburavi/snaputils/proto/transfer"

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

type BodyTransferStatusRequest struct {
	TransactionDate string           `json:"transactionDate, omitempty"`
	ReferenceNumber string           `json:"referenceNumber, omitempty"`
	Amount          *transfer.Amount `json:"amount, omitempty"`
}

type BodyTransferStatusResponse struct {
	Status  string                           `json:"status, omitempty"`
	Message string                           `json:"message, omitempty"`
	Data    *transfer.TransferStatusResponse `json:"data, omitempty"`
}
