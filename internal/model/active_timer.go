package model

import "time"

// TimerStatus representa o estado de um timer ativo.
type TimerStatus string

const (
	TimerRunning TimerStatus = "running"
	TimerPaused  TimerStatus = "paused"
)

// ActiveTimer representa o timer ativo de um usuário.
type ActiveTimer struct {
	ID             string      `json:"id"`
	UserID         string      `json:"user_id"`
	Status         TimerStatus `json:"status"`
	StartedAt      *time.Time  `json:"started_at"`      // instante em que entrou em running; nil quando paused
	ElapsedSeconds int         `json:"elapsed_seconds"` // segundos acumulados antes da última pausa
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
}
