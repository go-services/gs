---
title: "@http"
weight: 3
---

## @http

The `@http` annotation allows you to define HTTP endpoints in your service.


### How to use it?

You can use the `@http` annotation in interfaces that are annotated with `@service`.

The `@http` annotation needs to be added to a method in a `Service` interface.

```go
// Service is the interface that provides the methods.
//
// @service(name="example", base="/example")
type Service interface {
    // Abc method
    //
    // @http(method="post", path="/abc")
    Abc(context.Context, SumRequest) (*SumResponse, error)
}
```

### Parameters

The `@http` annotation accepts the following parameters:

- **method**: The HTTP method of the endpoint. Possible values are: GET, POST, PUT, PATCH, DELETE
  - **Required**: Yes
  - **Default**: None
- **path**: The path of the endpoint. Use the format specified in [go-chi/chi](https://go-chi.io/#/pages/routing)
  - **Required**: Yes
  - **Default**: None
- **responseFormat**: The response format of the service. Possible values are: JSON, XML
  - **Required**: No
  - **Default**: JSON
