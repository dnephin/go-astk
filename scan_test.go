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
	pkgs, err := packages.Load(cfg, "github.com/dnephin/astprune/internal/consumer")
	assert.NilError(t, err)

	symbols := Scan(pkgs[0], "github.com/dnephin/astprune/internal/target")
	expected := set{
		"TypeSymbol1": empty,
		"TypeSymbol2": empty,
		"TypeSymbol3": empty,
		"TypeSymbol4": empty,
		"TypeSymbol5": empty,
		"TypeSymbol6": empty,
		"TypeSymbol7": empty,
		"TypeSymbol8": empty,
		"FuncSymbol":  empty,
		"VarSymbol1":  empty,
		"VarSymbol2":  empty,
		"ConstSymbol": empty,
	}
	assert.DeepEqual(t, symbols, expected)
}
