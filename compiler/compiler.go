package compiler

import (
	"monkey/ast"
	"monkey/code"
)

type Compiler struct {
	instructions code.Instructions
	constatns    []object.Object
}

func New() *Compiler {
	return &Compiler{
		instructions: code.Instructions{},
		constants:    []object.Object{},
	}
}

func (c *Compiler) Compile(node ast.Node) error {
	return nil
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.instructions,
		Constants:    c.constatns,
	}
}

type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}