package main

import (
	"context"
	"encoding/json"
	"fmt"
	"snap_backend/backend_api"

	"github.com/aburavi/snaputils/async"
	"github.com/aburavi/snaputils/proto/backend"
)

type BackendServer struct {
	backend.UnimplementedBackendServer
}

// Privatekey implements backend.PrivateKeyServer
func (s *BackendServer) GetClientKey(ctx context.Context, req *backend.ClientKeyRequest) (*backend.ClientKeyResponse, error) {
	var future async.Future
	var err error
	data := backend.ClientKeyResponse{}

	future = async.Exec(func() (interface{}, error) {
		data_backend_clientkey, err := backend_api.ClientKey(ctx, req)
		return data_backend_clientkey, err
	})

	val, err := future.Await()
	if err != nil {
		fmt.Printf("service hit failed: %s\n", err)
		return nil, err
	}
	dval, err1 := json.Marshal(val)
	if err1 != nil {
		fmt.Printf("Encode data payload failed: %s\n", err1)
		return nil, err1
	}

	err2 := json.Unmarshal(dval, &data)
	if err2 != nil {
		fmt.Printf("Decode data payload failed: %s\n", err2)
		return nil, err2
	}
	//data.PrivateKey = "00000000000000000000"
	return &data, err

}

// UserKey implements backend.PrivateKeyServer
func (s *BackendServer) GetUserKey(ctx context.Context, req *backend.UserKeyRequest) (*backend.UserKeyResponse, error) {
	var future async.Future
	var err error
	data := backend.UserKeyResponse{}

	future = async.Exec(func() (interface{}, error) {
		data_backend_userkey, err := backend_api.UserKey(ctx, req)
		return data_backend_userkey, err
	})

	val, err := future.Await()
	if err != nil {
		fmt.Printf("service hit failed: %s\n", err)
		return nil, err
	}
	dval, err1 := json.Marshal(val)
	if err1 != nil {
		fmt.Printf("Encode data payload failed: %s\n", err1)
		return nil, err1
	}

	err2 := json.Unmarshal(dval, &data)
	if err2 != nil {
		fmt.Printf("Decode data payload failed: %s\n", err2)
		return nil, err2
	}
	//data.PrivateKey = "00000000000000000000"
	return &data, err

}
