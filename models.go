package tasks

import "time"

type (
	Task struct {
		ID      int64     `db:"id" json:"id"`
		Name    string    `db:"name" json:"name"`
		Desc    string    `db:"desc" json:"desc"`
		Expires time.Time `db:"expires" json:"expires"`
		Done    bool      `db:"done" json:"done"`
	}

	TasksUsers struct {
		UserID int64 `db:"user_id"`
		TaskID int64 `db:"task_id"`
	}
)
