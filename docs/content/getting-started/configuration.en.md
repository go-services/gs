---
title: "Configuration"
weight: 3
---

## The `gs.yaml` file
By default GS will look for a `gs.yaml` file in the current project directory. This file contains the configuration for the project.

If no `gs.yaml` file is found, GS will use the default configuration.

{{<code `gs.yaml` `yaml`>}}

```yaml
paths:
  gen: gen
  config: config
```

{{</code>}}
