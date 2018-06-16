package main

import (
	"log"
	"os"

	"github.com/UnnoTed/authenticaTed/generator"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "authenticaTed"
	app.Usage = "generates a auth system for the specified type"
	app.Action = func(c *cli.Context) error {
		log.Println("Running...")
		err := generator.Run()
		if err != nil {
			log.Fatal(err)
		}

		return nil
	}

	app.Run(os.Args)
}
