package main

import (
	"fmt"
	"hscli/client"
	"hscli/commands"
	"hscli/config"
	"hscli/logging"
	"log"
	"log/slog"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	c := client.NewClient()
	app := &cli.App{
		Name:                 client.ProgramName,
		Version:              client.ProgramVersion,
		Usage:                "CLI client for the HackerSchool API",
		EnableBashCompletion: true,
		// Load config from file or environment
		Before: func(cCtx *cli.Context) error {
			if cCtx.Bool("debug") {
				slog.SetLogLoggerLevel(slog.LevelDebug)
			}
			if err := config.LoadConfig(cCtx, c.Cfg, cCtx.String("config")); err != nil {
				return err
			}
			c.SetupJar()
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"f"},
				Value:   "",
				Usage:   "path to config file",
			},
			&cli.StringFlag{
				Name:        "root",
				Aliases:     []string{"r"},
				Value:       "",
				Usage:       "API root url        (overwrites file and HS_ROOT environment configs)",
				Destination: &c.Cfg.Root,
			},
			&cli.StringFlag{
				Name:        "user",
				Aliases:     []string{"u"},
				Value:       "",
				Usage:       "username            (overwrites file and HS_USER environment configs)",
				Destination: &c.Cfg.User,
			},
			&cli.StringFlag{
				Name:        "password",
				Aliases:     []string{"p"},
				Value:       "",
				Usage:       "user password       (overwrites file and HS_PASSWORD environment configs)",
				Destination: &c.Cfg.Password,
			},
			&cli.StringFlag{
				Name:        "cookie-jar",
				Aliases:     []string{"c"},
				Value:       "",
				Usage:       "cookie jar path     (overwrites file and HS_COOKIEJAR environment configs)",
				Destination: &c.Cfg.CookieJarPath,
			},

			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
				Usage:   "log debug information to the console",
			},
		},
		Action: func(cCtx *cli.Context) error {
			if cCtx.Args().Len() == 0 {
				return cli.ShowAppHelp(cCtx)
			}
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:      "mgetall",
				Usage:     "retrieve all members",
				UsageText: "mgetall [command options]",
				Action: func(cCtx *cli.Context) error {
					return commands.RunCommand(c,
						commands.WithLoginRetry(
							commands.GetMembers))
				},
			},
			{
				Name:      "mget",
				Usage:     "retrieve information of a member",
				UsageText: "mget [command options] <username>",
				Action: func(cCtx *cli.Context) error {
					if cCtx.Args().Len() == 0 {
						fmt.Println("Missing argument <username>")
						return nil
					}
					return commands.RunCommand(c,
						commands.WithLoginRetry(
							commands.GetMemberByUsername), cCtx.Args().First())
				},
			},
			{
				Name:      "mcreate",
				Usage:     "create member providing path to a json file",
				UsageText: "mcreate [commands options] <file>",
				Action: func(cCtx *cli.Context) error {
					if cCtx.Args().Len() == 0 {
						fmt.Println("Missing argument <file>")
						return nil
					}
					return commands.RunCommand(c,
						commands.WithLoginRetry(
							commands.CreateMember), cCtx.Args().First())
				},
			},
			{
				Name:  "mupdate",
				Usage: "update member information",
			},
			{
				Name:  "mremove",
				Usage: "remove member from the database",
			},
			{
				Name:  "mprojects",
				Usage: "get the projects a member is in",
			},
			{
				Name:      "pgetall",
				Usage:     "retrieve all projects",
				UsageText: "pgetall [command options]",
				Action: func(cCtx *cli.Context) error {
					return commands.RunCommand(c,
						commands.WithLoginRetry(
							commands.GetProjects))
				},
			},
			{
				Name:      "pget",
				Usage:     "retrieve information of a project",
				UsageText: "pget [command options] <id>",
				Action: func(cCtx *cli.Context) error {
					if cCtx.Args().Len() == 0 {
						fmt.Println("Missing argument <id>")
						return nil
					}
					return commands.RunCommand(c,
						commands.WithLoginRetry(
							commands.GetProjectByID), cCtx.Args().First())
				},
			},
			{
				Name:  "pupdate",
				Usage: "update information of a project",
			},
			{
				Name:  "pdelete",
				Usage: "delete project from the database",
			},
			{
				Name:  "pmembers",
				Usage: "get members in a project",
			},
			{
				Name:  "login",
				Usage: "login to the API, saving the cookie to the cookiejar",
				Action: func(cCtx *cli.Context) error {
					if err := c.Login(); err != nil {
						logging.LogDebug("%s", err)
						fmt.Println("Couldn't log in! Turn debug for more information")
						return nil
					}
					fmt.Println("Logged in sucessfully!")
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
