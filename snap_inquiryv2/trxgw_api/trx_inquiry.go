package trxgw_api

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	//"errors"
	"regexp"
	"strconv"

	"github.com/aburavi/snaputils/proto/inquiry"

	//"google.golang.org/grpc/status"
	"github.com/valyala/fasthttp"
)

func BalanceInquiry(ctx context.Context, req *inquiry.BalanceInquiryRequest) (*inquiry.BalanceInquiryResponse, error) {
	var base = os.Getenv("URL_TRXGTW_BASE")
	var ch = os.Getenv("CHANNELID")
	var uri = base + "/core/api/openapi/bin"

	data := inquiry.BalanceInquiryResponse{}
	grpcheader, _ := getGrpcHeader(ctx)
	// contoh Payload
	// HTTPMethod +”:“+ EndpointUrl +":"+ AccessToken +":“+ Lowercase(HexEncode(SHA- 256(minify(RequestBody))))+ ":“ + TimeStamp
	// POST:/api/v2.0/balance-inquiry:muhpwhwOkPRU9nNXYnyYHj8t54x3:8b4e9e83b5231cff4f84358ec8ca81951cfe9f999f635b1566452a501d5c23b2:2021-11-29T09:22:18.172+07:00

	clientid, terr := GetTokenUserId(grpcheader.Token)
	if terr != nil {
		fmt.Println("Error when parse token: %s", terr)
		data.ResponseCode = "4011101"
		data.ResponseMessage = fmt.Sprintf("Invalid Token (B2B): %s", terr)
		//derr := status.Errorf(4011100, data.ResponseMessage)
		return &data, nil
	}
	if grpcheader.XPartnerid != clientid {
		fmt.Println("Error Different CLientId: %s", terr)
		data.ResponseCode = "4011100"
		data.ResponseMessage = fmt.Sprintf("Unauthorized. clientid not matched")
		//derr := status.Errorf(4011100, data.ResponseMessage)
		return &data, nil
	}
	//Check Alphanumeric
	checkString := regexp.MustCompile(`\d+`).MatchString(req.PartnerReferenceNo)
	if req.PartnerReferenceNo != "" && !checkString {
		fmt.Println("Error when validation request")
		data.ResponseCode = "4001101"
		data.ResponseMessage = "Invalid field format, Alphanumeric Only"
		//derr := status.Errorf(4001100, data.ResponseMessage)
		return &data, nil
	}

	pth := strings.Split(grpcheader.UrlPath, "/")
	dpth := "/api/v2.0/" + pth[4]
	dchannel := grpcheader.XChannelid
	extid := grpcheader.XExternalId
	reffid := req.PartnerReferenceNo
	fmt.Println("uripath: " + dpth)

	//Check for Mandatory Field
	if req.AccountNo == "" {
		fmt.Println("Error when validation mandatory request")
		data.ResponseCode = "4001102"
		data.ResponseMessage = "Missing Mandatory Field: {accountNo}"
		//derr := status.Errorf(4001102, data.ResponseMessage)
		return &data, nil
	}

	sbody, serr := json.Marshal(req)
	if serr != nil {
		fmt.Println("error wrong body format.")
		data.ResponseCode = "5001102"
		data.ResponseMessage = fmt.Sprintf("External Server Error. %s", serr)
		//derr := status.Errorf(5001102, data.ResponseMessage)
		return &data, nil
	}
	if dchannel == "" || dchannel != ch {
		fmt.Println("Error when validation channel request")
		data.ResponseCode = "4001100"
		data.ResponseMessage = "Error when validation channel request"
		//derr := status.Errorf(4001100, data.ResponseMessage)
		return &data, nil
	}
	intergrpcrequest, dderr := GrpcSymetricValidasi(dpth, grpcheader.Token, string(sbody), clientid, grpcheader.XTimestamp, extid, reffid)
	if dderr != nil {
		if strings.Contains(dderr.Error(), "decryptionerror") {
			fmt.Println("Error when validation request: %s", dderr)
			data.ResponseCode = "4011100"
			data.ResponseMessage = fmt.Sprintf("Unauthorized. %s", dderr)
			//derr := status.Errorf(4011100, data.ResponseMessage)
			return &data, nil
		} else if strings.Contains(dderr.Error(), "Conflict") {
			fmt.Println("Error when validation request: %s", dderr)
			data.ResponseCode = "4091100"
			data.ResponseMessage = fmt.Sprintf("Conflict Unique externalid today")
			//derr := status.Errorf(4091100, data.ResponseMessage)
			return &data, nil
		} else {
			fmt.Println("Error when validation request: %s", dderr)
			data.ResponseCode = "4001100"
			data.ResponseMessage = fmt.Sprintf("Bad Request. %s", dderr)
			//derr := status.Errorf(4001100, data.ResponseMessage)
			return &data, nil
		}
	}

	rek := false
	for _, srcRek := range intergrpcrequest.Srcrekening {
		if req.AccountNo == srcRek {
			rek = true
			break
		} else {
			rek = false
		}
	}
	if !rek {
		fmt.Printf("Account is not allowed")
		data.ResponseCode = "4041111"
		data.ResponseMessage = "Invalid Account, not Allowed"
		//derr := status.Errorf(4041111, data.ResponseMessage)
		return &data, nil
	}

	bodyGenHash := intergrpcrequest.Signature
	fmt.Println("signaturebody: " + grpcheader.Signature)
	fmt.Println("bodyGenHash: " + bodyGenHash)
	fmt.Println("body: " + string(sbody))
	if bodyGenHash != grpcheader.Signature {
		fmt.Println("error body is not matched.")
		data.ResponseCode = "4011100"
		data.ResponseMessage = "Unauthorized. Signature body not matched"
		//derr := status.Errorf(4011100, data.ResponseMessage)
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
		fmt.Printf("Client API POST trx_type:BIN get failed: %s\n", err)
		data.ResponseCode = "5001102"
		data.ResponseMessage = "External Server Error trx_type BIN"
		//derr := status.Errorf(5001102, data.ResponseMessage)
		return &data, nil
	}

	b := dresp.Body()
	if dresp.StatusCode() != fasthttp.StatusCreated {
		fmt.Printf("load list data failed code=%d. [err=%v]\n", dresp.StatusCode(), strconv.Itoa(dresp.StatusCode())+", "+string(b))
		data.ResponseCode = "5001102"
		data.ResponseMessage = fmt.Sprintf("External Server Error. %s", string(b))
		//derr := status.Errorf(5001102, data.ResponseMessage)
		return &data, nil
	}

	err2 := json.Unmarshal(b, &data)
	if err2 != nil {
		fmt.Println("error wrong body format.")
		data.ResponseCode = "5001102"
		data.ResponseMessage = fmt.Sprintf("External Server Error. %s", err2)
		//derr := status.Errorf(5001102, data.ResponseMessage)
		return &data, nil
	}

	//ddata := []byte(`{"responseCode": "2001100","responseMessage": "Request has been processed successfully","referenceNo": "2020102977770000000009","partnerReferenceNo": "2020102900000000000001","accountNo": "115471119","name": "JONOMADE", "accountInfo": [{"balanceType": "Cash","amount": {"value": "200000.00","currency": "IDR"},"floatAmount": {"value": "50000.00","currency": "IDR"}, "holdAmount": {"value": "20000.00","currency": "IDR"},"availableBalance": {"value": "130000.00","currency": "IDR"}, "ledgerBalance": {"value": "30000.00","currency": "IDR"}, "currentMultilateralLimit": {"value": "10000.00","currency": "IDR"}, "registrationStatusCode": "0001","status": "0001"},{"balanceType": "Coins","amount": {"value": "200000.00","currency": "IDR"},"floatAmount": {"value": "50000.00","currency": "IDR"},"holdAmount": {"value": "20000.00","currency": "IDR"},"availableBalance": {"value": "130000.00","currency": "IDR"},"ledgerBalance": {"value": "30000.00","currency": "IDR"},"currentMultilateralLimit": {"value": "10000.00","currency": "IDR"},"registrationStatusCode": "0001","status": "0001"}],"additionalInfo": {"deviceId": "12345679237","channel": "mobilephone"}}`)
	//drsp, _ := json.Marshal(ddata)
	//        dderr := json.Unmarshal(ddata, &data)
	//        if dderr != nil {
	//                fmt.Printf("Decode auth data payload failed: %s\n", dderr)
	//                data.ResponseCode = "99"
	//                data.ResponseMessage = "Decode auth data payload failed"
	//                return &data, dderr
	//        }

	data.ResponseCode = "2001100"
	//data.ResponseMessage = "Success"
	return &data, nil
}

func InternalAccountInquiry(ctx context.Context, req *inquiry.InternalAccountInquiryRequest) (*inquiry.InternalAccountInquiryResponse, error) {
	var base = os.Getenv("URL_TRXGTW_BASE")
	var ch = os.Getenv("CHANNELID")
	var uri = base + "/core/api/openapi/ain"

	data := inquiry.InternalAccountInquiryResponse{}
	grpcheader, _ := getGrpcHeader(ctx)
	// contoh Payload
	// HTTPMethod +”:“+ EndpointUrl +":"+ AccessToken +":“+ Lowercase(HexEncode(SHA- 256(minify(RequestBody))))+ ":“ + TimeStamp
	// POST:/api/v2.0/balance-inquiry:muhpwhwOkPRU9nNXYnyYHj8t54x3:8b4e9e83b5231cff4f84358ec8ca81951cfe9f999f635b1566452a501d5c23b2:2021-11-29T09:22:18.172+07:00
	sbody, serr := json.Marshal(req)
	if serr != nil {
		fmt.Println("error wrong body format.")
		data.ResponseCode = "5001502"
		data.ResponseMessage = fmt.Sprintf("External Server Error. %s", serr)
		//derr := status.Errorf(5001502, data.ResponseMessage)
		return &data, nil
	}

	clientid, terr := GetTokenUserId(grpcheader.Token)
	if terr != nil {
		fmt.Println("Error when parse token: %s", terr)
		data.ResponseCode = "4011501"
		data.ResponseMessage = fmt.Sprintf("Invalid Token (B2B): %s", terr)
		//derr := status.Errorf(4011500, data.ResponseMessage)
		return &data, nil
	}
	if grpcheader.XPartnerid != clientid {
		fmt.Println("Error Different CLientId: %s", terr)
		data.ResponseCode = "4011500"
		data.ResponseMessage = fmt.Sprintf("Unauthorized. clientid not matched")
		//derr := status.Errorf(4011500, data.ResponseMessage)
		return &data, nil
	}
	//Check Alphanumeric
	checkString := regexp.MustCompile(`\d+`).MatchString(req.PartnerReferenceNo)
	if req.PartnerReferenceNo != "" && !checkString {
		fmt.Println("Error when validation request")
		data.ResponseCode = "4001501"
		data.ResponseMessage = "Invalid field format, Alphanumeric Only"
		//derr := status.Errorf(4001500, data.ResponseMessage)
		return &data, nil
	}

	pth := strings.Split(grpcheader.UrlPath, "/")
	dpth := "/api/v2.0/" + pth[4]
	dchannel := grpcheader.XChannelid
	extid := grpcheader.XExternalId
	reffid := req.PartnerReferenceNo
	fmt.Println("uripath: " + dpth)

	//Check for Mandatory Field
	if req.BeneficiaryAccountNo == "" {
		fmt.Println("Error when validation mandatory request")
		data.ResponseCode = "4001502"
		data.ResponseMessage = "Missing Mandatory Field: {beneficiaryAccountNo}"
		//derr := status.Errorf(4001502, data.ResponseMessage)
		return &data, nil
	}
	if dchannel == "" || dchannel != ch {
		fmt.Println("Error when validation channel request")
		data.ResponseCode = "4001500"
		data.ResponseMessage = "Error when validation channel request"
		//derr := status.Errorf(4001500, data.ResponseMessage)
		return &data, nil
	}
	intergrpcrequest, dderr := GrpcSymetricValidasi(dpth, grpcheader.Token, string(sbody), clientid, grpcheader.XTimestamp, extid, reffid)
	if dderr != nil {
		if strings.Contains(dderr.Error(), "decryptionerror") {
			fmt.Println("Error when validation request: %s", dderr)
			data.ResponseCode = "4011500"
			data.ResponseMessage = fmt.Sprintf("Unauthorized. %s", dderr)
			//derr := status.Errorf(4011500, data.ResponseMessage)
			return &data, nil
		} else if strings.Contains(dderr.Error(), "Conflict") {
			fmt.Println("Error when validation request: %s", dderr)
			data.ResponseCode = "4091500"
			data.ResponseMessage = fmt.Sprintf("Conflict Unique externalid today")
			//derr := status.Errorf(4091500, data.ResponseMessage)
			return &data, nil
		} else {
			fmt.Println("Error when validation request: %s", dderr)
			data.ResponseCode = "4001500"
			data.ResponseMessage = fmt.Sprintf("Bad Request. %s", dderr)
			//derr := status.Errorf(4001500, data.ResponseMessage)
			return &data, nil
		}
	}

	bodyGenHash := intergrpcrequest.Signature
	fmt.Println("signaturebody: " + grpcheader.Signature)
	if bodyGenHash != grpcheader.Signature {
		fmt.Println("error body is not matched.")
		data.ResponseCode = "4011500"
		data.ResponseMessage = "Unauthorized. Signature body not matched"
		//derr := status.Errorf(4011500, data.ResponseMessage)
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
		fmt.Printf("Client API POST trx_type:AIN get failed: %s\n", err)
		data.ResponseCode = "5001502"
		data.ResponseMessage = "External Server Error trx_type AIN"
		//derr := status.Errorf(5001502, data.ResponseMessage)
		return &data, nil
	}

	b := dresp.Body()
	if dresp.StatusCode() != fasthttp.StatusCreated {
		fmt.Printf("load list auth failed code=%d. [err=%v]\n", dresp.StatusCode(), strconv.Itoa(dresp.StatusCode())+", "+string(b))
		if strings.Contains(string(b), "ResponseInquiryAccount.message") {
			data.ResponseCode = "4031518"
			data.ResponseMessage = fmt.Sprintf("Inactive Account")
		} else {
			data.ResponseCode = "5001502"
			data.ResponseMessage = fmt.Sprintf("External Server Error trx_type AINX")
			//derr := status.Errorf(5001502, data.ResponseMessage)
		}
		return &data, nil
	}

	err1 := json.Unmarshal(b, &data)
	if err1 != nil {
		fmt.Printf("Decode data payload failed: %s\n", err1)
		data.ResponseCode = "5001502"
		data.ResponseMessage = fmt.Sprintf("External Server Error. %s", strconv.Itoa(dresp.StatusCode())+", "+string(b))
		//derr := status.Errorf(5001502, data.ResponseMessage)

		return &data, nil
	}

	//ddata := []byte(`{"responseCode": "2001500","responseMessage": "Request has been processed successfully","referenceNo": "2020102977770000000009","partnerReferenceNo": "2020102900000000000001","beneficiaryAccountName": "Yories Yolanda","beneficiaryAccountNo": "888801000157508","beneficiaryAccountStatus": "Rekening aktif","beneficiaryAccountType": "D",
	//"currency": "IDR","additionalInfo": {"deviceId": "12345679237","channel": "mobilephone"}}`)
	//drsp, _ := json.Marshal(ddata)
	//        dderr := json.Unmarshal(ddata, &data)
	//        if dderr != nil {
	//                fmt.Printf("Decode auth data payload failed: %s\n", dderr)
	//                data.ResponseCode = "99"
	//                data.ResponseMessage = "Decode auth data payload failed"
	//                return &data, dderr
	//        }
	if data.BeneficiaryAccountStatus == "not found" {
		data.ResponseCode = "4031518"
		data.BeneficiaryAccountStatus = "Inactive Account"
	} else if data.BeneficiaryAccountStatus == "found" {
		data.ResponseCode = "2001500"
		data.BeneficiaryAccountStatus = "Active Account"
	} else {
		data.ResponseCode = "4041500"
	}
	//data.ResponseCode = "2001500"
	//data.ResponseMessage = "Success"
	return &data, nil
}

func ExternalAccountInquiry(ctx context.Context, req *inquiry.ExternalAccountInquiryRequest) (*inquiry.ExternalAccountInquiryResponse, error) {
	var base = os.Getenv("URL_TRXGTW_BASE")
	var ch = os.Getenv("CHANNELID")
	var uri = base + "/sw/api/openapi/ainx"

	data := inquiry.ExternalAccountInquiryResponse{}
	grpcheader, _ := getGrpcHeader(ctx)
	// contoh Payload
	// HTTPMethod +”:“+ EndpointUrl +":"+ AccessToken +":“+ Lowercase(HexEncode(SHA- 256(minify(RequestBody))))+ ":“ + TimeStamp
	// POST:/api/v2.0/balance-inquiry:muhpwhwOkPRU9nNXYnyYHj8t54x3:8b4e9e83b5231cff4f84358ec8ca81951cfe9f999f635b1566452a501d5c23b2:2021-11-29T09:22:18.172+07:00
	sbody, serr := json.Marshal(req)
	if serr != nil {
		fmt.Println("error wrong body format.")
		data.ResponseCode = "5001602"
		data.ResponseMessage = fmt.Sprintf("External Server Error. %s", serr)
		//derr := status.Errorf(5001602, data.ResponseMessage)
		return &data, nil
	}

	clientid, terr := GetTokenUserId(grpcheader.Token)
	if terr != nil {
		fmt.Println("Error when parse token: %s", terr)
		data.ResponseCode = "4011601"
		data.ResponseMessage = fmt.Sprintf("Invalid Token (B2B): %s", terr)
		//derr := status.Errorf(4011600, data.ResponseMessage)
		return &data, nil
	}
	if grpcheader.XPartnerid != clientid {
		fmt.Println("Error Different CLientId: %s", terr)
		data.ResponseCode = "4011600"
		data.ResponseMessage = fmt.Sprintf("Unauthorized. clientid not matched")
		//derr := status.Errorf(4011600, data.ResponseMessage)
		return &data, nil
	}
	//Check Alphanumeric
	checkString := regexp.MustCompile(`\d+`).MatchString(req.PartnerReferenceNo)
	if req.PartnerReferenceNo != "" && !checkString {
		fmt.Println("Error when validation request")
		data.ResponseCode = "4001601"
		data.ResponseMessage = "Invalid field format, Alphanumeric Only"
		//derr := status.Errorf(4001600, data.ResponseMessage)
		return &data, nil
	}

	pth := strings.Split(grpcheader.UrlPath, "/")
	dpth := "/api/v2.0/" + pth[4]
	dchannel := grpcheader.XChannelid
	extid := grpcheader.XExternalId
	reffid := req.PartnerReferenceNo
	fmt.Println("uripath: " + dpth)

	//Check for Mandatory Field
	var mfield []string
	if req.BeneficiaryBankCode == "" {
		fmt.Println("Error when validation mandatory request {beneficiaryBankCode}")
		mfield = append(mfield, "beneficiaryBankCode")
	}
	if req.BeneficiaryAccountNo == "" {
		fmt.Println("Error when validation mandatory request {beneficiaryAccountNo}")
		mfield = append(mfield, "beneficiaryAccountNo")
	}
	if len(mfield) != 0 {
		data.ResponseCode = "4001602"
		data.ResponseMessage = fmt.Sprintf("Missing Mandatory Field: {%s}", strings.Join(mfield, ","))
		//derr := status.Errorf(4001602, data.ResponseMessage)
		return &data, nil
	}
	if dchannel == "" || dchannel != ch {
		fmt.Println("Error when validation channel request")
		data.ResponseCode = "4001600"
		data.ResponseMessage = "Error when validation channel request"
		//derr := status.Errorf(4001600, data.ResponseMessage)
		return &data, nil
	}
	intergrpcrequest, dderr := GrpcSymetricValidasi(dpth, grpcheader.Token, string(sbody), clientid, grpcheader.XTimestamp, extid, reffid)
	if dderr != nil {
		if strings.Contains(dderr.Error(), "decryptionerror") {
			fmt.Println("Error when validation request: %s", dderr)
			data.ResponseCode = "4011600"
			data.ResponseMessage = fmt.Sprintf("Unauthorized. %s", dderr)
			//derr := status.Errorf(4011600, data.ResponseMessage)
			return &data, nil
		} else if strings.Contains(dderr.Error(), "Conflict") {
			fmt.Println("Error when validation request: %s", dderr)
			data.ResponseCode = "4091600"
			data.ResponseMessage = fmt.Sprintf("Conflict Unique externalid today")
			//derr := status.Errorf(4091600, data.ResponseMessage)
			return &data, nil
		} else {
			fmt.Println("Error when validation request: %s", dderr)
			data.ResponseCode = "4001600"
			data.ResponseMessage = fmt.Sprintf("Bad Request. %s", dderr)
			//derr := status.Errorf(4001600, data.ResponseMessage)
			return &data, nil
		}
	}

	bodyGenHash := intergrpcrequest.Signature
	fmt.Println("signaturebody: " + grpcheader.Signature)
	if bodyGenHash != grpcheader.Signature {
		fmt.Println("error body is not matched.")
		data.ResponseCode = "4011600"
		data.ResponseMessage = "Unauthorized. Signature body not matched"
		//derr := status.Errorf(4001600, data.ResponseMessage)
		return &data, nil
	}

	srcRek := intergrpcrequest.Srcrekening[0]
	datareq := &ExternalInquiryRequest{
		AccountNo:            srcRek,
		BeneficiaryBankCode:  req.BeneficiaryAccountNo,
		BeneficiaryAccountNo: req.BeneficiaryAccountNo,
		PartnerReferenceNo:   req.PartnerReferenceNo,
		AdditionalInfo:       req.AdditionalInfo,
	}
	dtreq, derr1 := json.Marshal(datareq)
	if derr1 != nil {
		fmt.Println("error body can't parsed to json.")
		data.ResponseCode = "4001601"
		data.ResponseMessage = "error body can't parsed to json."
		//derr := status.Errorf(4001601, data.ResponseMessage)
		return &data, nil
	}

	dreq := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(dreq)
	dreq.SetRequestURI(uri)
	dreq.Header.SetMethod("POST")
	dreq.Header.Set("Content-Type", "application/json")
	dresp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(dresp)

	dreq.SetBody(dtreq)

	if err := fasthttp.Do(dreq, dresp); err != nil {
		fmt.Printf("Client API POST trx_type:AINX get failed: %s\n", err)
		data.ResponseCode = "5001602"
		data.ResponseMessage = "External Server Error trx_type AINX"
		//derr := status.Errorf(5001602, data.ResponseMessage)
		return &data, nil
	}

	b := dresp.Body()
	if dresp.StatusCode() != fasthttp.StatusCreated {
		fmt.Printf("load list auth failed code=%d. [err=%v]\n", dresp.StatusCode(), strconv.Itoa(dresp.StatusCode())+", "+string(b))
		if strings.Contains(string(b), "ResponseInquiryAccount.message") {
			data.ResponseCode = "4031618"
			data.ResponseMessage = fmt.Sprintf("Inactive Account")
		} else {
			data.ResponseCode = "5001602"
			data.ResponseMessage = fmt.Sprintf("External Server Error trx_type AINX")
			//derr := status.Errorf(5001502, data.ResponseMessage)
		}
	}

	err1 := json.Unmarshal(b, &data)
	if err1 != nil {
		fmt.Printf("Decode data payload failed: %s\n", err1)
		data.ResponseCode = "5001602"
		data.ResponseMessage = fmt.Sprintf("External Server Error. %s", strconv.Itoa(dresp.StatusCode())+", "+string(b))
		//derr := status.Errorf(5001502, data.ResponseMessage)

		return &data, nil
	}

	//ddata := []byte(`{"responseCode": "2001500","responseMessage": "Request has been processed successfully","referenceNo": "2020102977770000000009","partnerReferenceNo": "2020102900000000000001","beneficiaryAccountName": "Yories Yolanda","beneficiaryAccountNo": "888801000157508","beneficiaryAccountStatus": "Rekening aktif","beneficiaryAccountType": "D",
	//"currency": "IDR","additionalInfo": {"deviceId": "12345679237","channel": "mobilephone"}}`)
	//drsp, _ := json.Marshal(ddata)
	//        dderr := json.Unmarshal(ddata, &data)
	//        if dderr != nil {
	//                fmt.Printf("Decode auth data payload failed: %s\n", dderr)
	//                data.ResponseCode = "99"
	//                data.ResponseMessage = "Decode auth data payload failed"
	//                return &data, dderr
	//        }

	data.ResponseCode = "2001600"
	//data.ResponseMessage = "Success"
	return &data, nil
}
