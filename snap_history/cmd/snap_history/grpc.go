package main

import (
	"context"
	"encoding/json"
	"fmt"

	"snap_history/trxgw_api"

	"github.com/aburavi/snaputils/async"
	"github.com/aburavi/snaputils/proto/history"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type HistoryServer struct {
	history.UnimplementedHistoryServer
}

// TransactionHistoryList implements transactionhistorylist.GreeterServer
func (s *HistoryServer) GetTransactionHistoryList(ctx context.Context, req *history.TransactionHistoryListRequest) (*history.TransactionHistoryListResponse, error) {
	var future async.Future
	var err error
	data := history.TransactionHistoryListResponse{}

	future = async.Exec(func() (interface{}, error) {
		data_transaction_history_list, err := trxgw_api.TransactionHistoryList(ctx, req)
		return data_transaction_history_list, err
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
		derr := status.Errorf(4001201, er2)
		return nil, derr
	}

	err2 := json.Unmarshal(dval, &data)
	if err2 != nil {
		fmt.Printf("Decode data payload failed: %s\n", err2)
		er3 := fmt.Sprintf(" %s\n", err2)
		derr := status.Errorf(4001201, er3)
		return nil, derr
	}

	ecode := data.ResponseCode[:3]
	_ = grpc.SetHeader(ctx, metadata.Pairs("X-Http-Code", ecode))
	return &data, nil

}

// TransactionHistoryDetail implements transactionhistorydetail.GreeterServer
func (s *HistoryServer) GetTransactionHistoryDetail(ctx context.Context, req *history.TransactionHistoryDetailRequest) (*history.TransactionHistoryDetailResponse, error) {
	var future async.Future
	var err error
	data := history.TransactionHistoryDetailResponse{}

	future = async.Exec(func() (interface{}, error) {
		data_transaction_history_detail, err := trxgw_api.TransactionHistoryDetail(ctx, req)
		return data_transaction_history_detail, err
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
		derr := status.Errorf(4001201, er2)
		return nil, derr
	}

	err2 := json.Unmarshal(dval, &data)
	if err2 != nil {
		fmt.Printf("Decode data payload failed: %s\n", err2)
		er3 := fmt.Sprintf(" %s\n", err2)
		derr := status.Errorf(4001201, er3)
		return nil, derr
	}

	ecode := data.ResponseCode[:3]
	_ = grpc.SetHeader(ctx, metadata.Pairs("X-Http-Code", ecode))
	return &data, nil

}
