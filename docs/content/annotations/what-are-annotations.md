---
title: "Annotations"
weight: 1
---
## What are GS annotations?

GS annotations are comments that you can add to your code to tell GS how to generate the code for your services.

GS annotations use a the following format:

```go
// @annotationName(key="value", key="value")
```

The source code for parsing annotations can be found [here](https://github.com/go-services/annotation).

## How to use GS annotations?

GS annotations can be used to define services, cron jobs, and more. Add annotations as comments to your code in the appropriate places.

