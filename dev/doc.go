// Package dev contains reusable development helpers.
package dev

// These imports workaround `go mod vendor` prune.
//
// See https://github.com/golang/go/issues/26366.
import (
	_ "github.com/dohernandez/go-grpc-service/dev/makefiles"
	_ "github.com/dohernandez/go-grpc-service/dev/scripts"
)
