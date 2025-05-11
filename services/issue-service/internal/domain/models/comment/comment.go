package comment

import (
	"github.com/google/uuid"
	"time"
)

type Comment struct {
	ID        uuid.UUID `json:"id"`
	IssueID   uuid.UUID `json:"issue_id"`
	UserID    uuid.UUID `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
