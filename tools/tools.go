package tools

// This is a list of tools to be maintained; some non-main
// package in the same repo needs to be imported so they'll be managed in
// go.mod.

import (
	_ "github.com/golangci/golangci-lint/pkg/golinters"
	// _ "github.com/mailru/easyjson"
	// _ "github.com/valyala/quicktemplate"
	// _ "golang.org/x/tools/go/packages"
	_ "golang.org/x/tools/imports"
)
