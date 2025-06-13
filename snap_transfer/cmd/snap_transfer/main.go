package main

import (
	"fmt"
	"log"
	"net"

	"github.com/aburavi/snaputils/proto/transfer"

	"github.com/spf13/viper"
	grpc "google.golang.org/grpc"
)

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	//Host := viper.GetString("HOST")
	Port := viper.GetInt("PORT")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		panic(err)
	}

	grpcServer := grpc.NewServer()
	transfer.RegisterTransferServer(grpcServer, &TransferServer{})

	grpcServer.Serve(lis)
}
