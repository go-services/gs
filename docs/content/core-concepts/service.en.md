---
title: "Service"
weight: 1
---

## Services

Services are the building blocks of your project. They are the main components that you will use to build your app.

***GS*** uses [go-kit](https://gokit.io/) under the hood to stitch together the different parts of your service.

A service in ***GS*** is defined by a `Service` interface. This interface contains the methods that can be called from other services or clients.

```go
// Service is the interface that provides the Sum and Concat methods.
//
// @service(name="example", base="/example")
type Service interface {
	// Sum adds together two integers and returns the result.
	//
	// @http(method="post", path="/sum")
	Sum(context.Context, SumRequest) (*SumResponse, error)
	// Concat concatenates two strings and returns the result.
	//
	//@http(method="post", path="/concat")
	Concat(context.Context, ConcatRequest) (*ConcatResponse, error)
}
```

### How is a service designed under the hood?

Go kit tries to enforce strict **separation of concerns** through using middleware (or decorator) pattern.

{{< image src="/images/docs/onion.png" zoomable="true" >}}


Learn more about go-kit [here](https://gokit.io/faq/).



