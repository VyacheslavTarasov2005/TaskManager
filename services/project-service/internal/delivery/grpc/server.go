package grpc

import (
	"context"
	"project-service/internal/delivery/grpc/auth"
	"project-service/internal/delivery/grpc/errors"
	"project-service/internal/delivery/grpc/pb"
	serviceErrors "project-service/internal/service/errors"
	"project-service/internal/service/interfaces"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ProjectServer struct {
	pb.UnimplementedProjectServiceServer
	authClient         *auth.UserServiceClient
	projectService     interfaces.ProjectService
	projectUserService interfaces.ProjectUserService
}

func NewProjectServer(projectService interfaces.ProjectService, projectUserService interfaces.ProjectUserService, cli *auth.UserServiceClient) *ProjectServer {
	return &ProjectServer{
		projectService:     projectService,
		projectUserService: projectUserService,
		authClient:         cli,
	}
}

func (s *ProjectServer) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
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
			StatusCode: 400,
			Code:       "BadRequest",
			Errors:     map[string]string{"message": "Invalid user ID format"},
		})
	}

	projectId, err := s.projectService.Create(ctx, req.GetName(), userID)
	if err != nil {
		return nil, err
	}

	err = s.projectUserService.AddToProject(ctx, userID, *projectId, pb.Role_name[int32(pb.Role_OWNER)])
	if err != nil {
		return nil, errors.ParseError(err)
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
		UpdatedAt: nullableTimeToTimestamppb(project.UpdatedAt),
	}, nil
}

func (s *ProjectServer) GetMyProjects(ctx context.Context, req *pb.GetMyProjectsRequest) (*pb.GetMyProjectsResponse, error) {
	query := ""
	if req.Query != nil {
		query = *req.Query
	}

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
			StatusCode: 400,
			Code:       "BadRequest",
			Errors:     map[string]string{"message": "Invalid user ID format"},
		})
	}

	//изменить с получения проектов владельцем которого является пользователь на те в которых он находится
	projects, err := s.projectService.GetProjectsByUser(ctx, userID, query)
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
			UpdatedAt: nullableTimeToTimestamppb(project.UpdatedAt),
		})
	}
	return &pb.GetMyProjectsResponse{
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
			StatusCode: 400,
			Code:       "BadRequest",
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
			StatusCode: 403,
			Code:       "Forbidden",
			Errors:     map[string]string{"message": "User does not have access to modify this project"},
		})
	}
	updatedProject, err := s.projectService.Update(ctx, userID, projectID, req.NewName)
	if err != nil {
		return nil, errors.ParseError(err)
	}
	return &pb.Project{
		Id:        updatedProject.ID.String(),
		Name:      updatedProject.Name,
		OwnerId:   updatedProject.OwnerID.String(),
		CreatedAt: timestamppb.New(updatedProject.CreatedAt),
		UpdatedAt: nullableTimeToTimestamppb(project.UpdatedAt),
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
			StatusCode: 400,
			Code:       "BadRequest",
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
			StatusCode: 403,
			Code:       "Forbidden",
			Errors:     map[string]string{"message": "User does not have access to modify this project"},
		})
	}
	err = s.projectService.Delete(ctx, projectID)
	if err != nil {
		return nil, errors.ParseError(err)
	}

	err = s.projectUserService.DelteByProject(ctx, projectID)
	if err != nil {
		return nil, errors.ParseError(err)
	}
	return &emptypb.Empty{}, nil
}

func (s *ProjectServer) GetMyRole(ctx context.Context, req *pb.GetMyRoleRequest) (*pb.GetMyRoleResponse, error) {
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
			StatusCode: 400,
			Code:       "BadRequest",
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

	role, err := s.projectUserService.GetUserRole(ctx, projectID, userID)
	if err != nil {
		return nil, errors.ParseError(err)
	}
	if pb.Role_value[string(*role)] == int32(pb.Role_UNKNOWN) {
		return nil, errors.ParseError(serviceErrors.ApplicationError{
			StatusCode: 404,
			Code:       "NotFound",
			Errors:     map[string]string{"message": "Role not found for user in project"},
		})
	}
	return &pb.GetMyRoleResponse{
		Role: pb.Role(pb.Role_value[string(*role)]),
	}, nil
}

func (s *ProjectServer) KickFromProject(ctx context.Context, req *pb.KickFromProjectRequest) (*emptypb.Empty, error) {
	owner_id_val := ctx.Value("user_id")
	if owner_id_val == nil {
		return nil, errors.ParseError(serviceErrors.ApplicationError{
			StatusCode: 401,
			Code:       "Unauthorized",
			Errors:     map[string]string{"message": "User ID not found in context"},
		})
	}
	owner_id, err := uuid.Parse(owner_id_val.(string))
	if err != nil {
		return nil, errors.ParseError(serviceErrors.ApplicationError{
			StatusCode: 401,
			Code:       "Unauthorized",
			Errors:     map[string]string{"message": "Invalid user ID format"},
		})
	}

	project_id, err := uuid.Parse(req.ProjectId)
	if err != nil {
		return nil, errors.ParseError(serviceErrors.ApplicationError{
			StatusCode: 401,
			Code:       "Unauthorized",
			Errors:     map[string]string{"message": "Invalid project ID format"},
		})
	}

	user_id, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, errors.ParseError(serviceErrors.ApplicationError{
			StatusCode: 401,
			Code:       "Unauthorized",
			Errors:     map[string]string{"message": "Invalid project ID format"},
		})
	}
	//проверка на правильность id
	err = s.authClient.GetUserProfile(ctx, user_id.String())
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Ivalid user ID")
	}
	//проверка на права добавляющего
	role, err := s.projectUserService.GetUserRole(ctx, project_id, owner_id)
	if err != nil {
		return nil, errors.ParseError(err)
	}
	if string(*role) != pb.Role_name[int32(pb.Role_OWNER)] {
		return nil, errors.ParseError(serviceErrors.ApplicationError{
			StatusCode: 401,
			Code:       "Unauthorized",
			Errors:     map[string]string{"message": "Not enough rules to add user"},
		})
	}
	err = s.projectUserService.DeleteFromProject(ctx, project_id, user_id)
	if err != nil {
		return nil, errors.ParseError(err)
	}
	return &emptypb.Empty{}, nil
}

func (s *ProjectServer) AddToProject(ctx context.Context, req *pb.AddToProjectRequest) (*emptypb.Empty, error) {
	owner_id_val := ctx.Value("user_id")
	if owner_id_val == nil {
		return nil, errors.ParseError(serviceErrors.ApplicationError{
			StatusCode: 401,
			Code:       "Unauthorized",
			Errors:     map[string]string{"message": "User ID not found in context"},
		})
	}
	owner_id, err := uuid.Parse(owner_id_val.(string))
	if err != nil {
		return nil, errors.ParseError(serviceErrors.ApplicationError{
			StatusCode: 401,
			Code:       "Unauthorized",
			Errors:     map[string]string{"message": "Invalid user ID format"},
		})
	}

	project_id, err := uuid.Parse(req.ProjectId)
	if err != nil {
		return nil, errors.ParseError(serviceErrors.ApplicationError{
			StatusCode: 401,
			Code:       "Unauthorized",
			Errors:     map[string]string{"message": "Invalid project ID format"},
		})
	}

	user_id, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, errors.ParseError(serviceErrors.ApplicationError{
			StatusCode: 401,
			Code:       "Unauthorized",
			Errors:     map[string]string{"message": "Invalid project ID format"},
		})
	}
	//проверка на правильность id
	err = s.authClient.GetUserProfile(ctx, user_id.String())
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Ivalid user ID")
	}

	//проверка на права добавляющего
	role, err := s.projectUserService.GetUserRole(ctx, project_id, owner_id)
	if err != nil {
		return nil, errors.ParseError(err)
	}
	if string(*role) != pb.Role_name[int32(pb.Role_OWNER)] {
		return nil, errors.ParseError(serviceErrors.ApplicationError{
			StatusCode: 401,
			Code:       "Unauthorized",
			Errors:     map[string]string{"message": "Not enough rules to add user"},
		})
	}
	//добавление в проект
	err = s.projectUserService.AddToProject(ctx, user_id, project_id, pb.Role_name[int32(pb.Role_MEMBER)])
	if err != nil {
		return nil, errors.ParseError(err)
	}
	return &emptypb.Empty{}, nil
}

func nullableTimeToTimestamppb(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}
