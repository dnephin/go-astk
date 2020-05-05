package main

import (
	"flag"
	"fmt"
	"go/token"
	"os"
	"strings"

	"golang.org/x/tools/go/packages"
)

func main() {
	flags, opts := setupFlags(os.Args[0])
	switch err := flags.Parse(os.Args[1:]); {
	case err == flag.ErrHelp:
		return
	case err != nil:
		flags.Usage()
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(3)
	}

	opts.source = flags.Args()
	if err := run(opts); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func setupFlags(name string) (*flag.FlagSet, *options) {
	opts := &options{}
	flags := flag.NewFlagSet(name, flag.ContinueOnError)
	flags.StringVar(&opts.target, "target", "", "package to prune")
	flags.BoolVar(&opts.dryrun, "dry-run", false, "print symbols instead of pruning")
	return flags, opts
}

type options struct {
	source []string
	target string
	dryrun bool
}

func run(opts *options) error {
	cfg := &packages.Config{
		Mode: packages.LoadAllSyntax,
		Fset: token.NewFileSet(),
	}
	pkgs, err := packages.Load(cfg, opts.source...)
	if err != nil {
		return err
	}

	symbols := make(set)
	for _, pkg := range pkgs {
		if len(pkg.Errors) > 0 {
			return errPkgLoad(pkg)
		}

		symbols.extend(Scan(pkg, opts.target))
	}

	for symbol := range symbols {
		fmt.Println(symbol)
	}
	return nil
}

func errPkgLoad(pkg *packages.Package) error {
	buf := new(strings.Builder)
	for _, err := range pkg.Errors {
		buf.WriteString("\n" + err.Error())
	}
	return fmt.Errorf("failed to load package %v %v", pkg.PkgPath, buf.String())
}
