package xyz

// XyzService is an interface that represents a service
//
// @service()
type XyzService interface {
	// MyMethod is a method that does something
	//
	// @http(method="GET", path="/my-method")
	MyMethod() error
}

type xyzService struct {
}

func (s *xyzService) MyMethod() error {

	return nil
}

func NewXyzService() XyzService {
	return &xyzService{}
}
