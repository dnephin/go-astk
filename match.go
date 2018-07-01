package asttt

import (
	"go/ast"

	"golang.org/x/tools/go/ast/astutil"
)

type MatchNode func(node ast.Node, matcher MatchContext) bool

type ReplaceNode func(node ast.Node, matcher MatchContext) Action

type MatchContext struct {
	Imports map[string]string
}

type Result struct {
	matched bool
	root    ast.Node
	action  Action
}

func (r Result) Matched() bool {
	return r.matched
}

// TODO: remove action
func (r Result) Action() Action {
	return r.action
}

var NoMatch = Result{action: func(cursor *astutil.Cursor) {}}

// Action accepts a cursor and performs an operation (Delete, Replace, Insert)
// on the cursor.
type Action func(cursor *astutil.Cursor)

func ActionDelete(cursor *astutil.Cursor) {
	cursor.Delete()
}

func ActionReplace(node ast.Node) Action {
	return func(cursor *astutil.Cursor) {
		cursor.Replace(node)
	}
}

// TODO:
//func ActionInsertAfter
//func ActionInsertBefore

func SelectorExpr(matchX MatchNode, matchSel MatchNode) MatchNode {
	return func(node ast.Node, matcher MatchContext) bool {
		selExpr, ok := node.(*ast.SelectorExpr)
		if !ok {
			return false
		}
		return matchX(selExpr.X, matcher) && matchSel(selExpr.Sel, matcher)
	}
}

// PackageIdent matches an identifier that is an imported package. The package
// path is used to lookup the import alias to compare to the identifier name.
func PackageIdent(pkgPath string) MatchNode {
	return func(node ast.Node, matcher MatchContext) bool {
		ident, ok := node.(*ast.Ident)
		if !ok {
			return false
		}
		// TODO: better way to detect identifiers which are import names?
		if ident.Obj != nil {
			return false
		}
		return ident.Name == matcher.Imports[pkgPath]
	}
}

func IdentName(name string) MatchNode {
	return func(node ast.Node, matcher MatchContext) bool {
		ident, ok := node.(*ast.Ident)
		if !ok {
			return false
		}
		return ident.Name == name
	}
}

func Any() bool {

}
