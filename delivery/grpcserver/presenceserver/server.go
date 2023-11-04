package presenceserver

import (
	"context"
	"fmt"
	"game-app/contract/golang/presence"
	"game-app/param"
	"game-app/pkg/protobufmapper"
	"game-app/pkg/slice"
	"game-app/service/presenceservice"
	"log"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	presence.UnimplementedPresenceServiceServer
	svc presenceservice.Service
}

func New(svc presenceservice.Service) Server {
	return Server{
		UnimplementedPresenceServiceServer: presence.UnimplementedPresenceServiceServer{},
		svc:                                svc,
	}
}

func (s Server) GetPresence(ctx context.Context, req *presence.GetPresenceRequest) (*presence.GetPresenceResponse, error) {
	resp, err := s.svc.GetPresence(ctx, param.GetPresenceRequest{
		UserID: slice.MapFromUint64ToUint(req.GetUserIds())})
	if err != nil {
		return nil, err
	}

	return protobufmapper.MapGetPresenceResponseToProtobuf(resp), nil
}

func (s Server) Start() {
	// listener
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", 8086))
	if err != nil {
		panic(err)
	}

	// pb PresenceServer
	presenceSvcServer := Server{}

	// grpc Server
	grpcServer := grpc.NewServer()

	// register service to server
	presence.RegisterPresenceServiceServer(grpcServer, &presenceSvcServer)

	// Serve
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal("grpc server error")
	}

}
