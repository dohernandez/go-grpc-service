//go:build never
// +build never

package noprune

import (
	_ "github.com/bool64/dev"                                // Include CI/Dev scripts to project.
	_ "github.com/dohernandez/dev-grpc"                      // Include development grpc helpers to proje
	_ "github.com/dohernandez/go-grpc-service/dev/makefiles" // Include development makefiles to project.
	_ "github.com/dohernandez/go-grpc-service/dev/scripts"   // Include development scripts to project.
)
