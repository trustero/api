# Trustero Receptor SDK for Golang

This repository contains the SDK for creating a Receptor that integrates with
the Trustero Service.

---

## Documentation

- go.dev: https://pkg.go.dev/github.com/trustero/api/go/receptor_sdk

- [proto documentation](docs/receptor_v1/receptor.md)

## Installation

```
go get github.com/trustero/api/go@latest
```

## Example

There is an example of a receptor included in the [examples directory](go/examples/)

The included example implements a Receptor for the GitLab service. It produces
a list of users that have access to GitLab account you have provided credentials for.

The GitLab Receptor is implemented using the following functions:

1.  GetReceptorType
2.  GetKnownServices
3.  GetCredentialObj
4.  Verify
5.  Discover
6.  Report

Every Receptor will need to implement these 6 functions to integrate with Trustero.

The majority of your code will be written in the `Verify`, `Discover`, and `Report` functions
