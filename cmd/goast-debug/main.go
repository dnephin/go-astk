package main

import (
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/token"
	"os"
	"path/filepath"

	"golang.org/x/tools/go/loader"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error: "+err.Error())
		os.Exit(1)
	}
}

func run() error {
	flags, err := setupFlags(os.Args)
	switch {
	case err != nil:
		return err
	case len(flags.Args()) == 0:
		flags.Usage()
		return errors.New("missing position args")
	}

	config := loader.Config{}
	_, err = config.FromArgs(flags.Args(), true)
	if err != nil {
		return err
	}

	program, err := config.Load()
	if err != nil {
		return err
	}

	if len(program.InitialPackages()) == 0 {
		return errors.New("no package found")
	}

	for _, pkg := range program.InitialPackages() {
		for _, file := range pkg.Files {
			fmt.Printf("--> %s/%s\n",
				pkg.Pkg.Name(), filename(program.Fset, file))
			_ = ast.Print(program.Fset, file)
			fmt.Println()
		}
	}
	return nil
}

func filename(fileset *token.FileSet, file *ast.File) string {
	return filepath.Base(fileset.File(file.Pos()).Name())
}

func setupFlags(args []string) (*flag.FlagSet, error) {
	flags := flag.NewFlagSet(args[0], flag.ContinueOnError)
	flags.Usage = func() {
		fmt.Fprintf(os.Stderr,
			"Usage: %s <args>\n\n%s", args[0], loader.FromArgsUsage)
	}
	err := flags.Parse(args[1:])
	return flags, err
}
