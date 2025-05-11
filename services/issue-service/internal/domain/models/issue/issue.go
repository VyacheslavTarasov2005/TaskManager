package issue

import (
	"github.com/google/uuid"
	"time"
)

type Issue struct {
	ID          uuid.UUID     `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	ProjectID   uuid.UUID     `json:"project_id"`
	AssigneeID  uuid.UUID     `json:"assignee_id"`
	ReporterID  uuid.UUID     `json:"reporter_id"`
	Status      IssueStatus   `json:"status"`
	Priority    IssuePriority `json:"priority"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}
