package main

import (
	"fmt"
	"hscli/client"
	"hscli/commands"
	"hscli/config"
	"log"
	"log/slog"
	"os"

	"github.com/urfave/cli/v2"
)

const EX_USAGE = 64 // https://stackoverflow.com/questions/1101957/are-there-any-standard-exit-status-codes-in-linux

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
			if err := config.LoadConfig(c.Cfg, cCtx.String("config")); err != nil {
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
		Commands: []*cli.Command{
			{
				Name:      "mgetall",
				Usage:     "retrieve all members",
				UsageText: "mgetall [command options]",
				Action: func(cCtx *cli.Context) error {
					os.Exit(commands.RunCommand(c,
						commands.WithLoginRetry(
							commands.GetMembers)))
					return nil
				},
			},
			{
				Name:      "mget",
				Usage:     "retrieve information of a member",
				UsageText: "mget [command options] <username>",
				Action: func(cCtx *cli.Context) error {
					if cCtx.Args().Len() == 0 {
						fmt.Fprintf(os.Stderr, "Missing <username> argument\n")
						os.Exit(EX_USAGE)
					}
					os.Exit(commands.RunCommand(c,
						commands.WithLoginRetry(
							commands.GetMemberByUsername), cCtx.Args().First()))
					return nil
				},
			},
			{
				Name:      "mcreate",
				Usage:     "create member providing path to a json file",
				UsageText: "mcreate [commands options] <file>",
				Action: func(cCtx *cli.Context) error {
					if cCtx.Args().Len() == 0 {
						fmt.Fprintf(os.Stderr, "Missing <file> argument\n")
						os.Exit(EX_USAGE)
					}
					os.Exit(commands.RunCommand(c,
						commands.WithLoginRetry(
							commands.CreateMember), cCtx.Args().First()))
					return nil
				},
			},
			{
				Name:      "mupdate",
				Usage:     "update member information",
				UsageText: "mupdate [command options] <username> <file>",
				Action: func(cCtx *cli.Context) error {
					if cCtx.Args().Len() != 2 {
						fmt.Fprintf(os.Stderr, "Missing arguments\n")
						os.Exit(EX_USAGE)
					}
					os.Exit(commands.RunCommand(c,
						commands.WithLoginRetry(
							commands.UpdateMember), cCtx.Args().Get(0), cCtx.Args().Get(1)))
					return nil
				},
			},
			{
				Name:      "mremove",
				Usage:     "remove member from the database",
				UsageText: "mremove [command options] <username>",
				Action: func(cCtx *cli.Context) error {
					if cCtx.Args().Len() == 0 {
						fmt.Fprintf(os.Stderr, "Missing <username> argument\n")
						os.Exit(EX_USAGE)
					}
					os.Exit(commands.RunCommand(c,
						commands.WithLoginRetry(
							commands.RemoveMember), cCtx.Args().First()))
					return nil
				},
			},
			{
				Name:      "mprojects",
				Usage:     "get the projects a member is in",
				UsageText: "mprojects [command options] <username>",
				Action: func(cCtx *cli.Context) error {
					if cCtx.Args().Len() == 0 {
						fmt.Fprintf(os.Stderr, "Missing <username> argument\n")
						os.Exit(EX_USAGE)
					}
					os.Exit(commands.RunCommand(c,
						commands.WithLoginRetry(
							commands.GetMemberProjects), cCtx.Args().First()))
					return nil
				},
			},
			{
				Name:      "pgetall",
				Usage:     "retrieve all projects",
				UsageText: "pgetall [command options]",
				Action: func(cCtx *cli.Context) error {
					os.Exit(commands.RunCommand(c,
						commands.WithLoginRetry(
							commands.GetProjects)))
					return nil
				},
			},
			{
				Name:      "pget",
				Usage:     "retrieve information of a project",
				UsageText: "pget [command options] <id>",
				Action: func(cCtx *cli.Context) error {
					if cCtx.Args().Len() == 0 {
						fmt.Fprintf(os.Stderr, "Missing <id> argument\n")
						os.Exit(EX_USAGE)
					}
					os.Exit(commands.RunCommand(c,
						commands.WithLoginRetry(
							commands.GetProjectByID), cCtx.Args().First()))
					return nil
				},
			},
			{
				Name:      "pcreate",
				Usage:     "create a new project",
				UsageText: "pcreate [command options] <file>",
				Action: func(cCtx *cli.Context) error {
					if cCtx.Args().Len() == 0 {
						fmt.Fprintf(os.Stderr, "Missing <file> argument\n")
						os.Exit(EX_USAGE)
					}
					os.Exit(commands.RunCommand(c,
						commands.WithLoginRetry(
							commands.CreateProject), cCtx.Args().Get(0), cCtx.Args().Get(1)))
					return nil
				},
			},
			{
				Name:      "pupdate",
				Usage:     "update information of a project",
				UsageText: "pupdate [command options] <id> <file>",
				Action: func(cCtx *cli.Context) error {
					if cCtx.Args().Len() != 2 {
						fmt.Fprintf(os.Stderr, "Missing arguments\n")
						os.Exit(EX_USAGE)
					}
					os.Exit(commands.RunCommand(c,
						commands.WithLoginRetry(
							commands.UpdateProject), cCtx.Args().Get(0), cCtx.Args().Get(1)))
					return nil
				},
			},
			{
				Name:      "pdelete",
				Usage:     "delete project from the database",
				UsageText: "pdelete [command options] <id> ",
				Action: func(cCtx *cli.Context) error {
					if cCtx.Args().Len() == 0 {
						fmt.Fprintf(os.Stderr, "Missing <username> argument\n")
						os.Exit(EX_USAGE)
					}
					os.Exit(commands.RunCommand(c,
						commands.WithLoginRetry(
							commands.DeleteProject), cCtx.Args().First()))
					return nil
				},
			},
			{
				Name:      "pmembers",
				Usage:     "get members in a project",
				UsageText: "pmembers [command optinos] <id>",
				Action: func(cCtx *cli.Context) error {
					if cCtx.Args().Len() == 0 {
						fmt.Fprintf(os.Stderr, "Missing <id> argument\n")
						os.Exit(EX_USAGE)
					}
					os.Exit(commands.RunCommand(c,
						commands.WithLoginRetry(
							commands.GetProjectMembers), cCtx.Args().First()))
					return nil
				},
			},
			{
				Name:  "login",
				Usage: "login to the API, saving the cookie to the cookiejar",
				Action: func(cCtx *cli.Context) error {
					// Here we are simply making use of the way this "framework" is setup,
					// instead of writting a new command (which would just result in code duplication)
					// we simply pass it a fake command which returns Unauthorized at first and forces the
					// decorator to attempt a login, if it can do it, then we just return successful
					var retried bool = false
					os.Exit(commands.RunCommand(c, commands.WithLoginRetry(
						func(c *client.Client, args ...string) ([]byte, error) {
							if !retried {
								retried = true
								return nil, client.ErrUnauthorized
							} else {
								return []byte("Logged in successfully!\n"), nil
							}
						},
					)))
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
