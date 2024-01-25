---
title: "What is GS"
weight: 1
---
## What is GS?

GS is a code generator for Go services. It generates the boilerplate code for your services so you can focus on the business logic.

## Services

In GS, a service is a collection of methods that can be called from other services or clients. Currently, GS supports HTTP but we plan to add support for other protocols in the future.

GS uses [go-kit](https://gokit.io/) under the hood to stitch together the different parts of your service.

A service in GS is defined by a `Service` interface. This interface contains the methods that can be called from other services or clients.

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
    // @http(method="post", path="/concat")
    Concat(context.Context, ConcatRequest) (*ConcatResponse, error)
}
```

You only need to implement this interface and GS will generate the rest of the code for you.

Learn more about services [here](/core/service).

## Cron Jobs

GS can also generate cron jobs for you. A cron job is a service that runs periodically. It can be used to perform maintenance tasks or to send notifications.

A cron job in GS is defined by a `Cron` interface. This interface contains one method that will be called periodically.

```go
// Cron is the interface that provides the Run method.
//
// @cron(name="example", schedule="@every 10m")
type Cron interface {
    // Run is called periodically according to a schedule.
    Run()
}
```

