package grpcclient

import (
	"context"
	"fmt"
	presenceClient "game-app/adapter/presence"
	"game-app/param"

	"google.golang.org/grpc"
)

func main() {

	conn, err := grpc.Dial(":8086", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := presenceClient.New(conn)

	resp, err := client.GetPresence(context.Background(), param.GetPresenceRequest{
		UserID: []uint{1, 2, 3},
	})
	if err != nil {
		panic(err)
	}

	for _, item := range resp.Items {
		fmt.Println("itme :", item.UserID, item.Timestamp)
	}
}
