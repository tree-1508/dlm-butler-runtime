//go:build dld_runtime

package main

import (
	"github.com/itchio/butler/cmd/daemon"
	"github.com/itchio/butler/mansion"
)

// DLD-051 signed runtime exposes only the supervised daemon command.
// The normal upstream-compatible CLI remains available in builds without dld_runtime.
func registerCommands(ctx *mansion.Context) {
	daemon.Register(ctx)
}
