module github.com/trustero/monorepo/ntrced/multipartkit/multipart_usage

go 1.21

require (
	github.com/rs/zerolog v1.33.0
	github.com/trustero/api/go v0.0.0
	google.golang.org/protobuf v1.34.1

)

require (
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	golang.org/x/sys v0.20.0 // indirect
)

replace github.com/trustero/api/go => ../../../go
