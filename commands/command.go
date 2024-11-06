package commands

import (
	"errors"
	"fmt"
	"hscli/client"
	"hscli/logging"
	"os"
)

type Command func(c *client.Client, args ...string) ([]byte, error)

type CommandError struct {
	Cause   error  // the cause of the error
	Message string // message to be displayed to the console if the error occurs
}

// Create new CommandError with a display message and an error cause, cause might be set to nil
func NewCommandError(message string, cause error) error {
	return CommandError{
		Message: message,
		Cause:   cause,
	}
}

func (e CommandError) Unwrap() error {
	return e.Cause
}

func (e CommandError) Error() string {
	return e.Message
}

func WithLoginRetry(cmd Command) Command {
	return func(c *client.Client, args ...string) ([]byte, error) {
		rsp, err := cmd(c, args...) // run command
		if err != nil {             // command fails, see if unauthorized error
			var cmdErr CommandError
			if errors.As(err, &cmdErr) {
				if errors.Is(cmdErr.Cause, client.ErrUnauthorized) {
					_, err := Login(c) // attempt to login
					if err != nil {    // login fails
						return nil, err
					}
					return cmd(c, args...)
				}
			} else {
				// shouldn't happen, all commands should return CommandError
				return nil, err
			}
		}
		return rsp, err
	}
}

func DefaultLastArgumentToStdin(cmd Command) Command {
	return func(c *client.Client, args ...string) ([]byte, error) {
		if len(args) == 0 {
			args = append(args, "/dev/stdin")
			return cmd(c, args...)
		}

		// if last argument is not an existing file default to stdin
		lastArg := args[len(args)-1]
		if _, err := os.Stat(lastArg); err != nil {
			args = append(args, "/dev/stdin")
			return cmd(c, args...)
		}

		return cmd(c, args...)
	}
}

// Runs a command.
// Returns 0 on success, 1 on an domain related errors such as (unauthorized, resource doesn't exist, etc) and 2 on generic errors (no connection, etc)
func RunCommand(c *client.Client, cmd Command, args ...string) int {
	r, err := cmd(c, args...)
	if err != nil {
		var commandErr CommandError
		if errors.As(err, &commandErr) {
			if commandErr.Cause == nil || errors.Is(commandErr.Cause, client.ErrUnauthorized) { // business logic error
				fmt.Fprintf(os.Stdout, "%s\n", err)
				return 1
			} else {
				logging.LogDebug(commandErr.Cause.Error())
			}
		}

		fmt.Fprintf(os.Stderr, "%s\n", err)
		return 2
	}
	fmt.Fprintf(os.Stdout, "%s\n", string(r))
	return 0
}

// Example new command definition (can be anywhere in this package)
// 	func Command1(c *client.Client, args ...string) ([]byte, error) {
// 		// argument validation here
// 		rsp, err := c.Http.Get(c.Cfg.Root + "/endopint")
// 		if err != nil {
// 			return nil, fmt.Errorf("failed doing something!")
// 		}
// 		if rsp.StatusCode == http.StatusUnauthorized {
// 			return nil, client.ErrUnauthorized
// 		}
//		// other logic here
// 		return "result string (probably response json body)", nil
// 	}
