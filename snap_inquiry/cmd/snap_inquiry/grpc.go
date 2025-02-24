package main

import (
	"context"
	"encoding/json"
	"fmt"

	"snap_inquiry/trxgw_api"

	"github.com/aburavi/snaputils/async"
	"github.com/aburavi/snaputils/proto/inquiry"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type InquiryServer struct {
	inquiry.UnimplementedInquiryServer
}

// BalanceInquiry implements balanceinquiry.GreeterServer
func (s *InquiryServer) PostBalanceInquiry(ctx context.Context, req *inquiry.BalanceInquiryRequest) (*inquiry.BalanceInquiryResponse, error) {
	var future async.Future
	var err error
	data := inquiry.BalanceInquiryResponse{}

	future = async.Exec(func() (interface{}, error) {
		data_balance_inquiry, err := trxgw_api.BalanceInquiry(ctx, req)
		return data_balance_inquiry, err
	})

	val, err := future.Await()
	if err != nil {
		fmt.Printf("service hit failed: %s\n", err)
		return nil, err
	}

	dval, err1 := json.Marshal(val)
	if err1 != nil {
		fmt.Printf("Encode data payload failed: %s\n", err1)
		er2 := fmt.Sprintf(" %s\n", err1)
		derr := status.Errorf(4001101, er2)
		return nil, derr
	}

	err2 := json.Unmarshal(dval, &data)
	if err2 != nil {
		fmt.Printf("Decode data payload failed: %s\n", err2)
		er3 := fmt.Sprintf(" %s\n", err2)
		derr := status.Errorf(4001101, er3)
		return nil, derr
	}

	ecode := data.ResponseCode[:3]
	_ = grpc.SetHeader(ctx, metadata.Pairs("X-Http-Code", ecode))
	return &data, nil

}

// AccountInquiry implements balanceinquiry.GreeterServer
func (s *InquiryServer) PostInternalAccountInquiry(ctx context.Context, req *inquiry.InternalAccountInquiryRequest) (*inquiry.InternalAccountInquiryResponse, error) {
	var future async.Future
	var err error
	data := inquiry.InternalAccountInquiryResponse{}

	future = async.Exec(func() (interface{}, error) {
		data_account_inquiry, err := trxgw_api.InternalAccountInquiry(ctx, req)
		return data_account_inquiry, err
	})

	val, err := future.Await()
	if err != nil {
		fmt.Printf("service hit failed: %s\n", err)
		return nil, err
	}
	dval, err1 := json.Marshal(val)
	if err1 != nil {
		fmt.Printf("Encode data payload failed: %s\n", err1)
		er2 := fmt.Sprintf(" %s\n", err1)
		derr := status.Errorf(4001101, er2)
		return nil, derr
	}

	err2 := json.Unmarshal(dval, &data)
	if err2 != nil {
		fmt.Printf("Decode data payload failed: %s\n", err2)
		er3 := fmt.Sprintf(" %s\n", err2)
		derr := status.Errorf(4001101, er3)
		return nil, derr
	}

	ecode := data.ResponseCode[:3]
	_ = grpc.SetHeader(ctx, metadata.Pairs("X-Http-Code", ecode))
	return &data, nil
}

func (s *InquiryServer) PostExternalAccountInquiry(ctx context.Context, req *inquiry.ExternalAccountInquiryRequest) (*inquiry.ExternalAccountInquiryResponse, error) {
	var future async.Future
	var err error
	data := inquiry.ExternalAccountInquiryResponse{}

	future = async.Exec(func() (interface{}, error) {
		data_account_inquiry, err := trxgw_api.ExternalAccountInquiry(ctx, req)
		return data_account_inquiry, err
	})

	val, err := future.Await()
	if err != nil {
		fmt.Printf("service hit failed: %s\n", err)
		return nil, err
	}
	dval, err1 := json.Marshal(val)
	if err1 != nil {
		fmt.Printf("Encode data payload failed: %s\n", err1)
		er2 := fmt.Sprintf(" %s\n", err1)
		derr := status.Errorf(4001101, er2)
		return nil, derr
	}

	err2 := json.Unmarshal(dval, &data)
	if err2 != nil {
		fmt.Printf("Decode data payload failed: %s\n", err2)
		er3 := fmt.Sprintf(" %s\n", err2)
		derr := status.Errorf(4001101, er3)
		return nil, derr
	}

	ecode := data.ResponseCode[:3]
	_ = grpc.SetHeader(ctx, metadata.Pairs("X-Http-Code", ecode))
	return &data, nil

}
