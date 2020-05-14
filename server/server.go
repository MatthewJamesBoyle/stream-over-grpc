package server

import (
	"fmt"
	"io"
	"os"

	"github.com/matthewjamesboyle/stream-over-grpc/models/generated/pb"
	"github.com/pkg/errors"
)

type GRPCServer struct {
}

func (G GRPCServer) GetVideoData(stream pb.StreamingService_GetVideoDataServer) error {
	fo, err := os.Create("./output.tar")
	if err != nil {
		return errors.New("failed to create file")
	}
	defer fo.Close()

	for {
		s, err := stream.Recv()
		if err != nil {
			switch errors.Cause(err) {
			case io.EOF:
				return stream.SendAndClose(&pb.GetVideoDataResponse{StatusCode: pb.StatusCode_SUCCESS})
			default:
				return fmt.Errorf("error streaming: %w", err)
			}
		}
		fmt.Println("reading content")
		if _, err := fo.Write(s.Content); err != nil {
			return fmt.Errorf("failed to write to file %w", err)
		}
	}
}
