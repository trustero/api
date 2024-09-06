# Trustero Receptor SDK for Golang

![Build CI](https://github.com/trustero/api/actions/workflows/ci.yml/badge.svg)
[![CodeQL](https://github.com/trustero/api/actions/workflows/codeql.yml/badge.svg)](https://github.com/trustero/api/actions/workflows/codeql.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/trustero/api/go)](https://goreportcard.com/report/github.com/trustero/api/go)
[![Go Reference](https://pkg.go.dev/badge/github.com/trustero/api.svg)](https://pkg.go.dev/github.com/trustero/api/go/receptor_sdk)


ReceptorSDK documentation is available on [pkg.go.dev](https://pkg.go.dev/github.com/trustero/api/go/receptor_sdk). Users are advised to consult this ReceptorSDK documentation together with the comprehensive `Receptor Developer Guide` and [protobuf definitions](../docs/receptor_v1/receptor.md). To obtain a copy of the guide, please reach out to Trustero Support. The easiest way to learn about the SDK is to consult the set of [examples](examples/) built on top of the SDK. What follows is a subset of these examples that can be found useful as stand-alone programs.

| Example                                         | Description                                    |
| :---------------------------------------------- | :--------------------------------------------- |
| [GitLab Receptor](examples/gitlab_receptor/) | A Receptor that posts GitLab users as evidence |

ReceptorSDK is an open source Trustero project and contributions are welcome.

ReceptorSDK is periodically refreshed to reflect the newest additions to the Trustero API. Users of the SDK are advised to track the latest releases rather closely to ensure proper function in the unlikely event of an incompatible change to a Trustero API.

## Installation

```
go get github.com/trustero/api/go@latest
```



## Usage

The developer needs to implement the [Receptor interface](https://pkg.go.dev/github.com/trustero/api/go/receptor_sdk#Receptor) in their project to create a Receptor.

### Required Functions

As a developer, you will need to implement the following functions to have a working receptor:
1. func (r *Receptor) GetReceptorType() string {}
    - This function will return a string, signifying the receptor type (e.g. “trr-gitlab”)
2. func (r *Receptor) GetKnownServices() []string {}
    - This function will return an array of string, signifying a list of service types the receptor collects
3. func (r *Receptor) GetCredentialObj() (credentialObj interface{}) {}
4. func (r *Receptor) GetConfigObj() (configObj interface{}) {}
5. func (r *Receptor) GetConfigObjDesc() (configObjDesc interface{}) {}
6. func (r *Receptor) GetAuthMethods() (authMethods interface{}) {}
7. func (r *Receptor) Verify(credentials interface{}) (ok bool, err error) {}
8. func (r *Receptor) Discover(credentials interface{}) (svcs []*receptor_v1.ServiceEntity, err error) {}
9. func (r *Receptor) Report(credentials interface{}) (evidences []*receptor_sdk.Evidence, err error) {}
10. func (r *Receptor) GetEvidenceInfo() (evidences []*receptor_sdk.Evidence) {}

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

func (r *Receptor) GetConfigObj() (configObj interface{}){
	return
}

func (r *Receptor) GetConfigObjDesc() (configObjDesc interface{}) {
	return
}

func (r *Receptor)	GetConfigObjDesc() (configObjDesc interface{}) {
	return
}

func (r *Receptor) GetAuthMethods() (authMethods interface{}) {
	return
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

func (r *Receptor) GetEvidenceInfo() (evidences []*receptor_sdk.Evidence) {
   //YOUR CODE HERE
   return
}

func main() {
	cmd.Execute(&Receptor{})
}
```

A real-life example can be found in the [examples](examples/) directory.

## Testing A Receptor

You should be able to run your receptor code via the command line to confirm the Verify and Scan functions produce the correct output.
You can run the Receptor code with the `dryrun` flag and it will print the output to the console.
You can compile the Receptor code into a binary or run the main file directly.
If you run the main file directly, your command should look something like this:

```
go run main.go scan dryrun --find-evidence
```

This command will run the Verify, Discover, and Report functions that you wrote and print their output to the console. You should be able to see the final Evidences that are generated by the receptor.
