package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error: "+err.Error())
		os.Exit(1)
	}
}

func run() error {
	opts, err := setupFlags(os.Args)
	switch {
	case err == flag.ErrHelp:
		return nil
	case err != nil:
		return err
	}

	fileset := token.NewFileSet()
	f, err := parser.ParseFile(fileset, opts.path, nil, parser.AllErrors|opts.mode)
	if err != nil {
		return err
	}

	if opts.function != "" {
		return printFunctionAST(opts.function, fileset, f)
	}

	return ast.Print(fileset, f)
}

func printFunctionAST(name string, fileset *token.FileSet, file *ast.File) error {
	for _, decl := range file.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if funcDecl.Name == nil || funcDecl.Name.Name != name {
			continue
		}
		return ast.Print(fileset, funcDecl)
	}
	return fmt.Errorf("failed to find function %v in %v", name, file.Name.Name)
}

type options struct {
	path     string
	function string
	mode     parser.Mode
}

func setupFlags(args []string) (options, error) {
	o := options{}
	name, args := args[0], args[1:]
	flags := flag.NewFlagSet(name, flag.ContinueOnError)
	flags.Usage = usage(name)
	flags.Var(&parserModeFlagValue{value: &o.mode, field: parser.ParseComments}, "comments",
		"include comments in the printed AST")

	err := flags.Parse(args)
	if err != nil {
		return options{}, err
	}

	o.path = flags.Arg(0)
	o.function = flags.Arg(1)

	if o.path == "" {
		flags.Usage()
		return o, fmt.Errorf("missing required positional arg FILENAME")
	}
	return o, nil
}

func usage(name string) func() {
	return func() {

		fmt.Fprintf(os.Stderr, `
Print the AST source for the file or function in the file.   

Usage:
    %s FILENAME [FUNCTION]

`, name)
	}
}

type parserModeFlagValue struct {
	field parser.Mode
	value *parser.Mode
}

func (p *parserModeFlagValue) String() string {
	return "bool"
}

func (p *parserModeFlagValue) Set(_ string) error {
	*p.value |= p.field
	return nil
}

func (p *parserModeFlagValue) IsBoolFlag() bool {
	return true
}

var _ flag.Value = (*parserModeFlagValue)(nil)
