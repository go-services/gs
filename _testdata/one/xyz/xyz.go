package xyz

import "context"

type Request struct {
	Abc string `json:"abc"`
}

type Response struct {
	Abc string `json:"abc"`
}

// XyzService is an interface that represents a service
//
// @service()
type XyzService interface {
	// MyMethod is a method that does something
	//
	// @http(method="GET", path="/my-method")
	MyMethod(ctx context.Context, req Request) (*Response, error)
}

type xyzService struct {
}

func (s *xyzService) MyMethod(ctx context.Context, req Request) (*Response, error) {

	return &Response{
		Abc: req.Abc,
	}, nil
}

func NewXyzService() XyzService {
	return &xyzService{}
}
