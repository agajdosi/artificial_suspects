package main

// JUST FOR TESTING NOW
// GONNA BE A HELPER BINARY TO GENERATE AI DESCRIPTIONS ETC...

import (
	"fmt"
	"log"
	"os"
	"suspects/database"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "description",
				Aliases: []string{"d"},
				Usage:   "Generate description for image",
				Action:  generateDescription,
			},
			{
				Name:    "complete",
				Aliases: []string{"c"},
				Usage:   "complete a task on the list",
				Action: func(cCtx *cli.Context) error {
					fmt.Println("completed task: ", cCtx.Args().First())
					return nil
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func generateDescription(cCtx *cli.Context) error {
	fmt.Println("Description for: ", cCtx.Args().First())
	database.EnsureDBAvailable()

	service, err := database.GetService("OpenAI")
	if err != nil {
		return err
	}

	fmt.Println("Token:", service.Token)

	return nil
}
