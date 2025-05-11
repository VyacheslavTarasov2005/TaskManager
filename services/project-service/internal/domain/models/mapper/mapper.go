package mapper

import (
	"project-service/internal/domain/models"
	"project-service/internal/domain/models/dto"
)

func ToGetProject(project models.Project) dto.GetProject {
	return dto.GetProject{
		ID:        project.ID,
		Name:      project.Name,
		OwnerID:   project.OwnerID,
		CreatedAt: project.CreatedAt,
		UpdatedAt: project.UpdatedAt,
	}
}

func ToProject(dtoProject dto.CreateUpdateProject) models.Project {
	return models.Project{
		Name:    dtoProject.Name,
		OwnerID: dtoProject.OwnerID,
	}
}

func ToGetProjectList(projects []models.Project) []dto.GetProject {
	result := make([]dto.GetProject, len(projects))
	for i, p := range projects {
		result[i] = ToGetProject(p)
	}
	return result
}
