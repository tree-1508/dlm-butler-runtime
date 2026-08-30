package main

import (
	"context"
	"fmt"
	"os"

	"github.com/tree-1508/dlm-butler-runtime/dldruntime"
	"github.com/tree-1508/dlm-butler-runtime/dldruntime/internal/buildinfo"
)

func main() {
	if len(os.Args) == 2 && (os.Args[1] == "-V" || os.Args[1] == "--version") {
		fmt.Println(buildinfo.VersionString())
		return
	}
	if !validDaemonArgs(os.Args[1:]) {
		fmt.Fprintln(os.Stderr, "DLD runtime only supports: daemon --json --transport stdio")
		os.Exit(2)
	}
	store, err := dldruntime.OpenPlatformStore(os.Getenv("DLD_PROVIDER_STATE_DIR"))
	if err != nil {
		fmt.Fprintln(os.Stderr, "provider state initialization failed")
		os.Exit(1)
	}
	runtime, err := dldruntime.New(dldruntime.Config{
		Store:       store,
		APIBaseURL:  "https://api.itch.io",
		StorageRoot: os.Getenv("DLD_STORAGE_ROOT"),
		StagingRoot: os.Getenv("DLD_STAGING_ROOT"),
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "runtime initialization failed")
		os.Exit(1)
	}
	if err := dldruntime.ServeStdio(context.Background(), os.Stdin, os.Stdout, runtime); err != nil {
		fmt.Fprintln(os.Stderr, "stdio runtime terminated with an error")
		os.Exit(1)
	}
}

func validDaemonArgs(args []string) bool {
	if len(args) != 4 || args[0] != "daemon" {
		return false
	}
	seenJSON := false
	seenTransport := false
	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--json":
			seenJSON = true
		case "--transport":
			if i+1 >= len(args) || args[i+1] != "stdio" {
				return false
			}
			seenTransport = true
			i++
		default:
			return false
		}
	}
	return seenJSON && seenTransport
}
