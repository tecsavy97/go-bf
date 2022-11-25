package main

import (
	"strings"
	"testing"

	"github.com/tecsavy97/go-bf/brainfuck"
)

func TestHelloWorldRunBf(t *testing.T) {
	res := brainfuck.Interpret(strings.NewReader(">++++++++[<+++++++++>-]<.>++++[<+++++++>-]<+.+++++++..+++.>>++++++[<+++++++>-]<++.------------.>++++++[<+++++++++>-]<+.<.+++.------.--------.>>>++++[<++++++++>-]<+."))
	if res != "Hello, World!" {
		t.Fail()
	}
}

// Nth fibonacci number starts with 0,1
func fib(n uint32) uint32 {
	if n == 1 {
		return 1
	} else if n == 0 {
		return 0
	}
	return fib(n-1) + fib(n-2)
}

func TestCustomCommandRunBf(t *testing.T) {
	//Creates a new command List
	customCommands := brainfuck.NewCustomCommandList()
	//add current cells fibonacci number by 'f'
	cmd := brainfuck.CustomerOperands{
		Operands:  'f',
		Operation: fib,
	}
	err := customCommands.AddCustomCommand(cmd)
	if err != nil {
		t.Fatalf(err.Error())
	}
	res := brainfuck.Interpret(strings.NewReader("+++++++++++f--."))
	if res != "W" {
		t.Fail()
	}

}

func TestCustomCommandFaultyCharRunBf(t *testing.T) {
	customCommands := brainfuck.NewCustomCommandList()
	cmd := brainfuck.CustomerOperands{
		Operands:  'ƒ',
		Operation: fib,
	}
	err := customCommands.AddCustomCommand(cmd)
	if err.Error() != brainfuck.CustomCommandCannotExistInASCII {
		t.Fail()
	}
}
