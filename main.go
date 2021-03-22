package main

import (
	"flag"
	"fmt"
	"monkey/repl"
	"os"
	"os/user"
)

func main() {
	interpreter := flag.Bool("interpreter", false, "use interpreter instead of VM")
	flag.Parse()
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s! This is the Monkey programming language!\n", user.Username)
	if *interpreter {
		fmt.Printf("Interpreter mode\n")
	} else {
		fmt.Printf("Bytecode/VM mode\n")
	}
	fmt.Printf("Feel free to type in commands\n")
	repl.Start(os.Stdin, os.Stdout, *interpreter)
}
