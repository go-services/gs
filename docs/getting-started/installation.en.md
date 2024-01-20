---
title: "Installation"
weight: 1
---

## Get Started With GS

To get started with GS, you need to install go and set up your development environment. You can follow the [official Go installation guide](https://golang.org/doc/install) to install Go on your machine.

---


## Install GS
Once you have installed Go, you can install GS by running the following command:

{{<code `command line` `console`>}}

```bash
go get -u github.com/go-services/gs
```

{{</code>}}

---

## Create a new project

To create a new project, run the following command:

{{<code `command line` `console`>}}

```bash
gs new project my-project
```

{{</code>}}

This will create a new project in the `my-project` folder.

The project will contain a sample service. You can remove/modify it as you wish.

## Generate the project

To generate the project, run the following command:

{{<code `command line` `console`>}}

```bash
gs generate
```
{{</code>}}
