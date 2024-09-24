package entity

import (
	"time"
)

type TaskGroupDocument struct {
	ID           uint      `db:"id" json:"id"`
	TaskGroupID  uint      `db:"task_group_id" json:"task_group_id"`
	Document1URL string    `db:"document1_url" json:"document1_url"`
	Document2URL string    `db:"document2_url" json:"document2_url"`
	Document3URL string    `db:"document3_url" json:"document3_url"`
	Document4URL string    `db:"document4_url" json:"document4_url"`
	Document5URL string    `db:"document5_url" json:"document5_url"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

func NewTaskGroupDocuemnt(
	taskGroupID uint,
	document1URL string,
	document2URL string,
	document3URL string,
	document4URL string,
	document5URL string,
) *TaskGroupDocument {
	return &TaskGroupDocument{
		TaskGroupID:  taskGroupID,
		Document1URL: document1URL,
		Document2URL: document2URL,
		Document3URL: document3URL,
		Document4URL: document4URL,
		Document5URL: document5URL,
	}
}

type UpdateTaskGroupDocumentParam struct {
	DocumentType uint   `db:"document_type" json:"document_type" validate:"required"`
	TaskGroupID  uint   `db:"task_group_id" json:"task_group_id" validate:"required"`
	DocumentURL  string `db:"document_url" json:"document_url"`
}
