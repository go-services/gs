---
title: "Configuration"
weight: 3
---

## The `gs.json` file
By default GS will look for a `gs.json` file in the current project directory. This file contains the configuration for the project.

If no `gs.json` file is found, GS will use the default configuration.

{{<code `gs.json` `json`>}}

```json
{
  "paths": {
    "gen": "gen",
    "cmd": "cmd"
  }
}
```

{{</code>}}
