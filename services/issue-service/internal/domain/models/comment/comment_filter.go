package comment

import "github.com/google/uuid"

type CommentFilter struct {
	IssueID *uuid.UUID `json:"issue_id,omitempty"`
	UserID  *uuid.UUID `json:"user_id,omitempty"`
}
