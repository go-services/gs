package xyz

import "context"

// XyzService is an interface that represents a service
//
// @service(name="no_constructor")
type XyzService interface {

	// SecondMethod is a method that does something
	//
	// @http(method="GET", path="/my-method")
	SecondMethod(ctx context.Context) error
	// MyMethod is a method that does something
	//
	// @http(method="GET", path="/my-method")
	MyMethod(ctx context.Context) error
}

type xyzService struct {
}

func (s *xyzService) MyMethod(ctx context.Context) error {

	return nil
}
func (s *xyzService) SecondMethod(ctx context.Context) error {

	return nil
}

func NewXyzService() XyzService {
	return &xyzService{}
}
