package main

import (
	"context"
	"encoding/json"
	"fmt"

	"snap_transfer/trxgw_api"

	"github.com/aburavi/snaputils/async"

	"github.com/aburavi/snaputils/proto/transfer"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type TransferServer struct {
	transfer.UnimplementedTransferServer
}

// TransferIntraBank implements transferintraBank.GreeterServer
func (s *TransferServer) PostTransferIntraBank(ctx context.Context, req *transfer.TransferIntraBankRequest) (*transfer.TransferIntraBankResponse, error) {
	var future async.Future
	var err error
	data := transfer.TransferIntraBankResponse{}

	future = async.Exec(func() (interface{}, error) {
		data_transfer_intrabank, err := trxgw_api.TransferIntraBank(ctx, req)
		return data_transfer_intrabank, err
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
		derr := status.Errorf(4001301, er2)
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

// TransferInterBank implements transferinterbank.GreeterServer
func (s *TransferServer) PostTransferInterBank(ctx context.Context, req *transfer.TransferInterBankRequest) (*transfer.TransferInterBankResponse, error) {
	var future async.Future
	var err error
	data := transfer.TransferInterBankResponse{}

	future = async.Exec(func() (interface{}, error) {
		data_transfer_interbank, err := trxgw_api.TransferInterBank(ctx, req)
		return data_transfer_interbank, err
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
		derr := status.Errorf(4001301, er2)
		return nil, derr
	}

	err2 := json.Unmarshal(dval, &data)
	if err2 != nil {
		fmt.Printf("Decode data payload failed: %s\n", err2)
		er3 := fmt.Sprintf(" %s\n", err2)
		derr := status.Errorf(4001301, er3)
		return nil, derr
	}

	ecode := data.ResponseCode[:3]
	_ = grpc.SetHeader(ctx, metadata.Pairs("X-Http-Code", ecode))
	return &data, nil
}

// TransferStatus implements transferstatus.GreeterServer
func (s *TransferServer) PostTransferStatus(ctx context.Context, req *transfer.TransferStatusRequest) (*transfer.TransferStatusResponse, error) {
	var future async.Future
	var err error
	var data *transfer.TransferStatusResponse

	future = async.Exec(func() (interface{}, error) {
		data_transfer_status, err := trxgw_api.TransferStatus(ctx, req)
		return data_transfer_status, err
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
		derr := status.Errorf(4001301, er2)
		return nil, derr
	}

	err2 := json.Unmarshal(dval, &data)
	if err2 != nil {
		fmt.Printf("Decode data payload failed: %s\n", err2)
		er3 := fmt.Sprintf(" %s\n", err2)
		derr := status.Errorf(4001301, er3)
		return nil, derr
	}

	ecode := data.ResponseCode[:3]
	_ = grpc.SetHeader(ctx, metadata.Pairs("X-Http-Code", ecode))
	return data, nil
}
