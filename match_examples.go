package asttt

import "go/ast"

const (
	pkgTestifyAssert       = "github.com/stretchr/testify/assert"
	pkgGopkgTestifyAssert  = "gopkg.in/stretchr/testify.v1/assert"
	pkgTestifyRequire      = "github.com/stretchr/testify/require"
	pkgGopkgTestifyRequire = "gopkg.in/stretchr/testify.v1/require"
	pkgAssert              = "gotest.tools/assert"
	pkgCmp                 = "gotest.tools/assert/cmp"
)

func matchTestifyTestingT() {
	pkgRename := &pkgRename{}
	SelectorExpr(
		pkgRename.Match,
		IdentName("TestingT"),
	)
}

type pkgRename struct {
	replacement string
}

func (r *pkgRename) Match(node ast.Node, matcher MatchContext) bool {
	ident, ok := node.(*ast.Ident)
	if !ok {
		return false
	}
	// TODO: better way to detect identifiers which are import names?
	if ident.Obj != nil {
		return false
	}
	switch ident.Name {
	case matcher.Imports[pkgTestifyAssert]:
		r.replacement = "Check"
	case matcher.Imports[pkgGopkgTestifyAssert]:
		r.replacement = "Check"
	default:
		return false
	}
	return true
}
