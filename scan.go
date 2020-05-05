package main

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/packages"
)

type set map[string]struct{}

func (s set) extend(o set) {
	for k, v := range o {
		s[k] = v
	}
}

var empty struct{}

// Scan a package and return a list of all the exported symbols in pkg which
// are from the package target.
func Scan(pkg *packages.Package, target string) set {
	symbols := make(set)
	for _, file := range pkg.Syntax {
		imports := importsIndex(file.Imports, pkg.Imports)
		ast.Inspect(file, func(n ast.Node) bool {
			sel, ok := n.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			id, ok := sel.X.(*ast.Ident)
			if !ok {
				return true
			}

			if id.Obj != nil {
				return true
			}

			if id.Name == imports[target] {
				symbols[sel.Sel.Name] = empty
			}
			return true
		})
	}
	return symbols
}

// importsIndex returns a map of import package paths to the local alias used for
// the import. Imports without an alias are looked up in the pkgImports.
func importsIndex(specs []*ast.ImportSpec, pkgImports map[string]*packages.Package) map[string]string {
	result := make(map[string]string, len(specs))
	for _, spec := range specs {
		pkgPath := unquote(spec.Path.Value)
		if spec.Name != nil {
			result[pkgPath] = spec.Name.Name
			continue
		}

		// Lookup package name from pkgImports
		result[pkgPath] = pkgImports[pkgPath].Name
	}
	return result
}

func unquote(v string) string {
	return strings.Trim(v, `"`)
}
