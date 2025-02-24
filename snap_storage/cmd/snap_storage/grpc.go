package main

import (
	"context"
	//"encoding/json"
	"fmt"
	"snap_storage/storage"
	"time"

	"github.com/aburavi/snaputils/async"

	protostore "github.com/aburavi/snaputils/proto/storage"

	"google.golang.org/grpc/status"
)

type StorageServer struct {
	protostore.UnimplementedStorageServer
}

func (s *StorageServer) PostExternalId(ctx context.Context, req *protostore.ExternalIdRequest) (*protostore.ExternalIdResponse, error) {
	var future async.Future
	var err error
	data := protostore.ExternalIdResponse{}

	future = async.Exec(func() (interface{}, error) {
		rc := storage.NewStorage()
		t := time.Now()
		location := time.FixedZone("UTC+7", 0)
		datenow := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, location)
		push_externalid, err := rc.SetExternalId(ctx, req.ClientId, req.ExternalId, datenow.Unix(), 24*time.Hour)
		return push_externalid, err
	})

	val, err := future.Await()
	if err != nil {
		fmt.Printf("service hit failed: %s\n", err)
		er2 := fmt.Sprintf(" %s\n", err)
		derr := status.Errorf(4001401, er2)
		return nil, derr
	}
	fmt.Printf("SetExternalId: %s\n", val.(string))

	if val.(string) == "insert" {
		data.Status = true
		data.Description = "Allow ExternalId, used it"
	} else {
		data.Status = false
		data.Description = "Conflict ExternalId"
	}

	return &data, nil
}

func (s *StorageServer) PostRefferenceNo(ctx context.Context, req *protostore.ReffNoRequest) (*protostore.ReffNoResponse, error) {
	var future async.Future
	var err error
	data := protostore.ReffNoResponse{}

	future = async.Exec(func() (interface{}, error) {
		rc := storage.NewStorage()
		dtime, err := time.Parse("2006-01-02", req.Datetime[:10])
		if err != nil {
			fmt.Printf("string to time failed: %s\n", err)
			er2 := fmt.Sprintf(" %s\n", err)
			derr := status.Errorf(4001201, er2)
			return nil, derr
		}

		push_reffno, err := rc.SetRefferenceNo(ctx, req.ClientId, req.ReffNo, dtime.Unix(), 720*time.Hour)
		return push_reffno, err
	})

	val, err1 := future.Await()
	if err1 != nil {
		fmt.Printf("service hit failed: %s\n", err)
		er2 := fmt.Sprintf(" %s\n", err1)
		derr := status.Errorf(4001401, er2)
		return nil, derr
	}

	if val.(string) == "success" {
		data.Status = true
		data.Description = "Allow RefferenceNo"
	} else {
		data.Status = false
		data.Description = err.Error()
	}

	return &data, nil
}

func (s *StorageServer) PostTrxId(ctx context.Context, req *protostore.TrxIdRequest) (*protostore.TrxIdResponse, error) {
	var future async.Future
	var err error
	data := protostore.TrxIdResponse{}

	future = async.Exec(func() (interface{}, error) {
		rc := storage.NewStorage()
		push_trxid, err := rc.SetTrxId(ctx, req.ClientId, req.OriginalReffNo, req.ReffNo, 168*time.Hour)
		return push_trxid, err
	})

	val, err1 := future.Await()
	if err1 != nil {
		fmt.Printf("service hit failed: %s\n", err)
		er2 := fmt.Sprintf(" %s\n", err1)
		derr := status.Errorf(4001401, er2)
		return nil, derr
	}

	if val.(string) == "success" {
		data.Status = true
		data.ReffNo = req.ReffNo
	} else {
		data.Status = false
		data.ReffNo = err.Error()
	}

	return &data, nil
}

func (s *StorageServer) GetTrxId(ctx context.Context, req *protostore.TrxIdRequest) (*protostore.TrxIdResponse, error) {
	var future async.Future
	var err error
	data := protostore.TrxIdResponse{}

	future = async.Exec(func() (interface{}, error) {
		rc := storage.NewStorage()
		push_trxid, err := rc.GetTrxId(ctx, req.ClientId, req.OriginalReffNo)
		return push_trxid, err
	})

	val, err1 := future.Await()
	if err1 != nil {
		fmt.Printf("service hit failed: %s\n", err)
		er2 := fmt.Sprintf(" %s\n", err1)
		derr := status.Errorf(4001401, er2)
		return nil, derr
	}

	if val.(string) != "" {
		data.Status = true
		data.ReffNo = val.(string)
	} else {
		data.Status = false
		data.ReffNo = ""
	}

	return &data, nil
}
