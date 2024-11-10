package infra

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"optionhub-service/internal/config"
)

func UnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "no info in metadata")
	}

	userIDs, ok := md["uuid"]
	if !ok || len(userIDs) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "no uuid in metadata")
	}
	if len(userIDs) > 1 {
		return nil, status.Errorf(codes.Unauthenticated, "multiple uuids are not alowed in metadata")
	}

	ctx = context.WithValue(ctx, config.KeyUUID, userIDs[0])
	return handler(ctx, req)
}
