---
title: "@service"
weight: 2
---

## @service

The `@service` annotation allows you to define a service in your project.


### How to use it?

You can use the `@service` annotation in any go file in your project.

The `@service` annotation needs to be added to a `Service` interface.

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

The `@service` annotation accepts the following parameters:

- **name**: The name of the service. The value will be formatted to snakeCase
  - **Required**: No
  - **Default**: `snakeCase({interfaceName})`
- **base**: The base path of the service.
  - **Required**: No
  - **Default**: `/{serviceName}`
