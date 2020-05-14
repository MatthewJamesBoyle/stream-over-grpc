package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/matthewjamesboyle/stream-over-grpc/models/generated/pb"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

const (
	address   = "localhost:50051"
	sentValue = 1000000 //limit
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("didn't connect %v", err)
	}
	defer conn.Close()
	c := pb.NewStreamingServiceClient(conn)

	buf := make([]byte, 0)
	// open input file
	fi, err := os.Open("./g.tar")
	defer fi.Close()

	if err != nil {
		log.Fatalf("Not able to open file %v", err)
	}

	stat, err := fi.Stat()
	if err != nil {
		log.Fatalf("failed to get file stats: %v ",
			err)
	}

	ctx := context.Background()
	stream, err := c.GetVideoData(ctx)
	if err != nil {
		log.Fatalf(
			"failed to create upload stream for file %v",
			err)
		return
	}
	defer stream.CloseSend()

	buf = make([]byte, stat.Size())
	for {
		// read a chunk
		n, err := fi.Read(buf)
		if err != nil && err != io.EOF {
			err = errors.Wrapf(err,
				"failed to send chunk via stream")
			return
		}
		if n == 0 {
			break
		}
		var i int64
		for i = 0; i < ((stat.Size() / sentValue) * sentValue); i += sentValue {
			err = stream.Send(&pb.GetVideoDataRequest{
				Content: buf[i : i+sentValue],
			})
		}

		//deals with the remainder.
		if stat.Size()%sentValue > 0 {
			err = stream.Send(&pb.GetVideoDataRequest{
				Content: buf[((stat.Size() / sentValue) * sentValue):((stat.Size() / sentValue * sentValue) + (stat.Size() % sentValue))],
			})
		}

		if err != nil {
			err = errors.Wrapf(err,
				"failed to send chunk via stream")
			return
		}
	}

	st, err := stream.CloseAndRecv()
	if err != nil {
		err = errors.Wrapf(err,
			"failed to receive upstream status response")
		return
	}

	if st.StatusCode != pb.StatusCode_SUCCESS {
		err = fmt.Errorf(
			"upload failed: %w",
			err)
		return
	}
	log.Println("finished streaming successfully.")

	return
}
