package issue

import "github.com/google/uuid"

type IssueFilter struct {
	ProjectID  *uuid.UUID     `json:"project_id,omitempty"`
	AssigneeID *uuid.UUID     `json:"assignee_id,omitempty"`
	ReporterID *uuid.UUID     `json:"reporter_id,omitempty"`
	Status     *IssueStatus   `json:"status,omitempty"`
	Priority   *IssuePriority `json:"priority,omitempty"`
}
