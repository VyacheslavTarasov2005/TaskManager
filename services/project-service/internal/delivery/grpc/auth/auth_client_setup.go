package auth

import (
	"context"

	auth_pb "project-service/internal/delivery/grpc/auth/auth_pb" // путь к сгенерированным .pb.go

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type UserServiceClient struct {
	conn   *grpc.ClientConn
	client auth_pb.UserServiceClient
}

// NewUserServiceClient создает и возвращает клиента с открытым соединением
func NewUserServiceClient(addr string) (*UserServiceClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := auth_pb.NewUserServiceClient(conn)
	return &UserServiceClient{conn: conn, client: client}, nil
}

func (c *UserServiceClient) Close() error {
	return c.conn.Close()
}

func (c *UserServiceClient) GetClaims(ctx context.Context, accessToken string) (string, error) {
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+accessToken)

	resp, err := c.client.GetClaimsFromToken(ctx, &auth_pb.Empty{})
	if err != nil {
		return "", err
	}

	return resp.UserId, nil
}

func (c *UserServiceClient) GetUserProfile(ctx context.Context, user_id string) error {

	_, err := c.client.GetUserProfile(ctx, &auth_pb.GetUserProfileRequest{
		Id: user_id,
	})
	if err != nil {
		return err
	}

	return nil
}
