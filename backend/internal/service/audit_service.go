package service

import (
	"github.com/pingan/bastion/internal/model"
	"github.com/pingan/bastion/internal/repository"
)

type AuditService struct {
	repo *repository.AuditRepo
}

func NewAuditService(repo *repository.AuditRepo) *AuditService {
	return &AuditService{repo: repo}
}

func (s *AuditService) Search(userID, assetID int64, keyword string, page, size int) ([]model.AuditLog, int, error) {
	offset := (page - 1) * size
	logs, err := s.repo.Search(userID, assetID, keyword, size, offset)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.repo.Count(userID, assetID, keyword)
	if err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}
