---
title: "Services"
weight: 1
---

## Services

Services are the building blocks of your project. They are the main components that you will use to build your infrastructure.

Services are defined using the `@service` annotation. The annotation can be used on service interfaces.

```go
// Service is the interface that provides the Sum and Concat methods.
//
// @service(name="numbers", route="/numbers")
type Service interface {
	// Sum adds together two integers and returns the result.
	//
	// @http(method="post", route="/sum")
	Sum(context.Context, SumRequest) (*SumResponse, error)
	// Concat concatenates two strings and returns the result.
	//
	//@http(method="post", route="/concat")
	Concat(context.Context, ConcatRequest) (*ConcatResponse, error)
}
```
