package service

import (
	"errors"
	"time"

	"controltasks/internal/model"
	"controltasks/internal/repository"
)

var ErrTimerAlreadyActive = errors.New("timer already active")

type TimerService struct {
	repo *repository.TimerRepository
}

func NewTimerService(repo *repository.TimerRepository) *TimerService {
	return &TimerService{repo: repo}
}

// Get retorna o timer ativo do usuário, ou nil se não existir.
func (s *TimerService) Get(userID string) (*model.ActiveTimer, error) {
	return s.repo.FindByUserID(userID)
}

// Start cria um novo timer em estado running. Retorna ErrTimerAlreadyActive (HTTP 409) se já existe um.
func (s *TimerService) Start(userID string, initialSeconds int) (*model.ActiveTimer, error) {
	existing, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrTimerAlreadyActive
	}
	return s.repo.Create(userID, initialSeconds)
}

// Pause pausa o timer ativo, calculando o elapsed acumulado antes de persistir.
func (s *TimerService) Pause(userID string) (*model.ActiveTimer, error) {
	timer, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	if timer == nil {
		return nil, nil
	}

	newElapsed := timer.ElapsedSeconds
	if timer.StartedAt != nil {
		newElapsed += int(time.Since(*timer.StartedAt).Seconds())
	}

	return s.repo.UpdatePause(userID, newElapsed)
}

// Resume retoma o timer pausado.
func (s *TimerService) Resume(userID string) (*model.ActiveTimer, error) {
	return s.repo.UpdateResume(userID)
}

// Delete remove o timer do usuário.
func (s *TimerService) Delete(userID string) error {
	return s.repo.Delete(userID)
}
