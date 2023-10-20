package repository

import (
	"context"

	"github.com/google/uuid"
)

type (
	// tasks table
	Task struct {
		ID        uuid.UUID `db:"id"`
		UserID    uuid.UUID `db:"user_id"`
		Title     string    `db:"title"`
		IsDone    bool      `db:"is_done"`
		CreatedAt string    `db:"created_at"`
	}

	SearchTasksParams struct {
		UserID uuid.UUID
		Target string
	}

	CreateTaskParams struct {
		UserID uuid.UUID
		Title  string
	}

	UpdateTaskParams struct {
		ID     uuid.UUID
		Title  string
		IsDone bool
	}
)

func (r *Repository) GetTasks(ctx context.Context, UserID uuid.UUID) ([]Task, error) {
	tasks := []Task{}
	if err := r.db.SelectContext(ctx, &tasks, "SELECT * FROM tasks WHERE user_id = ?", UserID); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *Repository) SearchTasks(ctx context.Context, params SearchTasksParams) ([]Task, error) {
	tasks := []Task{}

	// TODO

	return tasks, nil
}

func (r *Repository) CreateTask(ctx context.Context, params CreateTaskParams) error {
	taskID := uuid.New()
	if _, err := r.db.ExecContext(ctx, "INSERT INTO tasks (id, user_id, title) VALUES (?, ?, ?)", taskID, params.UserID, params.Title); err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpdateTask(ctx context.Context, params UpdateTaskParams) error {
	if _, err := r.db.ExecContext(ctx, "UPDATE tasks SET title = ?, is_done = ? WHERE id = ?", params.Title, params.IsDone, params.ID); err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteTask(ctx context.Context, taskID uuid.UUID) error {
	if _, err := r.db.ExecContext(ctx, "DELETE FROM tasks WHERE id = ?", taskID); err != nil {
		return err
	}

	return nil
}
