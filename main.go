package main

import (
	"log"
	"net"

	"github.com/playsound/waveform/pb"
	"github.com/playsound/waveform/service"
	"google.golang.org/grpc"
)

func main() {
	waveformServer := service.NewWaveformServer()
	grpcServer := grpc.NewServer()
	pb.RegisterWaveformServiceServer(grpcServer, waveformServer)

	listener, err := net.Listen("tcp", "0.0.0.0:95270")
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
	grpcServer.Serve(listener)
}
