# Trustero Receptor SDK

![Build CI](https://github.com/trustero/api/actions/workflows/ci.yml/badge.svg)
[![CodeQL](https://github.com/trustero/api/actions/workflows/codeql.yml/badge.svg)](https://github.com/trustero/api/actions/workflows/codeql.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/trustero/api/go)](https://goreportcard.com/report/github.com/trustero/api/go)
[![Go Reference](https://pkg.go.dev/badge/github.com/trustero/api.svg)](https://pkg.go.dev/github.com/trustero/api/go/receptor_sdk)

This repository contains the SDK for creating a Receptor that integrates with
the Trustero Service.

ReceptorSDK documentation is available on [pkg.go.dev](https://pkg.go.dev/github.com/trustero/api/go/receptor_sdk). Users are advised to consult this ReceptorSDK documentation together with the comprehensive `Receptor Developer Guide` and [protobuf definitions](docs/receptor_v1/receptor.md). To obtain a copy of the guide, please reach out to Trustero Support. The easiest way to learn about the SDK is to consult the set of [examples](go/examples/) built on top of the SDK. What follows is a subset of these examples that can be found useful as stand-alone programs.

| Example                                         | Description                                    |
| :---------------------------------------------- | :--------------------------------------------- |
| [GitLab Receptor](go/examples/gitlab_receptor/) | A Receptor that posts GitLab users as evidence |

ReceptorSDK is an open source Trustero project and contributions are welcome.

ReceptorSDK is periodically refreshed to reflect the newest additions to the Trustero API. Users of the SDK are advised to track the latest releases rather closely to ensure proper function in the unlikely event of an incompatible change to a Trustero API.

## Installation

```
go get github.com/trustero/api/go@latest
```

## Usage

A developer does not need to initialize a client to interact with the Trustero service.

All communication to the Trustero service is handled by the ReceptorSDK internally.

The developer needs to implement the [Receptor interface](https://pkg.go.dev/github.com/trustero/api/go/receptor_sdk#Receptor) in their project to create a Receptor.

```go
package main

import (
	"github.com/trustero/api/go/receptor_sdk"
	"github.com/trustero/api/go/receptor_sdk/cmd"
	"github.com/trustero/api/go/receptor_v1"
)

const (
	receptorName = "trr-custom"
	serviceName = "Custom Service"
)

type Receptor struct {
  // YOUR CODE HERE
}

func (r *Receptor) GetReceptorType() string {
	return receptorName
}

func (r *Receptor) GetKnownServices() []string {
	return []string{serviceName}
}

func (r *Receptor) GetCredentialObj() (credentialObj interface{}) {
	return r
}

func (r *Receptor) Verify(credentials interface{}) (ok bool, err error) {
  // YOUR CODE HERE
	return
}

func (r *Receptor) Discover(credentials interface{}) (svcs []*receptor_v1.ServiceEntity, err error) {
  // YOUR CODE HERE
	return
}

func (r *Receptor) Report(credentials interface{}) (evidences []*receptor_sdk.Evidence, err error) {
  // YOUR CODE HERE
	return
}

func main() {
	cmd.Execute(&Receptor{})
}
```

A real-life example can be found in the [examples](go/examples/) directory.
