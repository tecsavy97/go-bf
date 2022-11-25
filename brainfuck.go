package main

import (
	"fmt"
	"log"
	"os"

	"github.com/tecsavy97/go-bf/brainfuck"
	"github.com/urfave/cli"
)

func main() {
	app := &cli.App{
		Name:  "Brainfuck Interpreter",
		Usage: "A Brainfuck cli interpreter",
		Action: func(c *cli.Context) error {
			if len(c.Args()) > 0 {
				file, err := os.Open(c.Args().Get(0))
				if err != nil {
					log.Fatal(err)
				}
				defer file.Close()
				fmt.Print(brainfuck.Interpret(file))
			} else {
				log.Fatal("Fatal error: No input file\n")
			}
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
