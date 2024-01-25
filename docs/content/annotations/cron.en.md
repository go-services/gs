---
title: "@cron"
weight: 4
---

## @cron

The `@cron` annotation allows you to define cron jobs in your project.

### How to use it?

You can use the `@cron` annotation in any go file in your project.

The `@cron` annotation needs to be added to a `Cron` interface.

```go
// Cron is the interface that provides the Run method.
//
// @cron(name="example", schedule="@every 10m")
type Cron interface {
    // Run is called periodically according to a schedule.
    Run()
}
```

### Parameters

The `@cron` annotation accepts the following parameters:

- **name**: The name of the cron job. The value will be formatted to snakeCase.
  - **Required**: No
  - **Default**: `snakeCase({interfaceName})`
- **schedule**: The schedule of the cron job, it accepts the same format as [robfig/cron](https://github.com/robfig/cron)
  - **Required**: Yes
  - **Default**: None
