package abc

import "context"

// AbcService is an interface that represents a service
//
// @service(name="service1", base="/abc")
type AbcService interface {
	// MyMethod is a method that does something
	//
	// @http(method="GET", path="/my-method")
	MyMethod(ctx context.Context) error
}

type abcService struct {
}

func (s *abcService) MyMethod(ctx context.Context) error {

	return nil
}

func NewXyzService() AbcService {
	return &abcService{}
}
