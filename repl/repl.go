package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/compiler"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"monkey/vm"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer, useInterpreter bool) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()
	macroEnv := object.NewEnvironment()

	constants := []object.Object{}
	symbolTable := compiler.NewSymbolTable()
	globals := make([]object.Object, vm.GlobalSize)

	for {
		fmt.Printf("%s", PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		if useInterpreter {
			// Interpreter
			evaluator.DefineMacros(program, macroEnv)
			expanded := evaluator.ExpandMacros(program, macroEnv)

			evaluated := evaluator.Eval(expanded, env)
			if evaluated != nil {
				io.WriteString(out, evaluated.Inspect())
				io.WriteString(out, "\n")
			}
		} else {
			// VM
			comp := compiler.NewWithState(symbolTable, constants)
			err := comp.Compile(program)
			if err != nil {
				fmt.Fprintf(out, "Woops! Compilation failed:\n%s\n", err)
				continue
			}
			machine := vm.NewWithState(comp.Bytecode(), globals)
			err = machine.Run()
			if err != nil {
				fmt.Fprintf(out, "Woops! Executing bytecode failed:\n%s\n", err)
				continue
			}
			lastPopped := machine.LastPoppedStackElem()
			io.WriteString(out, lastPopped.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
