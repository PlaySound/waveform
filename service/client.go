package service

import (
	"context"

	"github.com/playsound/waveform/pb"
	"google.golang.org/grpc"
)

// Client ...
type Client struct {
	service pb.WaveformServiceClient
}

// NewClient ...
func NewClient(cc *grpc.ClientConn) *Client {
	service := pb.NewWaveformServiceClient(cc)
	return &Client{service}
}

// Waveform ...
func (client *Client) Waveform(in *pb.TrackInput) (output *pb.TrackOutput, err error) {
	return client.service.Waveform(context.Background(), in)
}
