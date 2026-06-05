package service

import (
	"github.com/pingan/bastion/internal/model"
	"github.com/pingan/bastion/internal/repository"
)

type SessionService struct {
	repo *repository.SessionRepo
}

func NewSessionService(repo *repository.SessionRepo) *SessionService {
	return &SessionService{repo: repo}
}

func (s *SessionService) Create(session *model.Session) error {
	return s.repo.Create(session)
}

func (s *SessionService) Close(id int64) error {
	return s.repo.Close(id)
}

func (s *SessionService) List(limit, offset int) ([]model.Session, error) {
	return s.repo.FindAll(limit, offset)
}

func (s *SessionService) ActiveSessionCount() (int, error) {
	return s.repo.CountActive()
}

func (s *SessionService) TotalSessionCount() (int, error) {
	return s.repo.Count()
}
