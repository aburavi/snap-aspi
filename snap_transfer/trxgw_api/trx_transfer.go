package trxgw_api

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	//"errors"
	"reflect"
	"regexp"
	"strconv"

	"github.com/aburavi/snaputils/storage"

	"github.com/aburavi/snaputils/proto/transfer"

	"github.com/valyala/fasthttp"
)

func TransferIntraBank(ctx context.Context, req *transfer.TransferIntraBankRequest) (*transfer.TransferIntraBankResponse, error) {
	var base = os.Getenv("URL_TRXGTW_BASE")
	var ch = os.Getenv("CHANNELID")
	var uri = base + "/core/api/openapi/btrx"

	data := transfer.TransferIntraBankResponse{}
	grpcheader, _ := getGrpcHeader(ctx)
	// contoh Payload
	// HTTPMethod +”:“+ EndpointUrl +":"+ AccessToken +":“+ Lowercase(HexEncode(SHA- 256(minify(RequestBody))))+ ":“ + TimeStamp
	// POST:/api/v1.0/balance-inquiry:muhpwhwOkPRU9nNXYnyYHj8t54x3:8b4e9e83b5231cff4f84358ec8ca81951cfe9f999f635b1566452a501d5c23b2:2021-11-29T09:22:18.172+07:00

	clientid, terr := GetTokenUserId(grpcheader.Token)
	if terr != nil {
		fmt.Println("Error when parse token: %s", terr)
		data.ResponseCode = "4011701"
		data.ResponseMessage = fmt.Sprintf("Invalid Token (B2B): %s", terr)
		//derr := status.Errorf(4011700, data.ResponseMessage)
		return &data, nil
	}
	if grpcheader.XPartnerid != clientid {
		fmt.Println("Error Different CLientId: %s", terr)
		data.ResponseCode = "4011700"
		data.ResponseMessage = fmt.Sprintf("Unauthorized. clientid not matched")
		//derr := status.Errorf(4011700, data.ResponseMessage)
		return &data, nil
	}
	pth := strings.Split(grpcheader.UrlPath, "/")
	dpth := "/api/v1.0/transfer-" + pth[4]
	dchannel := grpcheader.XChannelid
	extid := grpcheader.XExternalId
	reffid := req.PartnerReferenceNo
	fmt.Println("uripath: " + dpth)

	//Check for Mandatory Field
	var mfield []string
	if req.PartnerReferenceNo == "" {
		fmt.Println("Error when validation mandatory request {partnerReferenceNo}")
		mfield = append(mfield, "partnerReferenceNo")
	}
	if req.Amount == nil || string(reflect.TypeOf(req.Amount).Kind()) == "string" {
		fmt.Println("Error when validation mandatory request {amount}")
		mfield = append(mfield, "amount")
	}
	if req.Amount.Value == "" {
		fmt.Println("Error when validation mandatory request {amount.value}")
		mfield = append(mfield, "amount.value")
	}
	if req.Amount.Currency == "" {
		fmt.Println("Error when validation mandatory request {amount.currency}")
		mfield = append(mfield, "amount.currency")
	}
	if req.BeneficiaryAccountNo == "" {
		fmt.Println("Error when validation mandatory request {beneficiaryAccountNo}")
		mfield = append(mfield, "beneficiaryAccountNo")
	}
	if req.SourceAccountNo == "" {
		fmt.Println("Error when validation mandatory request {sourceAccountNo}")
		mfield = append(mfield, "sourceAccountNo")
	}
	if req.TransactionDate == "" {
		fmt.Println("Error when validation mandatory request {transactionDate}")
		mfield = append(mfield, "transactionDate")
	}

	if len(mfield) != 0 {
		data.ResponseCode = "4001702"
		data.ResponseMessage = fmt.Sprintf("Missing Mandatory Field: {%s}", strings.Join(mfield, ","))
		//derr := status.Errorf(4001702, data.ResponseMessage)
		return &data, nil
	}

	//Check Alphanumeric
	checkString := regexp.MustCompile(`\d+`).MatchString(req.PartnerReferenceNo)
	if req.PartnerReferenceNo != "" && !checkString {
		fmt.Println("Error when validation request")
		data.ResponseCode = "4001700"
		data.ResponseMessage = "Bad Request, Alphanumeric Only"
		//derr := status.Errorf(4001700, data.ResponseMessage)
		return &data, nil
	}

	if dchannel == "" || dchannel != ch {
		fmt.Println("Error when validation channel request")
		data.ResponseCode = "4001700"
		data.ResponseMessage = "Error when validation channel request"
		//derr := status.Errorf(4001700, data.ResponseMessage)
		return &data, nil
	}

	intergrpcrequest, derr := GrpcValidasi(dpth, grpcheader.Token, grpcheader.Signature, clientid, extid, reffid)
	if derr != nil {
		if strings.Contains(derr.Error(), "decryptionerror") {
			fmt.Println("Error when validation request: %s", derr)
			data.ResponseCode = "4011700"
			data.ResponseMessage = fmt.Sprintf("Unauthorized. %s", derr)
			//derr := status.Errorf(4011700, data.ResponseMessage)
			return &data, nil
		} else if strings.Contains(derr.Error(), "Conflict") {
			fmt.Println("Error when validation request: %s", derr)
			data.ResponseCode = "4091700"
			data.ResponseMessage = fmt.Sprintf("Conflict Unique externalid today")
			//derr := status.Errorf(4091700, data.ResponseMessage)
			return &data, nil
		} else {
			fmt.Println("Error when validation request: %s", derr)
			data.ResponseCode = "4001700"
			data.ResponseMessage = fmt.Sprintf("Bad Request. %s", derr)
			//derr := status.Errorf(4001700, data.ResponseMessage)
			return &data, nil
		}
	}

	rek := false
	for _, srcRek := range intergrpcrequest.Srcrekening {
		if req.SourceAccountNo == srcRek {
			rek = true
			break
		} else {
			rek = false
		}
	}
	if !rek {
		fmt.Printf("Account is not Allowed")
		data.ResponseCode = "4041711"
		data.ResponseMessage = "Invalid Account, not Allowed"
		//derr := status.Errorf(4041711, data.ResponseMessage)
		return &data, nil
	}
	sbody, serr := json.Marshal(req)
	if serr != nil {
		fmt.Println("error wrong body format.")
		data.ResponseCode = "5001702"
		data.ResponseMessage = fmt.Sprintf("External Server Error. %s", serr)
		//derr := status.Errorf(5001702, data.ResponseMessage)
		return &data, nil
	}

	ds := strings.Split(intergrpcrequest.Signature, ":")
	//httpMethod := ds[0]
	//endpointUrl := ds[1]
	reqBodyHash := ds[2]
	//timestamp := ds[3]

	bodyGenHash := hashBody(string(sbody))
	fmt.Println("hashbody: " + bodyGenHash)
	if reqBodyHash != bodyGenHash {
		fmt.Println("error body is not matched.")
		data.ResponseCode = "4011700"
		data.ResponseMessage = "Unauthorized. Signature body not matched"
		//derr := status.Errorf(4011700, data.ResponseMessage)
		return &data, nil
	}

	dreq := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(dreq)
	dreq.SetRequestURI(uri)
	dreq.Header.SetMethod("POST")
	dreq.Header.Set("Content-Type", "application/json")
	dresp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(dresp)

	dreq.SetBody(sbody)

	if err := fasthttp.Do(dreq, dresp); err != nil {
		fmt.Printf("Client API POST trx_type:BTRX get failed: %s\n", err)
		data.ResponseCode = "5001702"
		data.ResponseMessage = "External Server Error trx_type BTRX"
		//derr := status.Errorf(5001702, data.ResponseMessage)
		return &data, nil
	}

	b := dresp.Body()
	if dresp.StatusCode() != fasthttp.StatusCreated {
		fmt.Printf("load list data failed code=%d. [err=%v]\n", dresp.StatusCode(), strconv.Itoa(dresp.StatusCode())+", "+string(b))
		data.ResponseCode = "5001702"
		data.ResponseMessage = fmt.Sprintf("External Server Error trx_type AIN")
		//derr := status.Errorf(5001502, data.ResponseMessage)
		return &data, nil
	}

	err2 := json.Unmarshal(b, &data)
	if err2 != nil {
		fmt.Printf("Decode auth data payload failed: %s\n", err2)
		data.ResponseCode = "5001702"
		data.ResponseMessage = fmt.Sprintf("External Server Error. %s", err2)
		//derr := status.Errorf(5001702, data.ResponseMessage)
		return &data, nil
	}

	//ddata := []byte(`{"responseCode":"2001700","responseMessage":"Request has been processed successfully","referenceNo":"2020102977770000000009","partnerReferenceNo":"2020102900000000000001","amount":{"value":"12345678.00","currency":"IDR"},"beneficiaryAccountNo":"888801000003301","currency":"IDR","customerReference":"10052019","sourceAccount":"888801000157508","transactionDate":"2019-07-03T12:08:56-07:00","additionalInfo":{"deviceId":"12345679237","channel":"mobilephone"}}`)
	//drsp, _ := json.Marshal(ddata)
	//        dderr := json.Unmarshal(ddata, &data)
	//        if dderr != nil {
	//                fmt.Printf("Decode auth data payload failed: %s\n", dderr)
	//                data.ResponseCode = "99"
	//                data.ResponseMessage = "Decode auth data payload failed"
	//                return &data, dderr
	//        }
	fmt.Printf("transferIntrabank Response: %s\n", data.ResponseCode)
	dtrxid, err6 := storage.GrpcSetTrxId(clientid, req.PartnerReferenceNo, data.ReferenceNo)
	if err6 != nil {
		fmt.Println(err6.Error())
		//data.ResponseCode = "4091300"
		//data.ResponseMessage = fmt.Sprintf("External Server Error. %s", err6)
		//return &data, nil
	}
	if !dtrxid.Status {
		fmt.Println("Saving TrxId Failed")
		//data.ResponseCode = "4091300"
		//data.ResponseMessage = fmt.Sprintf("External Server Error. %s", err6)
		//return &data, nil
	}
	if (data.ResponseCode == "400") && strings.Contains(data.ResponseMessage, "Tidak Cukup") {
		data.ResponseCode = "4031714"
		data.ResponseMessage = fmt.Sprintf("Insufficient Funds")
		return &data, nil
	}
	data.ResponseCode = "2001700"
	//data.ResponseMessage = "Success"
	if data.AdditionalInfo.DeviceId == "" && data.AdditionalInfo.Channel == "" {
		data.AdditionalInfo = nil
	}
	return &data, nil
}

func TransferInterBank(ctx context.Context, req *transfer.TransferInterBankRequest) (*transfer.TransferInterBankResponse, error) {
	var base = os.Getenv("URL_TRXGTW_BASE")
	var ch = os.Getenv("CHANNELID")
	var uri = base + "/core/api/openapi/btrx"

	data := transfer.TransferInterBankResponse{}
	grpcheader, _ := getGrpcHeader(ctx)
	// contoh Payload
	// HTTPMethod +”:“+ EndpointUrl +":"+ AccessToken +":“+ Lowercase(HexEncode(SHA- 256(minify(RequestBody))))+ ":“ + TimeStamp
	// POST:/api/v1.0/balance-inquiry:muhpwhwOkPRU9nNXYnyYHj8t54x3:8b4e9e83b5231cff4f84358ec8ca81951cfe9f999f635b1566452a501d5c23b2:2021-11-29T09:22:18.172+07:00

	clientid, terr := GetTokenUserId(grpcheader.Token)
	if terr != nil {
		fmt.Println("Error when parse token: %s", terr)
		data.ResponseCode = "4011801"
		data.ResponseMessage = fmt.Sprintf("Invalid Token (B2B): %s", terr)
		//derr := status.Errorf(4011800, data.ResponseMessage)
		return &data, nil
	}
	if grpcheader.XPartnerid != clientid {
		fmt.Println("Error Different CLientId: %s", terr)
		data.ResponseCode = "4011800"
		data.ResponseMessage = fmt.Sprintf("Unauthorized. clientid not matched")
		//derr := status.Errorf(4011800, data.ResponseMessage)
		return &data, nil
	}
	pth := strings.Split(grpcheader.UrlPath, "/")
	dpth := "/api/v1.0/transfer-" + pth[4]
	dchannel := grpcheader.XChannelid
	extid := grpcheader.XExternalId
	reffid := req.PartnerReferenceNo
	fmt.Println("uripath: " + dpth)

	//Check for Mandatory Field
	var mfield []string
	if req.PartnerReferenceNo == "" {
		fmt.Println("Error when validation mandatory request {partnerReferenceNo}")
		mfield = append(mfield, "partnerReferenceNo")
	}
	if req.Amount == nil || string(reflect.TypeOf(req.Amount).Kind()) == "string" {
		fmt.Println("Error when validation mandatory request {amount}")
		mfield = append(mfield, "amount")
	}
	if req.Amount.Value == "" {
		fmt.Println("Error when validation mandatory request {amount.value}")
		mfield = append(mfield, "amount.value")
	}
	if req.Amount.Currency == "" {
		fmt.Println("Error when validation mandatory request {amount.currency}")
		mfield = append(mfield, "amount.currency")
	}
	if req.BeneficiaryAccountName == "" {
		fmt.Println("Error when validation mandatory request {beneficiaryAccountName}")
		mfield = append(mfield, "beneficiaryAccountName")
	}
	if req.BeneficiaryAccountNo == "" {
		fmt.Println("Error when validation mandatory request {beneficiaryAccountNo}")
		mfield = append(mfield, "beneficiaryAccountNo")
	}
	if req.BeneficiaryBankCode == "" {
		fmt.Println("Error when validation mandatory request {beneficiaryBankCode}")
		mfield = append(mfield, "beneficiaryBankCode")
	}
	if req.SourceAccountNo == "" {
		fmt.Println("Error when validation mandatory request {sourceAccountNo}")
		mfield = append(mfield, "sourceAccountNo")
	}
	if req.TransactionDate == "" {
		fmt.Println("Error when validation mandatory request {transactionDate}")
		mfield = append(mfield, "transactionDate")
	}

	if len(mfield) != 0 {
		data.ResponseCode = "4001802"
		data.ResponseMessage = fmt.Sprintf("Missing Mandatory Field: {%s}", strings.Join(mfield, ","))
		//derr := status.Errorf(4001802, data.ResponseMessage)
		return &data, nil
	}

	//Check Alphanumeric
	checkString := regexp.MustCompile(`\d+`).MatchString(req.PartnerReferenceNo)
	if req.PartnerReferenceNo != "" && !checkString {
		fmt.Println("Error when validation request")
		data.ResponseCode = "4001801"
		data.ResponseMessage = "Invalid field format, Alphanumeric Only"
		//derr := status.Errorf(4001800, data.ResponseMessage)
		return &data, nil
	}

	if dchannel == "" || dchannel != ch {
		fmt.Println("Error when validation channel request")
		data.ResponseCode = "4001800"
		data.ResponseMessage = "Error when validation channel request"
		//derr := status.Errorf(4001800, data.ResponseMessage)
		return &data, nil
	}

	intergrpcrequest, derr := GrpcValidasi(dpth, grpcheader.Token, grpcheader.Signature, clientid, extid, reffid)
	if derr != nil {
		if strings.Contains(derr.Error(), "decryptionerror") {
			fmt.Println("Error when validation request: %s", derr)
			data.ResponseCode = "4011800"
			data.ResponseMessage = fmt.Sprintf("Unauthorized. %s", derr)
			//derr := status.Errorf(4011800, data.ResponseMessage)
			return &data, nil
		} else if strings.Contains(derr.Error(), "Conflict") {
			fmt.Println("Error when validation request: %s", derr)
			data.ResponseCode = "4091800"
			data.ResponseMessage = fmt.Sprintf("Conflict Unique externalid today")
			//derr := status.Errorf(4091800, data.ResponseMessage)
			return &data, nil
		} else {
			fmt.Println("Error when validation request: %s", derr)
			data.ResponseCode = "4001800"
			data.ResponseMessage = fmt.Sprintf("Bad Request. %s", derr)
			//derr := status.Errorf(4001800, data.ResponseMessage)
			return &data, nil
		}
	}

	rek := false
	for _, srcRek := range intergrpcrequest.Srcrekening {
		if req.SourceAccountNo == srcRek {
			rek = true
			break
		} else {
			rek = false
		}
	}
	if !rek {
		fmt.Printf("Account is not Allowed")
		data.ResponseCode = "4041811"
		data.ResponseMessage = "Invalid Account, not Allowed"
		//derr := status.Errorf(4041811, data.ResponseMessage)
		return &data, nil
	}

	sbody, serr := json.Marshal(req)
	if serr != nil {
		fmt.Println("error wrong body format.")
		data.ResponseCode = "5001802"
		data.ResponseMessage = fmt.Sprintf("External Server Error. %s", serr)
		//derr := status.Errorf(5001802, data.ResponseMessage)
		return &data, nil
	}

	ds := strings.Split(intergrpcrequest.Signature, ":")
	//httpMethod := ds[0]
	//endpointUrl := ds[1]
	reqBodyHash := ds[2]
	//timestamp := ds[3]

	bodyGenHash := hashBody(string(sbody))
	fmt.Println("body: " + string(sbody))
	fmt.Println("hashbody: " + bodyGenHash)
	if reqBodyHash != bodyGenHash {
		fmt.Println("error body is not matched.")
		data.ResponseCode = "4011800"
		data.ResponseMessage = "Unauthorized. Signature body not matched"
		//derr := status.Errorf(4011800, data.ResponseMessage)
		return &data, nil
	}

	dreq := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(dreq)
	dreq.SetRequestURI(uri)
	dreq.Header.SetMethod("POST")
	dreq.Header.Set("Content-Type", "application/json")
	dresp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(dresp)

	dreq.SetBody(sbody)

	if err := fasthttp.Do(dreq, dresp); err != nil {
		fmt.Printf("Client API POST trx_type:BTRX get failed: %s\n", err)
		data.ResponseCode = "5001802"
		data.ResponseMessage = "External Server Error trx_type BTRX"
		//derr := status.Errorf(5001802, data.ResponseMessage)
		return &data, nil
	}

	b := dresp.Body()
	if dresp.StatusCode() != fasthttp.StatusCreated {
		fmt.Printf("load list data failed code=%d. [err=%v]\n", dresp.StatusCode(), string(b))
		if (dresp.StatusCode() == 400) && strings.Contains(string(b), "Tidak Cukup") {
			data.ResponseCode = "4031814"
			data.ResponseMessage = fmt.Sprintf("Saldo tidak cukup")
		} else {
			data.ResponseCode = "5001802"
			data.ResponseMessage = fmt.Sprintf("External Server Error trx_type AIN")
			//derr := status.Errorf(5001502, data.ResponseMessage)
		}
		return &data, nil
	}

	err2 := json.Unmarshal(b, &data)
	if err2 != nil {
		fmt.Printf("Decode data payload failed: %s\n", err2)
		data.ResponseCode = "5001802"
		data.ResponseMessage = fmt.Sprintf("External Server Error. %s", err2)
		//derr := status.Errorf(5001802, data.ResponseMessage)
		return &data, nil
	}

	//ddata := []byte(`{"responseCode":"2001800","responseMessage":"Request has been processed successfully","referenceNo":"2020102977770000000009","partnerReferenceNo":"2020102900000000000001","amount":{"value":"12345678.00","currency":"IDR"},"beneficiaryAccountNo":"888801000003301","beneficiaryBankCode":"002","sourceAccountNo":"888801000157508","traceNo":"10052019","additionalInfo":{"deviceId":"12345679237","channel":"mobilephone"}}`)
	//drsp, _ := json.Marshal(ddata)
	//        ddderr := json.Unmarshal(ddata, &data)
	//        if ddderr != nil {
	//                fmt.Printf("Decode auth data payload failed: %s\n", ddderr)
	//                data.ResponseCode = "99"
	//                data.ResponseMessage = "Decode auth data payload failed"
	//                return &data, ddderr
	//        }
	dtrxid, err6 := storage.GrpcSetTrxId(clientid, req.PartnerReferenceNo, data.ReferenceNo)
	if err6 != nil {
		fmt.Println(err6.Error())
		//data.ResponseCode = "4091300"
		//data.ResponseMessage = fmt.Sprintf("External Server Error. %s", err6)
		//return &data, nil
	}
	if !dtrxid.Status {
		fmt.Println("Saving TrxId Failed")
		//data.ResponseCode = "4091300"
		//data.ResponseMessage = fmt.Sprintf("External Server Error. %s", err6)
		//return &data, nil
	}
	data.ResponseCode = "2001800"
	//data.ResponseMessage = "Success"
	if data.AdditionalInfo.DeviceId == "" && data.AdditionalInfo.Channel == "" {
		data.AdditionalInfo = nil
	}
	return &data, nil
}

func TransferStatus(ctx context.Context, req *transfer.TransferStatusRequest) (*transfer.TransferStatusResponse, error) {
	var base = os.Getenv("URL_TRXGTW_BASE")
	var ch = os.Getenv("CHANNELID")
	var uri = base + "/core/api/openapi/txa"

	data := transfer.TransferStatusResponse{}
	oriresponse := BodyTransferStatusResponse{}
	grpcheader, _ := getGrpcHeader(ctx)
	// contoh Payload
	// HTTPMethod +”:“+ EndpointUrl +":"+ AccessToken +":“+ Lowercase(HexEncode(SHA- 256(minify(RequestBody))))+ ":“ + TimeStamp
	// POST:/api/v1.0/balance-inquiry:muhpwhwOkPRU9nNXYnyYHj8t54x3:8b4e9e83b5231cff4f84358ec8ca81951cfe9f999f635b1566452a501d5c23b2:2021-11-29T09:22:18.172+07:00

	clientid, terr := GetTokenUserId(grpcheader.Token)
	if terr != nil {
		fmt.Println("Error when parse token: %s", terr)
		data.ResponseCode = "4013601"
		data.ResponseMessage = fmt.Sprintf("Invalid Token (B2B): %s", terr)
		//derr := status.Errorf(4013600, data.ResponseMessage)
		return &data, terr
	}
	if grpcheader.XPartnerid != clientid {
		fmt.Println("Error Different ClientId: %s", terr)
		data.ResponseCode = "4013600"
		data.ResponseMessage = fmt.Sprintf("Unauthorized. clientid not matched")
		//derr := status.Errorf(4013600, data.ResponseMessage)
		return &data, nil
	}
	//Check Alphanumeric
	checkString := regexp.MustCompile(`\d+`).MatchString(req.OriginalPartnerReferenceNo)
	if req.OriginalPartnerReferenceNo != "" && !checkString {
		fmt.Println("Error when validation request")
		data.ResponseCode = "4003601"
		data.ResponseMessage = "Invalid field format, Alphanumeric Only"
		//derr := status.Errorf(4003600, data.ResponseMessage)
		return &data, nil
	}

	pth := strings.Split(grpcheader.UrlPath, "/")
	dpth := "/api/v1.0/transfer-" + pth[4]
	dchannel := grpcheader.XChannelid
	extid := grpcheader.XExternalId
	reffid := "None"
	fmt.Println("uripath: " + dpth)
	//Check for Mandatory Field
	var mfield []string
	if req.ServiceCode == "" {
		fmt.Println("Error when validation mandatory request {serviceCode}")
		mfield = append(mfield, "serviceCode")
	}
	if req.Amount == nil || string(reflect.TypeOf(req.Amount).Kind()) == "string" {
		fmt.Println("Error when validation mandatory request {amount}")
		mfield = append(mfield, "amount")
	}
	if req.Amount.Value == "" {
		fmt.Println("Error when validation mandatory request {amount.value}")
		mfield = append(mfield, "amount.value")
	}
	if req.Amount.Currency == "" {
		fmt.Println("Error when validation mandatory request {amount.currency}")
		mfield = append(mfield, "amount.currency")
	}

	if len(mfield) != 0 {
		data.ResponseCode = "4003602"
		data.ResponseMessage = fmt.Sprintf("Missing Mandatory Field: {%s}", strings.Join(mfield, ","))
		//derr := status.Errorf(4003602, data.ResponseMessage)
		return &data, nil
	}

	if dchannel == "" || dchannel != ch {
		fmt.Println("Error when validation channel request")
		data.ResponseCode = "4003600"
		data.ResponseMessage = "Error when validation channel request"
		//derr := status.Errorf(4003600, data.ResponseMessage)
		return &data, nil
	}

	intergrpcrequest, derr := GrpcValidasi(dpth, grpcheader.Token, grpcheader.Signature, clientid, extid, reffid)
	if derr != nil {
		if strings.Contains(derr.Error(), "decryptionerror") {
			fmt.Println("Error when validation request: %s", derr)
			data.ResponseCode = "4013600"
			data.ResponseMessage = fmt.Sprintf("Unauthorized. %s", derr)
			//derr := status.Errorf(4013600, data.ResponseMessage)
			return &data, nil
		} else if strings.Contains(derr.Error(), "Conflict") {
			fmt.Println("Error when validation request: %s", derr)
			data.ResponseCode = "4093600"
			data.ResponseMessage = fmt.Sprintf("Conflict Unique externalid today")
			//derr := status.Errorf(4093600, data.ResponseMessage)
			return &data, nil
		} else {
			fmt.Println("Error when validation request: %s", derr)
			data.ResponseCode = "4003600"
			data.ResponseMessage = fmt.Sprintf("Bad Request. %s", derr)
			//derr := status.Errorf(4003600, data.ResponseMessage)
			return &data, nil
		}
	}

	//rek := false
	//for _, srcRek := range intergrpcrequest.Srcrekening {
	//        if req.SourceAccountNo == srcRek {
	//                rek = true
	//                break
	//        } else {
	//                rek = false
	//        }
	//}
	//if !rek {
	//        fmt.Printf("Account is not Allowed")
	//        data.ResponseCode = "4041311"
	//        data.ResponseMessage = "Invalid Account, not Allowed"
	//        return &data, errors.New("Invalid Account, not Allowed")
	//}
	sbody, serr := json.Marshal(req)
	if serr != nil {
		fmt.Println("error wrong body format.")
		data.ResponseCode = "5003602"
		data.ResponseMessage = fmt.Sprintf("External Server Error. %s", serr)
		//derr := status.Errorf(5003602, data.ResponseMessage)
		return &data, nil
	}
	//Get TrxId from redis storage

	reffno := ""
	if req.OriginalReferenceNo == "" {
		dtrxid, _ := storage.GrpcGetTrxId(clientid, req.OriginalPartnerReferenceNo)
		reffno = dtrxid.ReffNo
	} else {
		reffno = req.OriginalReferenceNo
	}
	dbody := &BodyTransferStatusRequest{
		TransactionDate: req.TransactionDate,
		ReferenceNumber: reffno,
		Amount:          req.Amount,
	}

	sdbody, err := json.Marshal(dbody)
	if err != nil {
		fmt.Println("error wrong body format.")
		data.ResponseCode = "4003601"
		data.ResponseMessage = "error wrong body format."
		//derr := status.Errorf(4003601, data.ResponseMessage)
		return &data, nil
	}

	ds := strings.Split(intergrpcrequest.Signature, ":")
	//httpMethod := ds[0]
	//endpointUrl := ds[1]
	reqBodyHash := ds[2]
	//timestamp := ds[3]

	fmt.Printf("request transfer status: " + string(sdbody))
	bodyGenHash := hashBody(string(sbody))
	fmt.Println("hashbody: " + bodyGenHash)
	if reqBodyHash != bodyGenHash {
		fmt.Println("error body is not match..")
		data.ResponseCode = "4013600"
		data.ResponseMessage = "Unauthorized. Signature body not matched"
		//derr := status.Errorf(4013600, data.ResponseMessage)
		return &data, nil
	}

	dreq := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(dreq)
	dreq.SetRequestURI(uri)
	dreq.Header.SetMethod("POST")
	dreq.Header.Set("Content-Type", "application/json")
	dresp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(dresp)

	dreq.SetBody(sdbody)

	if err := fasthttp.Do(dreq, dresp); err != nil {
		fmt.Printf("Client API POST trx_type:BTRX get failed: %s\n", err)
		data.ResponseCode = "5003602"
		data.ResponseMessage = "External Server Error trx_type BTRX"
		//derr := status.Errorf(5003602, data.ResponseMessage)
		return &data, nil
	}

	b := dresp.Body()
	if dresp.StatusCode() != fasthttp.StatusCreated {
		fmt.Printf("load list data failed code=%d. [err=%v]\n", dresp.StatusCode(), strconv.Itoa(dresp.StatusCode())+", "+string(b))
		data.ResponseCode = "5003602"
		data.ResponseMessage = fmt.Sprintf("External Server Error. %s", string(b))
		//derr := status.Errorf(5003602, data.ResponseMessage)
		return &data, nil
	}
	fmt.Println("Status Response: %s", string(b))

	err2 := json.Unmarshal(b, &oriresponse)
	if err2 != nil {
		fmt.Printf("Decode data payload failed: %s\n", err2)
		data.ResponseCode = "5003602"
		data.ResponseMessage = fmt.Sprintf("External Server Error. %s", err2)
		//derr := status.Errorf(5003602, data.ResponseMessage)
		return &data, nil
	}

	//ddata := []byte(`{"responseCode":"2002200","responseMessage":"Request has been processed successfully","referenceNo":"2020102977770000000009","partnerReferenceNo":"2020102900000000000001","amount":{"value":"12345678.00","currency":"IDR"},"beneficiaryAccountName":"Yories Yolanda","beneficiaryAccountNo":"888801000003301","beneficiaryAccountType":"D","beneficiaryBankCode":"002","currency":"IDR","customerReference":"10052019","sourceAccountNo":"888801000157508","traceNo":"10052019","transactionDate":"2019-07-03T12:08:56-07:00","transactionStatus":"00","transactionStatusDesc":"success","additionalInfo":{"deviceId":"12345679237","channel":"mobilephone"}}`)
	//drsp, _ := json.Marshal(ddata)
	//        dddderr := json.Unmarshal(ddata, &data)
	//        if dddderr != nil {
	//                fmt.Printf("Decode auth data payload failed: %s\n", dddderr)
	//                data.ResponseCode = "99"
	//                data.ResponseMessage = "Decode auth data payload failed"
	//                return &data, dddderr
	//        }

	if oriresponse.Status == "76" {
		data.ResponseCode = "4043601"
		data.ResponseMessage = oriresponse.Message
	} else {
		data.ResponseCode = oriresponse.Data.ResponseCode
		data.ResponseMessage = oriresponse.Data.ResponseMessage
	}
	data.OriginalReferenceNo = oriresponse.Data.OriginalReferenceNo
	data.OriginalExternalId = oriresponse.Data.OriginalExternalId
	data.ServiceCode = "17"
	data.TransactionDate = oriresponse.Data.TransactionDate
	data.Amount = oriresponse.Data.Amount
	data.BeneficiaryAccountNo = oriresponse.Data.BeneficiaryAccountNo
	data.BeneficiaryBankCode = oriresponse.Data.BeneficiaryBankCode
	data.Currency = oriresponse.Data.Currency
	data.PreviousResponseCode = oriresponse.Data.PreviousResponseCode
	data.ReferenceNumber = oriresponse.Data.ReferenceNumber
	data.SourceAccountNo = oriresponse.Data.SourceAccountNo
	data.TransactionId = oriresponse.Data.TransactionId
	data.LatestTransactionStatus = oriresponse.Data.LatestTransactionStatus
	data.TransactionStatusDesc = oriresponse.Data.TransactionStatusDesc
	data.AdditionalInfo = oriresponse.Data.AdditionalInfo
	//if oriresponse.Data.AdditionalInfo.DeviceId == "" && oriresponse.Data.AdditionalInfo.Channel == "" {
	//        data.AdditionalInfo = nil
	//} else {
	//        data.AdditionalInfo = oriresponse.Data.AdditionalInfo
	//}
	return &data, nil
}
