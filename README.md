# hs-cli
# Usage
```
NAME:
   hs-cli - CLI client for the HackerSchool API

USAGE:
   hs-cli [global options] command [command options]

VERSION:
   0.0.1

COMMANDS:
   mgetall    retrieve all members
   mget       retrieve information of a member
   mcreate    create member providing path to a json file
   mupdate    update member information
   mremove    remove member from the database
   mprojects  get the projects a member is in
   pgetall    retrieve all projects
   pget       retrieve information of a project
   pcreate    create a new project
   pupdate    update information of a project
   pdelete    delete project from the database
   pmembers   get members in a project
   login      login to the API, saving the cookie to the cookiejar
   help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config value, -f value      path to config file
   --root value, -r value        API root url        (overwrites file and HS_ROOT environment configs)
   --user value, -u value        username            (overwrites file and HS_USER environment configs)
   --password value, -p value    user password       (overwrites file and HS_PASSWORD environment configs)
   --cookie-jar value, -c value  cookie jar path     (overwrites file and HS_COOKIEJAR environment configs)
   --debug, -d                   log debug information to the console (default: false)
   --help, -h                    show help
   --version, -v                 print the version
```
# Configuration
If configuration values are not provided as CLI options the program will first attempt to load them from the `--config` option and, after that, look for the environment variables and overwrite any options set by the configuration file. This way the order of preference is CLI args < environment variables < config file.

Example `config.yaml` file:
```yml
root: https://api.hackerschool.dev
user: username 
password: password
cookiejar: ./cookiejar.json
```
Example `.env` file:
```sh
export HS_ROOT="https://api.hackerschool.dev"
export HS_USER="username"
export HS_PASSWORD="password"
export HS_HS_COOKIERJAR="/home/my/cookiejar.json"
```

# Examples

Example using CLI options:
```sh
hscli -r https://api.hackerschool.dev -u username -p password -c cookiejar.json -d login 
```
Example using `config.yaml` file:
```sh
hscli --config config.yaml -d login
```

Example using `.env` file:
```sh
source .env; hscli -d login
```
The `source` command only needs to be ran once per shell session.

# Development
Commands are expected to be ran a certain way to keep structure somewhat coherent.
To build new commands start by adding a new `Command` entry into the `urfave/cli/v2` object in `main.go`:
```go
{
    Name:      "newcommand",
    Usage:     "message displayed in -h flag of main program",
    UsageText: "message displayed in -h flag of command",
    Action: func(cCtx *cli.Context) error {
        if cCtx.Args().Len() == 0 {
            fmt.Fprintf(os.Stderr, "Missing <mandatory argument> argument\n")
            os.Exit(EX_USAGE) // exit code to specify CLI misusage 
        }
        os.Exit(commands.RunCommand(c, YourCommandHandlerFuncHere, cCtx.Args()...))
        return nil
    },
},
```
Now, in the `commands` package feel free to define this function anywhere you'd like, new files might be created if appropriate.
```go
func YourCommandHandlerFuncHere(c *client.Client, args ...string) ([]byte, error) {
    // You can validate arguments here
    if len(args) != 1 || !strings.HasPrefix(args[0], "https://") {
        return nil, NewCommandError("Invalid URL!", nil)
    }

    // Do your command logic here
    rsp, err := c.Http.Get(c.Cfg.Root + "/mycommandendpoint")
    if err != nil {
        return nil, NewCommandError("Failed requesting server", fmt.Errorf("http.Get %s: %w", "mycommandendpoint"; err)
    }

    // Here you can validate if any business logic error might have occured, for example you might
    // want to print something special for some HTTP codes, 404, 401, etc. 

    // Here you can return a generic domain error whose result will be printed to the user (considering the 
    // response JSON data is in the rspJson variable)
    if rsp.StatusCode != http.StatusOK {
        return nil, NewCommandError(fmt.Sprintf("%d %s\n%s", rsp.StatusCode, http.StatusText(rsp.StatusCode), string(rspJson)), nil)
    }

    // Return the raw JSON response from the server
    return rspData, nil
}
```

The `CommandError` has two usecases. 

One is to report generic errors not really related to the API domain, for example, lost connection, invalid URL provided by the user, etc. To report these errors just return a new instance with a descriptive message for the user and an error that wraps the original error, the "stack trace" will be printed if running in debug.

For errors related to business logic, for example, such as a user not being authenticated or a resource not existing, which don't have an underlying error caused at runtime, just create the `CommandErrror` with a `nil` cause parameter and the desired message to be displayed.

# TODO
- [ ] add usage documentation (similar to `-h` usage string but with examples)
- [ ] missing commands: `pcreate`, `pupdate`, `pdelete`
- [ ] development documentation (`command.go` has a comment with new command definition examples)
- [ ] define output format  (only printing json strings as of know)
