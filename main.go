package main

import (
	"fmt"
	"hscli/client"
	"hscli/commands"
	"hscli/config"
	"io"
	"log"
	"log/slog"
	"net/http"
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
							commands.GetMemberByUsername), cCtx.Args().Slice()...))
					return nil
				},
			},
			{
				Name:      "mcreate",
				Usage:     "create member",
				UsageText: "mcreate [commands options] [<file>]",
				Action: func(cCtx *cli.Context) error {
					os.Exit(commands.RunCommand(c,
						commands.WithLoginRetry(
							commands.DefaultLastArgumentToStdin(
								commands.CreateMember)), cCtx.Args().Slice()...))
					return nil
				},
			},
			{
				Name:      "mupdate",
				Usage:     "update member information",
				UsageText: "mupdate [command options] <username> [<file>]",
				Action: func(cCtx *cli.Context) error {
					if cCtx.Args().Len() < 1 {
						fmt.Fprintf(os.Stderr, "Missing arguments\n")
						os.Exit(EX_USAGE)
					}
					os.Exit(commands.RunCommand(c,
						commands.WithLoginRetry(
							commands.DefaultLastArgumentToStdin(
								commands.UpdateMember)), cCtx.Args().Slice()...))
					return nil
				},
			},
			{
				Name:      "mdelete",
				Usage:     "delete member from the database",
				UsageText: "mdelete [command options] <username>",
				Action: func(cCtx *cli.Context) error {
					if cCtx.Args().Len() < 1 {
						fmt.Fprintf(os.Stderr, "Missing <username> argument\n")
						os.Exit(EX_USAGE)
					}
					os.Exit(commands.RunCommand(c,
						commands.WithLoginRetry(
							commands.DeleteMember), cCtx.Args().Slice()...))
					return nil
				},
			},
			{
				Name:      "mprojects",
				Usage:     "get the projects a member is in",
				UsageText: "mprojects [command options] <username>",
				Action: func(cCtx *cli.Context) error {
					if cCtx.Args().Len() < 1 {
						fmt.Fprintf(os.Stderr, "Missing <username> argument\n")
						os.Exit(EX_USAGE)
					}
					os.Exit(commands.RunCommand(c,
						commands.WithLoginRetry(
							commands.GetMemberProjects), cCtx.Args().Slice()...))
					return nil
				},
			},
			{
				Name:      "mlogo",
				Usage:     "get member logo",
				UsageText: "mgetlogo [command options] <username>",
				Action: func(cCtx *cli.Context) error {
					if cCtx.Args().Len() < 1 {
						fmt.Fprintf(os.Stderr, "Missing <username> argument\n")
						os.Exit(EX_USAGE)
					}
					os.Exit(commands.RunCommand(c,
						commands.WithLoginRetry(
							commands.GetMemberLogo), cCtx.Args().Slice()...))
					return nil
				},
			},
			{
				Name:      "mtags",
				Usage:     "get member tags",
				UsageText: "mgetlogo [command options] <username>",
				Action: func(cCtx *cli.Context) error {
					if cCtx.Args().Len() == 0 {
						fmt.Fprintf(os.Stderr, "Missing <username> argument\n")
						os.Exit(EX_USAGE)
					}
					os.Exit(commands.RunCommand(c,
						commands.WithLoginRetry(
							commands.GetTags), cCtx.Args().Slice()...))
					return nil
				},
			},
			{
				Name:      "maddproject",
				Usage:     "add a project to a member",
				UsageText: "maddproject [commands options] <username> <proj_name> [<file>]",
				Action: func(cCtx *cli.Context) error {
					if cCtx.Args().Len() < 2 {
						fmt.Fprintf(os.Stderr, "Missing arguments\n")
						os.Exit(EX_USAGE)
					}
					os.Exit(commands.RunCommand(c,
						commands.WithLoginRetry(
							commands.DefaultLastArgumentToStdin(
								commands.AddProject)), cCtx.Args().Slice()...))
					return nil
				},
			},
			{
				Name:      "maddlogo",
				Usage:     "upload member logo",
				UsageText: "mupdatelogo [command options] <username> [<logo>]",
				Action: func(cCtx *cli.Context) error {
					if cCtx.Args().Len() < 1 {
						fmt.Fprintf(os.Stderr, "Missing <username> arguments\n")
						os.Exit(EX_USAGE)
					}
					os.Exit(commands.RunCommand(c,
						commands.WithLoginRetry(
							commands.DefaultLastArgumentToStdin(
								commands.UpdateMemberLogo)), cCtx.Args().Slice()...))
					return nil
				},
			},
			{
				Name:      "maddtag",
				Usage:     "add member tag",
				UsageText: "maddtag [command options] <username> [<file>]",
				Action: func(cCtx *cli.Context) error {
					if cCtx.Args().Len() < 1 {
						fmt.Fprintf(os.Stderr, "Missing <username> arguments\n")
						os.Exit(EX_USAGE)
					}
					os.Exit(commands.RunCommand(c,
						commands.WithLoginRetry(
							commands.DefaultLastArgumentToStdin(
								commands.AddTag)), cCtx.Args().Slice()...))
					return nil
				},
			},
			{
				Name:      "mdeltag",
				Usage:     "delete member tag",
				UsageText: "mdeltag [command options] <username> [<file>]",
				Action: func(cCtx *cli.Context) error {
					if cCtx.Args().Len() < 1 {
						fmt.Fprintf(os.Stderr, "Missing <username> arguments\n")
						os.Exit(EX_USAGE)
					}
					os.Exit(commands.RunCommand(c,
						commands.WithLoginRetry(
							commands.DefaultLastArgumentToStdin(
								commands.DeleteTag)), cCtx.Args().Slice()...))
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
				UsageText: "pget [command options] <proj_name>",
				Action: func(cCtx *cli.Context) error {
					if cCtx.Args().Len() < 1 {
						fmt.Fprintf(os.Stderr, "Missing <proj_name> argument\n")
						os.Exit(EX_USAGE)
					}
					os.Exit(commands.RunCommand(c,
						commands.WithLoginRetry(
							commands.GetProjectByID), cCtx.Args().Slice()...))
					return nil
				},
			},
			{
				Name:      "pcreate",
				Usage:     "create a new project",
				UsageText: "pcreate [command options] [<file>]",
				Action: func(cCtx *cli.Context) error {
					os.Exit(commands.RunCommand(c,
						commands.WithLoginRetry(
							commands.DefaultLastArgumentToStdin(
								commands.CreateProject)), cCtx.Args().Slice()...))
					return nil
				},
			},
			{
				Name:      "pupdate",
				Usage:     "update information of a project",
				UsageText: "pupdate [command options] <proj_name> [<file>]",
				Action: func(cCtx *cli.Context) error {
					if cCtx.Args().Len() < 1 {
						fmt.Fprintf(os.Stderr, "Missing arguments\n")
						os.Exit(EX_USAGE)
					}
					os.Exit(commands.RunCommand(c,
						commands.WithLoginRetry(
							commands.DefaultLastArgumentToStdin(
								commands.UpdateProject)), cCtx.Args().Slice()...))
					return nil
				},
			},
			{
				Name:      "pdelete",
				Usage:     "delete project from the database",
				UsageText: "pdelete [command options] <proj_name> ",
				Action: func(cCtx *cli.Context) error {
					if cCtx.Args().Len() < 1 {
						fmt.Fprintf(os.Stderr, "Missing <proj_name> argument\n")
						os.Exit(EX_USAGE)
					}
					os.Exit(commands.RunCommand(c,
						commands.WithLoginRetry(
							commands.DeleteProject), cCtx.Args().Slice()...))
					return nil
				},
			},
			{
				Name:      "pmembers",
				Usage:     "get members in a project",
				UsageText: "pmembers [command options] <proj_name>",
				Action: func(cCtx *cli.Context) error {
					if cCtx.Args().Len() < 1 {
						fmt.Fprintf(os.Stderr, "Missing <proj_name> argument\n")
						os.Exit(EX_USAGE)
					}
					os.Exit(commands.RunCommand(c,
						commands.WithLoginRetry(
							commands.GetProjectMembers), cCtx.Args().Slice()...))
					return nil
				},
			},
			{
				Name:      "plogo",
				Usage:     "get project logo",
				UsageText: "pmembers [command options] <proj_name>",
				Action: func(cCtx *cli.Context) error {
					if cCtx.Args().Len() < 1 {
						fmt.Fprintf(os.Stderr, "Missing <proj_name> argument\n")
						os.Exit(EX_USAGE)
					}
					os.Exit(commands.RunCommand(c,
						commands.WithLoginRetry(
							commands.GetProjectLogo), cCtx.Args().Slice()...))
					return nil
				},
			},
			{
				Name:      "paddmember",
				Usage:     "add member to a project",
				UsageText: "pmembers [command options] <proj_name> <username>",
				Action: func(cCtx *cli.Context) error {
					if cCtx.Args().Len() < 2 {
						fmt.Fprintf(os.Stderr, "Missing argument\n")
						os.Exit(EX_USAGE)
					}
					os.Exit(commands.RunCommand(c,
						commands.WithLoginRetry(
							commands.GetProjectMembers), cCtx.Args().Slice()...))
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
					os.Exit(commands.RunCommand(c, commands.Login))
					return nil
				},
			},
			{
				Name:  "logout",
				Usage: "logout off the API, clearing the session",
				Action: func(cCtx *cli.Context) error {
					os.Exit(commands.RunCommand(c,
						func(c *client.Client, args ...string) ([]byte, error) {
							rsp, err := c.Http.Get(c.Cfg.Root + "/logout")
							if err != nil {
								return nil, commands.NewCommandError("Failed requesting server", fmt.Errorf("http.Get %s: %w", c.Cfg.Root+"/logout", err))
							}
							defer rsp.Body.Close()
							rspData, err := io.ReadAll(rsp.Body)
							if err != nil {
								return nil, commands.NewCommandError("Failed receiving server response", fmt.Errorf("io.ReadAll: %w", err))
							}
							if rsp.StatusCode != http.StatusOK {
								return nil, commands.NewCommandError(fmt.Sprintf("%d %s\n%s", rsp.StatusCode, http.StatusText(rsp.StatusCode), string(rspData)), nil)
							}
							return rspData, nil
						}))
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
