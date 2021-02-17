package service

import (
	"context"

	"github.com/playsound/waveform/pb"
)

// WaveformServer ...
type WaveformServer struct {
}

// NewWaveformServer ...
func NewWaveformServer() *WaveformServer {
	return &WaveformServer{}
}

// Waveform ...
func (server *WaveformServer) Waveform(ctx context.Context, input *pb.TrackInput) (output *pb.TrackOutput, err error) {
	duration, waveformURL, err := MP3ToJSON(input.Path)
	if err != nil {
		return
	}
	output = &pb.TrackOutput{
		Duration:    duration,
		WaveformURL: waveformURL,
	}
	return
}
