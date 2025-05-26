package grpc

import (
	"context"
	"project-service/internal/delivery/grpc/errors"
	"project-service/internal/delivery/grpc/pb"
	serviceErrors "project-service/internal/service/errors"
	"project-service/internal/service/interfaces"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ProjectServer struct {
	pb.UnimplementedProjectServiceServer
	projectService interfaces.ProjectService
}

func NewProjectServer(projectService interfaces.ProjectService) *ProjectServer {
	return &ProjectServer{
		projectService: projectService,
	}
}

func (s *ProjectServer) CreateProject(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	userIDVal := ctx.Value("user_id")
	if userIDVal == nil {
		return nil, errors.ParseError(serviceErrors.ApplicationError{
			StatusCode: 401,
			Code:       "Unauthorized",
			Errors:     map[string]string{"message": "User ID not found in context"},
		})
	}
	userID, err := uuid.Parse(userIDVal.(string))
	if err != nil {
		return nil, errors.ParseError(serviceErrors.ApplicationError{
			StatusCode: 401,
			Code:       "Unauthorized",
			Errors:     map[string]string{"message": "Invalid user ID format"},
		})
	}

	projectId, err := s.projectService.Create(ctx, req.GetName(), userID)
	if err != nil {
		return nil, err
	}

	return &pb.CreateResponse{
		ProjectId: projectId.String(),
	}, nil
}

func (s *ProjectServer) GetProject(ctx context.Context, req *pb.GetProjectRequest) (*pb.Project, error) {
	projectID, err := uuid.Parse(req.ProjectId)
	if err != nil {
		return nil, errors.ParseError(serviceErrors.ApplicationError{
			StatusCode: 401,
			Code:       "Unauthorized",
			Errors:     map[string]string{"message": "Invalid project ID format"},
		})
	}
	project, err := s.projectService.GetProject(ctx, projectID)
	if err != nil {
		return nil, errors.ParseError(err)
	}

	return &pb.Project{
		Id:        project.ID.String(),
		Name:      project.Name,
		OwnerId:   project.OwnerID.String(),
		CreatedAt: timestamppb.New(project.CreatedAt),
		UpdatedAt: nullableTimeToTimestamppb(&project.UpdatedAt),
	}, nil
}

func (s *ProjectServer) GetProjectsByUser(ctx context.Context, req *pb.GetProjectsByUserRequest) (*pb.GetProjectsByUserResponse, error) {
	userIDVal := ctx.Value("user_id")
	if userIDVal == nil {
		return nil, errors.ParseError(serviceErrors.ApplicationError{
			StatusCode: 401,
			Code:       "Unauthorized",
			Errors:     map[string]string{"message": "User ID not found in context"},
		})
	}
	userID, err := uuid.Parse(userIDVal.(string))
	if err != nil {
		return nil, errors.ParseError(serviceErrors.ApplicationError{
			StatusCode: 401,
			Code:       "Unauthorized",
			Errors:     map[string]string{"message": "Invalid user ID format"},
		})
	}

	projects, err := s.projectService.GetProjectsByUser(ctx, userID)
	if err != nil {
		return nil, errors.ParseError(err)
	}

	projectsResponse := make([]*pb.Project, 0)

	for _, project := range projects {
		projectsResponse = append(projectsResponse, &pb.Project{
			Id:        project.ID.String(),
			Name:      project.Name,
			OwnerId:   project.OwnerID.String(),
			CreatedAt: timestamppb.New(project.CreatedAt),
			UpdatedAt: nullableTimeToTimestamppb(&project.UpdatedAt),
		})
	}
	return &pb.GetProjectsByUserResponse{
		Projects: projectsResponse,
	}, nil
}

func (s *ProjectServer) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.Project, error) {
	userIDVal := ctx.Value("user_id")
	if userIDVal == nil {
		return nil, errors.ParseError(serviceErrors.ApplicationError{
			StatusCode: 401,
			Code:       "Unauthorized",
			Errors:     map[string]string{"message": "User ID not found in context"},
		})
	}
	userID, err := uuid.Parse(userIDVal.(string))
	if err != nil {
		return nil, errors.ParseError(serviceErrors.ApplicationError{
			StatusCode: 401,
			Code:       "Unauthorized",
			Errors:     map[string]string{"message": "Invalid user ID format"},
		})
	}

	projectID, err := uuid.Parse(req.ProjectId)
	if err != nil {
		return nil, errors.ParseError(serviceErrors.ApplicationError{
			StatusCode: 401,
			Code:       "Unauthorized",
			Errors:     map[string]string{"message": "Invalid project ID format"},
		})
	}

	project, err := s.projectService.GetProject(ctx, projectID)
	if err != nil {
		return nil, errors.ParseError(err)
	}
	if project.OwnerID != userID {
		return nil, errors.ParseError(serviceErrors.ApplicationError{
			StatusCode: 401,
			Code:       "Unauthorized",
			Errors:     map[string]string{"message": "User don't have acces to edit this project"},
		})
	}
	updatedProject, err := s.projectService.Update(ctx, projectID, req.NewName)
	if err != nil {
		return nil, errors.ParseError(err)
	}
	return &pb.Project{
		Id:        updatedProject.ID.String(),
		Name:      updatedProject.Name,
		OwnerId:   updatedProject.OwnerID.String(),
		CreatedAt: timestamppb.New(updatedProject.CreatedAt),
		UpdatedAt: nullableTimeToTimestamppb(&updatedProject.UpdatedAt),
	}, nil
}

func (s *ProjectServer) Delete(ctx context.Context, req *pb.DeleteRequest) (*emptypb.Empty, error) {
	userIDVal := ctx.Value("user_id")
	if userIDVal == nil {
		return nil, errors.ParseError(serviceErrors.ApplicationError{
			StatusCode: 401,
			Code:       "Unauthorized",
			Errors:     map[string]string{"message": "User ID not found in context"},
		})
	}
	userID, err := uuid.Parse(userIDVal.(string))
	if err != nil {
		return nil, errors.ParseError(serviceErrors.ApplicationError{
			StatusCode: 401,
			Code:       "Unauthorized",
			Errors:     map[string]string{"message": "Invalid user ID format"},
		})
	}

	projectID, err := uuid.Parse(req.ProjectId)
	if err != nil {
		return nil, errors.ParseError(serviceErrors.ApplicationError{
			StatusCode: 401,
			Code:       "Unauthorized",
			Errors:     map[string]string{"message": "Invalid project ID format"},
		})
	}

	project, err := s.projectService.GetProject(ctx, projectID)
	if err != nil {
		return nil, errors.ParseError(err)
	}
	if project.OwnerID != userID {
		return nil, errors.ParseError(serviceErrors.ApplicationError{
			StatusCode: 401,
			Code:       "Unauthorized",
			Errors:     map[string]string{"message": "User don't have acces to edit this project"},
		})
	}
	err = s.projectService.Delete(ctx, projectID)
	if err != nil {
		return nil, errors.ParseError(err)
	}
	return &emptypb.Empty{}, nil
}

func nullableTimeToTimestamppb(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	} else {
		return nil
	}
}
