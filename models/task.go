package models

import (
	"time"
	"gorm.io/gorm"
)

// TaskStatus mendefinisikan tipe data untuk status pengerjaan sebuah task
type TaskStatus string

const (
	TaskTodo       TaskStatus = "todo"
	TaskInProgress TaskStatus = "in_progress"
	TaskDone       TaskStatus = "done"
)

// Task merepresentasikan entitas tugas yang berelasi dengan project dan user
type Task struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	ProjectID    uint       `gorm:"not null" json:"project_id"`
	CreatedByID  uint       `gorm:"not null" json:"created_by_id"`  // Siapa yang buat? (PM/User)
	AssignedToID uint       `gorm:"not null" json:"assigned_to_id"` // Siapa yang ngerjain?
	
	Title        string     `gorm:"type:varchar(255);not null" json:"title"`
	Description  string     `gorm:"type:text" json:"description"`
	Status       TaskStatus `gorm:"type:varchar(20);default:'todo'" json:"status"`
	
	// Relasi untuk reverse query (Optional tapi berguna)
	Timesheets   []Timesheet `gorm:"foreignKey:TaskID" json:"timesheets,omitempty"`

	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
