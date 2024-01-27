package xyz

import "context"

// XyzService is an interface that represents a service
type XyzService interface {
	MyMethod(ctx context.Context) error
}

type xyzService struct {
}

func (s *xyzService) MyMethod(ctx context.Context) error {

	return nil
}

func NewXyzService() XyzService {
	return &xyzService{}
}
