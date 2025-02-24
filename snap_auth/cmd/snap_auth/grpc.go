package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aburavi/snaputils/async"
	"github.com/aburavi/snaputils/keycloakapi"
	"github.com/aburavi/snaputils/proto/authv1"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type AuthV1Server struct {
	authv1.UnimplementedAuthV1Server
}

func (s *AuthV1Server) GetAuthV1(ctx context.Context, req *authv1.AuthV1Request) (*authv1.AuthV1Response, error) {
	var future async.Future
	var err error
	data := authv1.AuthV1Response{}

	future = async.Exec(func() (interface{}, error) {
		data_auth, err := keycloakapi.KeycloakAuthV1(ctx, req)
		return data_auth, err
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

	ecode := data.ResponseCode[:3]
	_ = grpc.SetHeader(ctx, metadata.Pairs("X-Http-Code", ecode))
	return &data, nil
}

func (s *AuthV1Server) GetRefreshAuthV1(ctx context.Context, req *authv1.RefreshAuthV1Request) (*authv1.RefreshAuthV1Response, error) {
	var future async.Future
	var err error
	data := authv1.RefreshAuthV1Response{}

	future = async.Exec(func() (interface{}, error) {
		data_auth, err := keycloakapi.KeycloakRefreshAuthV1(ctx, req)
		return data_auth, err
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

	ecode := data.ResponseCode[:3]
	_ = grpc.SetHeader(ctx, metadata.Pairs("X-Http-Code", ecode))
	return &data, nil
}

func (s *AuthV1Server) GetResourceSetUriV1(ctx context.Context, req *authv1.ResourceSetUriV1Request) (*authv1.ResourceSetUriV1Response, error) {
	var future async.Future
	var err error
	data := authv1.ResourceSetUriV1Response{}

	future = async.Exec(func() (interface{}, error) {
		data_auth, err := keycloakapi.KeycloakCheckUriV1Access(ctx, req)
		return data_auth, err
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
	return &data, err
}

func (s *AuthV1Server) GetResourceSetAttributeV1(ctx context.Context, req *authv1.ResourceSetAttributeV1Request) (*authv1.ResourceSetAttributeV1Response, error) {
	var future async.Future
	var err error
	data := authv1.ResourceSetAttributeV1Response{}

	future = async.Exec(func() (interface{}, error) {
		data_auth, err := keycloakapi.KeycloakCheckAttributeV1Access(ctx, req)
		return data_auth, err
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
	return &data, err
}
