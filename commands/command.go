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
		r, err := cmd(c, args...)
		if errors.Is(err, client.ErrUnauthorized) {
			logging.LogDebug("Retrying command with login")
			if err := c.Login(); err != nil {
				if errors.Is(err, client.ErrUnauthorized) {
					return nil, NewCommandError("Unauthorized!", nil)
				}
				return nil, NewCommandError("Failed to log in", err)
			}
			return cmd(c, args...)
		}
		return r, err
	}
}

// Runs a command.
// Returns 0 on success, 1 on an domain related errors such as (unauthorized, resource doesn't exist, etc) and 2 on generic errors (no connection, etc)
func RunCommand(c *client.Client, cmd Command, args ...string) int {
	r, err := cmd(c, args...)
	if err != nil {
		var retCode int = 2

		var commandErr CommandError
		if errors.As(err, &commandErr) {
			if commandErr.Cause == nil { // business logic error
				retCode = 1
			} else {
				logging.LogDebug(commandErr.Cause.Error())
			}
		}

		fmt.Fprintf(os.Stderr, "%s\n", err)
		return retCode
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
