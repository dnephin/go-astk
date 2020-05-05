package main

import (
	"go/token"
	"testing"

	"golang.org/x/tools/go/packages"
	"gotest.tools/v3/assert"
)

func TestScan(t *testing.T) {
	cfg := &packages.Config{
		Mode: packages.LoadAllSyntax,
		Fset: token.NewFileSet(),
	}
	pkgs , err := packages.Load(cfg, "gotest.tools/v3/icmd")
	assert.NilError(t, err)

	symbols := Scan(pkgs[0], "time")
	t.Log(symbols)
	t.Fail()
}
