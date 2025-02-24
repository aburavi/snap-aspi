package main

import (
	"context"
	"encoding/json"
	"fmt"
	"snap_signature/signature"

	"github.com/aburavi/snaputils/async"
	protosig "github.com/aburavi/snaputils/proto/signature"

	"google.golang.org/grpc/status"
)

type SignatureServer struct {
	protosig.UnimplementedSignatureServer
}

func (s *SignatureServer) PostAuthSignature(ctx context.Context, req *protosig.AuthSignatureRequest) (*protosig.SignatureResponse, error) {
	var future async.Future
	var err error
	data := protosig.SignatureResponse{}

	future = async.Exec(func() (interface{}, error) {
		signature, err := signature.AuthSignature(ctx, req)
		return signature, err
	})

	val, err := future.Await()
	if err != nil {
		fmt.Printf("service hit failed: %s\n", err)
		er2 := fmt.Sprintf(" %s\n", err)
		derr := status.Errorf(4001201, er2)
		return nil, derr
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

	return &data, err
}

func (s *SignatureServer) PostTrxSignature(ctx context.Context, req *protosig.TrxSignatureRequest) (*protosig.SignatureResponse, error) {
	var future async.Future
	var err error
	data := protosig.SignatureResponse{}

	future = async.Exec(func() (interface{}, error) {
		signature, err := signature.TrxSignature(ctx, req)
		return signature, err
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

	return &data, err
}
