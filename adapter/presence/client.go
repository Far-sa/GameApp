package presence

import (
	"context"
	"game-app/contract/goproto/presence"
	"game-app/param"
	"game-app/pkg/protobufmapper"
	"game-app/pkg/slice"

	"google.golang.org/grpc"
)

type Config struct{}

type Client struct {
	address string
	// config
}

func New(address string) Client {
	return Client{address: address}
}

func (c Client) GetPresence(ctx context.Context, req param.GetPresenceRequest) (param.GetPresenceResponse, error) {

	// TODO: use richerror
	// TODO: find the best practice for reliable communication
	//! create new connection for every single method call
	conn, err := grpc.Dial(c.address, grpc.WithInsecure())
	if err != nil {
		return param.GetPresenceResponse{}, err
	}
	defer conn.Close()

	client := presence.NewPresenceServiceClient(conn)

	resp, err := client.GetPresence(ctx,
		&presence.GetPresenceRequest{
			UserIds: slice.MapFromUintToUint64(req.UserID),
		})
	if err != nil {
		return param.GetPresenceResponse{}, err
	}

	return protobufmapper.MapGetPresenceResponseFromProtobuf(resp), nil
}
