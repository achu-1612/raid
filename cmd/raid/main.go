package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/achu-1612/raid"
	"github.com/urfave/cli/v3"

	"github.com/jedib0t/go-pretty/v6/table"
)

func main() {
	cmd := &cli.Command{
		Commands:    getCommands(),
		Name:        "raid",
		Usage:       "simulate RAID operations",
		Description: "A command-line tool to simulate RAID operations and manage tasks.",
		Authors:     []any{"aka.achu.1612@gmail.com"},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func getCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:   "init",
			Usage:  "initialize RAID state",
			Action: raidInitAction,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "dir",
					Aliases: []string{"d"},
					Usage:   "path to the base directory",
					Value:   ".",
				},
			},
		},
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "list raid configurations",
			Action:  raidListAction,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "dir",
					Aliases: []string{"d"},
					Usage:   "path to the base directory",
					Value:   ".",
				},
			},
		},
		{
			Name:   "create",
			Usage:  "create a new RAID configuration",
			Action: raidCreateAction,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "type",
					Aliases:  []string{"t"},
					Usage:    "raid type (e.g. RAID0, RAID1, RAID5, RAID10)",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "name",
					Aliases:  []string{"n"},
					Usage:    "raid name",
					Required: true,
				},
				&cli.StringFlag{
					Name:    "dir",
					Aliases: []string{"d"},
					Usage:   "path to the base directory",
					Value:   ".",
				},
			},
		},
		// {
		// 	Name:    "template",
		// 	Aliases: []string{"t"},
		// 	Usage:   "options for task templates",
		// 	Commands: []*cli.Command{
		// 		{
		// 			Name:  "add",
		// 			Usage: "add a new template",
		// 			Action: func(ctx context.Context, cmd *cli.Command) error {
		// 				fmt.Println("new task template: ", cmd.Args().First())
		// 				return nil
		// 			},
		// 		},
		// 		{
		// 			Name:  "remove",
		// 			Usage: "remove an existing template",
		// 			Action: func(ctx context.Context, cmd *cli.Command) error {
		// 				fmt.Println("removed task template: ", cmd.Args().First())
		// 				return nil
		// 			},
		// 		},
		// 	},
		// },
	}
}

func raidListAction(ctx context.Context, cmd *cli.Command) error {
	dir := cmd.String("dir")

	raids, err := raid.LoadRAIDState(dir)
	if err != nil {
		return err
	}

	if len(raids) == 0 {
		fmt.Println("No RAID configurations found.")

		return nil
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleColoredBright)
	t.AppendHeader(table.Row{"#", "Name", "Type", "Drives"})

	i := 1

	for _, raid := range raids {
		t.AppendRow(table.Row{
			i,
			raid.Name,
			raid.RAIDType,
			strings.Join(raid.Drives, ", "),
		})

		i++
	}

	t.Render()

	return nil
}

func raidInitAction(ctx context.Context, cmd *cli.Command) error {
	dir := cmd.String("dir")
	if err := raid.InitializeState(dir); err != nil {
		return err
	}

	fmt.Printf("Initialized RAID state in directory: %s\n", dir)

	return nil
}

func raidCreateAction(ctx context.Context, cmd *cli.Command) error {
	dir := cmd.String("dir")
	raidType := raid.RAIDType(cmd.String("type"))
	name := cmd.String("name")

	if !raidType.IsValid() {
		return fmt.Errorf("invalid RAID type: %s", raidType)
	}

	newRaid, err := raid.New(dir, raidType, name)
	if err != nil {
		return fmt.Errorf("failed to create RAID: %w", err)
	}

	fmt.Printf("Created new RAID configuration: %s of type %s\n", newRaid.Name, newRaid.Type())

	return nil
}
