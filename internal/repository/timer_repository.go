package repository

import (
	"database/sql"

	"controltasks/internal/model"
)

type TimerRepository struct {
	db *sql.DB
}

func NewTimerRepository(db *sql.DB) *TimerRepository {
	return &TimerRepository{db: db}
}

// scanTimer lê uma linha da tabela active_timers.
func scanTimer(row interface {
	Scan(...any) error
}) (*model.ActiveTimer, error) {
	var t model.ActiveTimer
	err := row.Scan(
		&t.ID, &t.UserID, &t.Status,
		&t.StartedAt, &t.ElapsedSeconds,
		&t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// FindByUserID retorna o timer ativo do usuário, ou nil se não existir.
func (r *TimerRepository) FindByUserID(userID string) (*model.ActiveTimer, error) {
	row := r.db.QueryRow(`
		SELECT id, user_id, status, started_at, elapsed_seconds, created_at, updated_at
		FROM active_timers
		WHERE user_id = $1`, userID)

	t, err := scanTimer(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return t, err
}

// Create insere um novo timer em estado running para o usuário.
func (r *TimerRepository) Create(userID string, initialSeconds int) (*model.ActiveTimer, error) {
	row := r.db.QueryRow(`
		INSERT INTO active_timers (user_id, status, started_at, elapsed_seconds)
		VALUES ($1, 'running', NOW(), $2)
		RETURNING id, user_id, status, started_at, elapsed_seconds, created_at, updated_at`,
		userID, initialSeconds)

	return scanTimer(row)
}

// UpdatePause pausa o timer e persiste o elapsed acumulado.
func (r *TimerRepository) UpdatePause(userID string, elapsedSeconds int) (*model.ActiveTimer, error) {
	row := r.db.QueryRow(`
		UPDATE active_timers
		SET status = 'paused',
		    elapsed_seconds = $2,
		    started_at = NULL,
		    updated_at = NOW()
		WHERE user_id = $1
		RETURNING id, user_id, status, started_at, elapsed_seconds, created_at, updated_at`,
		userID, elapsedSeconds)

	t, err := scanTimer(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return t, err
}

// UpdateResume retoma o timer pausado, atualizando started_at para NOW().
func (r *TimerRepository) UpdateResume(userID string) (*model.ActiveTimer, error) {
	row := r.db.QueryRow(`
		UPDATE active_timers
		SET status = 'running',
		    started_at = NOW(),
		    updated_at = NOW()
		WHERE user_id = $1
		RETURNING id, user_id, status, started_at, elapsed_seconds, created_at, updated_at`,
		userID)

	t, err := scanTimer(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return t, err
}

// Delete remove o timer do usuário.
func (r *TimerRepository) Delete(userID string) error {
	res, err := r.db.Exec(`DELETE FROM active_timers WHERE user_id = $1`, userID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
