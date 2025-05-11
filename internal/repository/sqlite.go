package repository

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(dbPath string) (*SQLiteRepository, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			login TEXT UNIQUE,
			password TEXT
		);
		
		CREATE TABLE IF NOT EXISTS tasks (
			id TEXT PRIMARY KEY,
			user_id INTEGER,
			expression TEXT,
			arg1 TEXT,
			arg2 TEXT,
			operation TEXT,
			result REAL,
			status TEXT DEFAULT 'pending',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			started_at DATETIME,
			completed_at DATETIME,
			FOREIGN KEY(user_id) REFERENCES users(id)
		);
		
		CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
		CREATE INDEX IF NOT EXISTS idx_tasks_user ON tasks(user_id);
	`); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return &SQLiteRepository{db: db}, nil
}

func (r *SQLiteRepository) CreateUser(login, password string) (int64, error) {
	res, err := r.db.Exec(
		"INSERT INTO users (login, password) VALUES (?, ?)",
		login, password,
	)
	if err != nil {
		return 0, ErrUserExists
	}
	return res.LastInsertId()
}

func (r *SQLiteRepository) Authenticate(login, password string) (*User, error) {
	var user User
	err := r.db.QueryRow(
		"SELECT id, login, password FROM users WHERE login = ?",
		login,
	).Scan(&user.ID, &user.Login, &user.Password)

	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query user: %w", err)
	}
	if user.Password != password {
		return nil, ErrInvalidPassword
	}
	return &user, nil
}

func (r *SQLiteRepository) CreateTask(userID int, expr string) (string, error) {
	taskID := generateTaskID()
	_, err := r.db.Exec(
		"INSERT INTO tasks (id, user_id, expression, status) VALUES (?, ?, ?, 'pending')",
		taskID, userID, expr,
	)
	if err != nil {
		return "", fmt.Errorf("failed to create task: %w", err)
	}
	return taskID, nil
}

func (r *SQLiteRepository) GetPendingTask() (*Task, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	row := tx.QueryRow(`
		SELECT id, user_id, expression, arg1, arg2, operation 
		FROM tasks 
		WHERE status = 'pending' 
		ORDER BY created_at ASC 
		LIMIT 1
	`)

	var task Task
	err = row.Scan(&task.ID, &task.UserID, &task.Expression, &task.Arg1, &task.Arg2, &task.Operation)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan task: %w", err)
	}

	_, err = tx.Exec(`
		UPDATE tasks 
		SET status = 'in_progress', 
		    started_at = CURRENT_TIMESTAMP 
		WHERE id = ?
	`, task.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update task status: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &task, nil
}

func (r *SQLiteRepository) SaveResult(id string, result float64) error {
	_, err := r.db.Exec(`
		UPDATE tasks 
		SET status = 'completed', 
		    result = ?,
		    completed_at = CURRENT_TIMESTAMP 
		WHERE id = ?
	`, result, id)
	if err != nil {
		return fmt.Errorf("failed to save result: %w", err)
	}
	return nil
}

func (r *SQLiteRepository) GetUserTasks(userID int) ([]Task, error) {
	rows, err := r.db.Query(`
		SELECT id, expression, status, result, created_at, completed_at
		FROM tasks
		WHERE user_id = ?
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(
			&task.ID,
			&task.Expression,
			&task.Status,
			&task.Result,
			&task.CreatedAt,
			&task.CompletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *SQLiteRepository) GetTaskByID(id string) (*Task, error) {
	var task Task
	err := r.db.QueryRow(`
		SELECT id, user_id, expression, status, result, 
		       created_at, started_at, completed_at
		FROM tasks
		WHERE id = ?
	`, id).Scan(
		&task.ID,
		&task.UserID,
		&task.Expression,
		&task.Status,
		&task.Result,
		&task.CreatedAt,
		&task.StartedAt,
		&task.CompletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrTaskNotFound
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	return &task, nil
}

func generateTaskID() string {
	return fmt.Sprintf("task_%d", time.Now().UnixNano())
}
