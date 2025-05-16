package grpc

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
	"user-service/internal/delivery/grpc/errors"
	"user-service/internal/delivery/grpc/pb"
	"user-service/internal/service/interfaces"
)

type UserServer struct {
	pb.UnimplementedUserServiceServer
	userService interfaces.UserService
	authService interfaces.AuthService
}

func NewUserServer(userService interfaces.UserService, authService interfaces.AuthService) *UserServer {
	return &UserServer{
		userService: userService,
		authService: authService,
	}
}

func (s *UserServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {
	accessToken, refreshToken, err := s.userService.Register(ctx, req.Name, req.Email, req.Password)
	if err != nil {
		return nil, status.Error(errors.ParseError(err), err.Error())
	}

	return &pb.AuthResponse{
		AccessToken:  *accessToken,
		RefreshToken: refreshToken.String(),
	}, nil
}

func (s *UserServer) Login(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	accessToken, refreshToken, err := s.userService.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, status.Error(errors.ParseError(err), err.Error())
	}

	return &pb.AuthResponse{
		AccessToken:  *accessToken,
		RefreshToken: refreshToken.String(),
	}, nil
}

func (s *UserServer) GetClaimsFromToken(ctx context.Context, _ *pb.Empty) (*pb.ClaimsResponse, error) {
	userIDVal := ctx.Value("user_id")
	if userIDVal == nil {
		return nil, status.Error(codes.Unauthenticated, "user ID not found in context")
	}

	return &pb.ClaimsResponse{
		UserId: userIDVal.(string),
	}, nil
}

func (s *UserServer) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.AuthResponse, error) {
	tokenUUID, err := uuid.Parse(req.RefreshToken)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid refresh token format")
	}

	accessToken, refreshToken, err := s.authService.RefreshToken(ctx, tokenUUID)
	if err != nil {
		return nil, status.Error(errors.ParseError(err), err.Error())
	}

	return &pb.AuthResponse{
		AccessToken:  *accessToken,
		RefreshToken: refreshToken.String(),
	}, nil
}

func (s *UserServer) GetMyProfile(ctx context.Context, _ *pb.Empty) (*pb.UserResponse, error) {
	userIDVal := ctx.Value("user_id")
	if userIDVal == nil {
		return nil, status.Error(codes.Unauthenticated, "user ID not found in context")
	}

	userID, err := uuid.Parse(userIDVal.(string))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user ID format")
	}

	user, err := s.userService.GetProfile(ctx, userID)
	if err != nil {
		return nil, status.Error(errors.ParseError(err), err.Error())
	}

	return &pb.UserResponse{
		Id:        user.ID.String(),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: nullableTimeToTimestamppb(user.UpdatedAt),
		IsDeleted: user.IsDeleted,
		Name:      user.Name,
		Email:     user.Email,
	}, nil
}

func (s *UserServer) GetUserProfile(ctx context.Context, req *pb.GetUserProfileRequest) (*pb.UserResponse, error) {
	userID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user ID format")
	}

	user, err := s.userService.GetProfile(ctx, userID)
	if err != nil {
		return nil, status.Error(errors.ParseError(err), err.Error())
	}

	return &pb.UserResponse{
		Id:        user.ID.String(),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: nullableTimeToTimestamppb(user.UpdatedAt),
		IsDeleted: user.IsDeleted,
		Name:      user.Name,
		Email:     user.Email,
	}, nil
}

func (s *UserServer) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UserResponse, error) {
	userIDVal := ctx.Value("user_id")
	if userIDVal == nil {
		return nil, status.Error(codes.Unauthenticated, "user ID not found in context")
	}

	userID, err := uuid.Parse(userIDVal.(string))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user ID format")
	}

	user, err := s.userService.UpdateProfile(ctx, userID, req.Name, req.Email)
	if err != nil {
		return nil, status.Error(errors.ParseError(err), err.Error())
	}

	return &pb.UserResponse{
		Id:        user.ID.String(),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: nullableTimeToTimestamppb(user.UpdatedAt),
		IsDeleted: user.IsDeleted,
		Name:      user.Name,
		Email:     user.Email,
	}, nil
}

func (s *UserServer) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.Empty, error) {
	userIDVal := ctx.Value("user_id")
	if userIDVal == nil {
		return nil, status.Error(codes.Unauthenticated, "user ID not found in context")
	}

	userID, err := uuid.Parse(userIDVal.(string))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user ID format")
	}

	err = s.userService.ChangePassword(ctx, userID, req.OldPassword, req.NewPassword)
	if err != nil {
		return nil, status.Error(errors.ParseError(err), err.Error())
	}

	return &pb.Empty{}, nil
}

func (s *UserServer) DeleteMe(ctx context.Context, _ *pb.Empty) (*pb.Empty, error) {
	userIDVal := ctx.Value("user_id")
	if userIDVal == nil {
		return nil, status.Error(codes.Unauthenticated, "user ID not found in context")
	}

	userID, err := uuid.Parse(userIDVal.(string))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user ID format")
	}

	err = s.userService.DeleteUser(ctx, userID)
	if err != nil {
		return nil, status.Error(errors.ParseError(err), err.Error())
	}

	return &pb.Empty{}, nil
}

func (s *UserServer) RecoverAccount(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	accessToken, refreshToken, err := s.userService.RecoverAccount(ctx, req.Email, req.Password)
	if err != nil {
		return nil, status.Error(errors.ParseError(err), err.Error())
	}

	return &pb.AuthResponse{
		AccessToken:  *accessToken,
		RefreshToken: refreshToken.String(),
	}, nil
}

func nullableTimeToTimestamppb(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	} else {
		return nil
	}
}
