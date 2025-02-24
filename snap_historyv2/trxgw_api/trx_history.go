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

	"github.com/aburavi/snaputils/proto/history"

	"github.com/valyala/fasthttp"
)

func TransactionHistoryList(ctx context.Context, req *history.TransactionHistoryListRequest) (*history.TransactionHistoryListResponse, error) {
	var base = os.Getenv("URL_TRXGTW_BASE")
	var ch = os.Getenv("CHANNELID")
	var uri = base + "/core/api/openapi/hin"

	data := history.TransactionHistoryListResponse{}
	grpcheader, _ := getGrpcHeader(ctx)
	// contoh Payload
	// HTTPMethod +”:“+ EndpointUrl +":"+ AccessToken +":“+ Lowercase(HexEncode(SHA- 256(minify(RequestBody))))+ ":“ + TimeStamp
	// POST:/api/v2.0/balance-inquiry:muhpwhwOkPRU9nNXYnyYHj8t54x3:8b4e9e83b5231cff4f84358ec8ca81951cfe9f999f635b1566452a501d5c23b2:2021-11-29T09:22:18.172+07:00

	clientid, terr := GetTokenUserId(grpcheader.Token)
	if terr != nil {
		fmt.Println("Error when parse token: %s", terr)
		data.ResponseCode = "4011201"
		data.ResponseMessage = fmt.Sprintf("Invalid Token (B2B): %s", terr)
		//derr := status.Errorf(4011200, data.ResponseMessage)
		return &data, nil
	}
	if grpcheader.XPartnerid != clientid {
		fmt.Println("Error Different CLientId: %s", terr)
		data.ResponseCode = "4011200"
		data.ResponseMessage = fmt.Sprintf("Unauthorized. clientid not matched")
		//derr := status.Errorf(4011200, data.ResponseMessage)
		return &data, nil
	}
	//Check Alphanumeric
	checkString := regexp.MustCompile(`\d+`).MatchString(req.PartnerReferenceNo)
	if req.PartnerReferenceNo != "" && !checkString {
		fmt.Println("Error when validation request")
		data.ResponseCode = "4001201"
		data.ResponseMessage = "Invalid field format, Alphanumeric Only"
		//derr := status.Errorf(4001200, data.ResponseMessage)
		return &data, nil
	}

	pth := strings.Split(grpcheader.UrlPath, "/")
	dpth := "/api/v2.0/" + pth[4]
	dchannel := grpcheader.XChannelid
	extid := grpcheader.XExternalId
	reffid := req.PartnerReferenceNo
	fmt.Println("uripath: " + dpth)
	sbody, serr := json.Marshal(req)
	if serr != nil {
		fmt.Println("error wrong body format.")
		data.ResponseCode = "5001202"
		data.ResponseMessage = fmt.Sprintf("External Server Error. %s", serr)
		//derr := status.Errorf(5001202, data.ResponseMessage)
		return &data, nil
	}
	if dchannel == "" || dchannel != ch {
		fmt.Println("Error when validation channel request")
		data.ResponseCode = "4001200"
		data.ResponseMessage = "Error when validation channel request"
		//derr := status.Errorf(4001200, data.ResponseMessage)
		return &data, nil
	}
	intergrpcrequest, derr := GrpcSymetricValidasi(dpth, grpcheader.Token, string(sbody), clientid, extid, reffid)
	if derr != nil {
		if strings.Contains(derr.Error(), "decryptionerror") {
			fmt.Println("Error when validation request: %s", derr)
			data.ResponseCode = "4011200"
			data.ResponseMessage = fmt.Sprintf("Unauthorized. %s", derr)
			//derr := status.Errorf(4011200, data.ResponseMessage)
			return &data, nil
		} else if strings.Contains(derr.Error(), "Conflict") {
			fmt.Println("Error when validation request: %s", derr)
			data.ResponseCode = "4091200"
			data.ResponseMessage = fmt.Sprintf("Conflict Unique externalid today")
			//derr := status.Errorf(4091200, data.ResponseMessage)
			return &data, nil
		} else {
			fmt.Println("Error when validation request: %s", derr)
			data.ResponseCode = "4001200"
			data.ResponseMessage = fmt.Sprintf("Bad Request. %s", derr)
			//derr := status.Errorf(4001200, data.ResponseMessage)
			return &data, nil
		}
	}

	bodyGenHash := intergrpcrequest.Signature
	fmt.Println("signaturebody: " + grpcheader.Signature)
	if bodyGenHash != grpcheader.Signature {
		fmt.Println("error body is not matched.")
		data.ResponseCode = "4011200"
		data.ResponseMessage = "Unauthorized. Signature body not matched"
		//derr := status.Errorf(4011200, data.ResponseMessage)
		return &data, nil
	}

	dinfo := &Info{
		DeviceId: req.AdditionalInfo.DeviceId,
		Channel:  req.AdditionalInfo.Channel,
	}
	dbody := &bodyHinListRequest{
		AccountNo:          intergrpcrequest.Srcrekening[0],
		PartnerReferenceNo: req.PartnerReferenceNo,
		FromDateTime:       req.FromDateTime,
		ToDateTime:         req.ToDateTime,
		PageSize:           req.PageSize,
		PageNumber:         req.PageNumber,
		AdditionalInfo:     *dinfo,
	}

	sdbody, err := json.Marshal(dbody)
	if err != nil {
		fmt.Println("error wrong body format.")
		data.ResponseCode = "5001202"
		data.ResponseMessage = fmt.Sprintf("External Server Error. %s", err)
		//derr := status.Errorf(5001202, data.ResponseMessage)
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
		fmt.Printf("Client API POST trx_type:HIN get failed: %s\n", err)
		data.ResponseCode = "5001202"
		data.ResponseMessage = "External Server Error trx_type HIN"
		//derr := status.Errorf(5001202, data.ResponseMessage)
		return &data, nil
	}

	b := dresp.Body()
	if dresp.StatusCode() != fasthttp.StatusCreated {
		fmt.Printf("load list data failed code=%d. [err=%v]\n", dresp.StatusCode(), strconv.Itoa(dresp.StatusCode())+", "+string(b))
		data.ResponseCode = "5001202"
		data.ResponseMessage = fmt.Sprintf("External Server Error. %s", string(b))
		//derr := status.Errorf(5001202, data.ResponseMessage)
		return &data, nil
	}

	err2 := json.Unmarshal(b, &data)
	if err2 != nil {
		fmt.Printf("Decode data payload failed: %s\n", err2)
		data.ResponseCode = "5001202"
		data.ResponseMessage = fmt.Sprintf("External Server Error. %s", err2)
		//derr := status.Errorf(5001202, data.ResponseMessage)
		return &data, nil
	}

	//ddata := []byte(`{"responseCode":"2001200","responseMessage":"Request has been processed successfully","referenceNo":"2020102977770000000009","partnerReferenceNo":"2020102900000000000001","detailData":[{"dateTime":"2019-07-03T12:08:56-07:00","amount":{"value":"12345678.00","currency":"IDR"},"remark":"Payment to Warung Ikan Bakar","sourceOfFunds":[{"source":"BALANCE","amount":{"value":"10000.00","currency":"IDR"}}],"status":"SUCCESS","type":"PAYMENT","additionalInfo":{"deviceId":"12345679237","channel":"mobilephone"}}],"additionalInfo":{"deviceId":"12345679237","channel":"mobilephone"}}`)
	//        //drsp, _ := json.Marshal(ddata)
	//        derr := json.Unmarshal(ddata, &data)
	//        if derr != nil {
	//                fmt.Printf("Decode auth data payload failed: %s\n", derr)
	//                data.ResponseCode = "99"
	//                data.ResponseMessage = "Decode auth data payload failed"
	//                return &data, derr
	//        }
	//fmt.Printf("response history list: " + string(b))
	//data.ResponseCode = "2001200"
	//data.ResponseMessage = "Success"
	return &data, nil
}

func TransactionHistoryDetail(ctx context.Context, req *history.TransactionHistoryDetailRequest) (*history.TransactionHistoryDetailResponse, error) {
	var base = os.Getenv("URL_TRXGTW_BASE")
	var ch = os.Getenv("CHANNELID")
	var uri = base + "/core/api/openapi/hin"

	data := history.TransactionHistoryDetailResponse{}
	grpcheader, _ := getGrpcHeader(ctx)
	// contoh Payload
	// HTTPMethod +”:“+ EndpointUrl +":"+ AccessToken +":“+ Lowercase(HexEncode(SHA- 256(minify(RequestBody))))+ ":“ + TimeStamp
	// POST:/api/v2.0/balance-inquiry:muhpwhwOkPRU9nNXYnyYHj8t54x3:8b4e9e83b5231cff4f84358ec8ca81951cfe9f999f635b1566452a501d5c23b2:2021-11-29T09:22:18.172+07:00
	clientid, terr := GetTokenUserId(grpcheader.Token)
	if terr != nil {
		fmt.Println("Error when parse token: %s", terr)
		data.ResponseCode = "4011301"
		data.ResponseMessage = fmt.Sprintf("Invalid Token (B2B): %s", terr)
		//derr := status.Errorf(4011300, data.ResponseMessage)
		return &data, nil
	}
	if grpcheader.XPartnerid != clientid {
		fmt.Println("Error Different CLientId: %s", terr)
		data.ResponseCode = "4011300"
		data.ResponseMessage = fmt.Sprintf("Unauthorized. clientid not matched")
		//derr := status.Errorf(4011300, data.ResponseMessage)
		return &data, nil
	}
	pth := strings.Split(grpcheader.UrlPath, "/")
	dpth := "/api/v2.0/" + pth[4]
	dchannel := grpcheader.XChannelid
	extid := grpcheader.XExternalId
	reffid := "None"
	fmt.Println("uripath: " + dpth)

	//Check for Mandatory Field
	checkString := regexp.MustCompile(`\d+`).MatchString(req.OriginalPartnerReferenceNo)
	if req.OriginalPartnerReferenceNo == "" || !checkString {
		fmt.Println("Error when validation mandatory request")
		data.ResponseCode = "4001302"
		data.ResponseMessage = "Missing Mandatory Field: {originalPartnerReferenceNo}"
		//derr := status.Errorf(4001302, data.ResponseMessage)
		return &data, nil
	}

	sbody, serr := json.Marshal(req)
	if serr != nil {
		fmt.Println("error wrong body format.")
		data.ResponseCode = "5001302"
		data.ResponseMessage = fmt.Sprintf("External Server Error. %s", serr)
		//derr := status.Errorf(5001302, data.ResponseMessage)
		return &data, nil
	}
	if dchannel == "" || dchannel != ch {
		fmt.Println("Error when validation channel request")
		data.ResponseCode = "4001300"
		data.ResponseMessage = "Error when validation channel request"
		//derr := status.Errorf(4001300, data.ResponseMessage)
		return &data, nil
	}
	intergrpcrequest, derr := GrpcSymetricValidasi(dpth, grpcheader.Token, string(sbody), clientid, extid, reffid)
	if derr != nil {
		if strings.Contains(derr.Error(), "decryptionerror") {
			fmt.Println("Error when validation request: %s", derr)
			data.ResponseCode = "4011300"
			data.ResponseMessage = fmt.Sprintf("Unauthorized. %s", derr)
			//derr := status.Errorf(4011300, data.ResponseMessage)
			return &data, nil
		} else if strings.Contains(derr.Error(), "Conflict") {
			fmt.Println("Error when validation request: %s", derr)
			data.ResponseCode = "4091300"
			data.ResponseMessage = fmt.Sprintf("Conflict Unique externalid today")
			//derr := status.Errorf(4091300, data.ResponseMessage)
			return &data, nil
		} else {
			fmt.Println("Error when validation request: %s", derr)
			data.ResponseCode = "4001300"
			data.ResponseMessage = fmt.Sprintf("Bad Request. %s", derr)
			//derr := status.Errorf(4001300, data.ResponseMessage)
			return &data, nil
		}
	}

	bodyGenHash := intergrpcrequest.Signature
	fmt.Println("signaturebody: " + grpcheader.Signature)
	if bodyGenHash != grpcheader.Signature {
		fmt.Println("error body is not matched.")
		data.ResponseCode = "4011300"
		data.ResponseMessage = "Unauthorized. Signature body not matched"
		//derr := status.Errorf(4011300, data.ResponseMessage)
		return &data, nil
	}

	dbody := &bodyHinDetailRequest{
		AccNbr: intergrpcrequest.Srcrekening[0],
		TrxNbr: req.OriginalPartnerReferenceNo,
	}
	sdbody, err := json.Marshal(dbody)
	if err != nil {
		fmt.Println("error wrong body format.")
		data.ResponseCode = "4001301"
		data.ResponseMessage = "error wrong body format."
		//derr := status.Errorf(4001301, data.ResponseMessage)
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

	if err1 := fasthttp.Do(dreq, dresp); err != nil {
		fmt.Printf("Client API POST trx_type:HIN get failed: %s\n", err1)
		data.ResponseCode = "5001302"
		data.ResponseMessage = "External Server Error trx_type HIN"
		//derr := status.Errorf(5001302, data.ResponseMessage)
		return &data, nil
	}

	b := dresp.Body()
	if dresp.StatusCode() != fasthttp.StatusCreated {
		fmt.Printf("load list data failed code=%d. [err=%v]\n", dresp.StatusCode(), strconv.Itoa(dresp.StatusCode())+", "+string(b))
		data.ResponseCode = "5001302"
		data.ResponseMessage = fmt.Sprintf("External Server Error. %s", string(b))
		//derr := status.Errorf(5001302, data.ResponseMessage)
		return &data, nil
	}

	err2 := json.Unmarshal(b, &data)
	if err2 != nil {
		fmt.Printf("Decode data payload failed: %s\n", err2)
		data.ResponseCode = "5001302"
		data.ResponseMessage = fmt.Sprintf("External Server Error. %s", err2)
		//derr := status.Errorf(5001302, data.ResponseMessage)
		return &data, nil
	}

	//ddata := []byte(`{"responseCode":"2001200","responseMessage":"Request has been processed successfully","referenceNo":"2020102977770000000009","partnerReferenceNo":"2020102900000000000001","detailData":[{"dateTime":"2019-07-03T12:08:56-07:00","amount":{"value":"12345678.00","currency":"IDR"},"remark":"Payment to Warung Ikan Bakar","sourceOfFunds":[{"source":"BALANCE","amount":{"value":"10000.00","currency":"IDR"}}],"status":"SUCCESS","type":"PAYMENT","additionalInfo":{"deviceId":"12345679237","channel":"mobilephone"}}],"additionalInfo":{"deviceId":"12345679237","channel":"mobilephone"}}`)
	//drsp, _ := json.Marshal(ddata)
	//        derr := json.Unmarshal(ddata, &data)
	//        if derr != nil {
	//                fmt.Printf("Decode auth data payload failed: %s\n", derr)
	//                data.ResponseCode = "99"
	//                data.ResponseMessage = "Decode auth data payload failed"
	//                return &data, derr

	//data.ResponseCode = "2001200"
	//data.ResponseMessage = "Success"
	return &data, nil
}
