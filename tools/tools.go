package tools

// This is a list of tools to be maintained; some non-main
// package in the same repo needs to be imported so they'll be managed in
// go.mod.

import (
	_ "github.com/golangci/golangci-lint/pkg/golinters" // for golangci-lint
	_ "golang.org/x/tools/imports"                      // for goimports
	// _ "github.com/mailru/easyjson"  // for easyjson
	// _ "github.com/valyala/quicktemplate"  // for qtc
	// _ "golang.org/x/tools/go/packages"  // for stringer
)
