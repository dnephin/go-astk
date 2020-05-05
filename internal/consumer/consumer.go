package consumer

import (
	"fmt"

	"github.com/dnephin/astprune/internal/target"
)

func symbolInReturnValue() *target.TypeSymbol1 {
	return nil
}

func symbolInFuncArg(_ target.TypeSymbol2) {}

type foo struct {
	target.TypeSymbol3

	symbolInStructField target.TypeSymbol4
}

var symbolInPkgVar *target.TypeSymbol5

func symbolInFunctionBody() {
	fmt.Println(target.VarSymbol1)
	target.FuncSymbol()

	_ = target.VarSymbol2
	_ = target.ConstSymbol

	var _ target.TypeSymbol8
}

type foop interface {
	SymbolInInterfaceMethodReturn() target.TypeSymbol6
	SymbolInInterfaceMethodArg(target.TypeSymbol7)
}
