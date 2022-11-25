package brainfuck

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/tecsavy97/go-bf/stack"
)

// ProgramLength to define the length that for the brainfuck program to Run
// taken from the documentation
const ProgramLength = 30000

var (
	//the data on which the code will execute
	program = make([]uint32, ProgramLength)
	//Posotion of pointer on the program
	pointer = 0
	//Main Stack to keep track of '{' and '}' operands
	mainStack stack.Stack
	//Loopstack to keep track of loops
	loopStack stack.Stack

	// Customer Commands shouldRedirect
	customCommandList = make(customCommands, 0)
)

const (
	//Error
	CustommCommandAlreadyPresentErr = "Custom Command already exists"
	CustomCommandCannotExistInASCII = "Custom Command cannot exist in the ASCII Table"
	CustomCommandAbsentErr          = "Custom Command is not present"
)

// the customer commands need to have
type CustomerOperands struct {
	Operands  int
	Operation func(uint32) uint32
}

type customCommands map[int]CustomerOperands

// Returns New CustomCommand Object to add or remove
func NewCustomCommandList() customCommands {
	return customCommandList
}

// Receiver function to add custom AddCustomCommand
// it will also check whether command can be used for the library or not
func (c customCommands) AddCustomCommand(custom CustomerOperands) (err error) {
	switch custom.Operands {
	case '>', '<', '+', '-', '.', ',', '[', '}':
		err = errors.New("Existing Operands cant be added")
		return
	}
	if _, ok := c[custom.Operands]; ok {
		return errors.New(CustommCommandAlreadyPresentErr)
	}
	if custom.Operands > 127 {
		return errors.New(CustomCommandCannotExistInASCII)
	}
	customCommandList[custom.Operands] = custom
	return
}

// Removes existing customer command
// It will give error of the command is not present in the
func (c customCommands) RemoveCustomCommand(command int) (err error) {
	if _, ok := c[command]; !ok {
		return errors.New(CustomCommandAbsentErr)
	}
	delete(customCommandList, command)
	return
}

// Give a pretty print og the custom commands if present
func (c customCommands) PreviewAllCustomCommands() {
	if len(customCommandList) > 0 {
		list, _ := json.MarshalIndent(customCommandList, "", "\t")
		fmt.Println(string(list))
	}
}

// Executes the operands that are passed by the interpreter
func execute(op string, index *int, program *[]uint32) (result string) {
	switch []byte(op)[0] {
	// Move the pointer to the right
	case '>':
		if *index == ProgramLength-1 {
			break
		}
		(*index)++
	// Move the pointer to the left
	case '<':
		if *index == 0 {
			*index = ProgramLength - 1
			break
		}
		(*index)--
	// Increment the memory cell under the pointer
	case '+':
		(*program)[*index]++
	// Decrement the memory cell under the pointer
	case '-':
		if (*program)[*index] == 0 {
			(*program)[*index] = 255
			break
		}
		(*program)[*index]--
	// Output the character signified by the cell at the pointer
	case '.':
		result = fmt.Sprintf("%c", (*program)[*index])
	// Input a character and store it in the cell at the pointer
	case ',':
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Text() != "" {
			input, err := strconv.ParseUint(scanner.Text(), 10, 32)
			if err != nil {
				log.Fatal(err)
			}
			(*program)[*index] = uint32(input)
		}
	default:
		if len(customCommandList) > 0 {
			opInt := []rune(op)
			if v, ok := customCommandList[int(opInt[0])]; ok {
				if v.Operation != nil {
					(*program)[*index] = v.Operation((*program)[*index])
				}
			}
		}
	}
	return
}

func Interpret(stream io.Reader) (result string) {
	var op string
	buf := make([]byte, 1)
	for {
		if loopStack.Len() > 0 {
			op = loopStack.Pop()
		} else {
			_, err := io.ReadFull(stream, buf)
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Fatal(err)
				break
			}
			op = string(buf)
		}
		switch []byte(op)[0] {
		case '>', '<', '+', '-', '.', ',':
			result += execute(op, &pointer, &program)
			mainStack.Push(op)
			break
		case '[':
			mainStack.Push(op)
			break
		case ']':
			mainStack.Push(op)
			if program[pointer] > 0 {
				innerLoop := 0
				firsttimehit := false
				for {
					operation := mainStack.Pop()
					if operation == "" {
						break
					}
					loopStack.Push(operation)
					// nested loops
					if operation == "]" && firsttimehit {
						innerLoop++
					}
					if operation == "[" {
						if innerLoop == 0 {
							break
						} else {
							innerLoop--
						}
					}
					firsttimehit = true
				}
			}
		default:
			if len(customCommandList) > 0 {
				opInt := []rune(op)
				if _, ok := customCommandList[int(opInt[0])]; ok {
					result += execute(op, &pointer, &program)
					mainStack.Push(op)
					break
				}
			}
		}
	}
	return
}
