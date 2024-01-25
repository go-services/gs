---
title: "Cron"
weight: 2
---

## Cron Jobs

GS can also generate cron jobs for you. A cron job is a service that runs periodically. It can be used to perform maintenance tasks or to send notifications.

A cron job in GS is defined by a `Cron` interface. This interface contains one method that will be called periodically.

Under the hood, GS uses [robfig/cron](https://github.com/robfig/cron) to schedule the execution of the `Run` method.

```go
// Cron is the interface that provides the Run method.
//
// @cron(name="example", schedule="@every 10m")
type Cron interface {
    // Run is called periodically according to a schedule.
    Run()
}
```

You only need to implement this interface and GS will generate the rest of the code for you.
