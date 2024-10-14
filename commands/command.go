package commands

import (
	"errors"
	"fmt"
	"hscli/client"
	"hscli/logging"
)

var ErrUnauthorized = errors.New("Unauthorized")

type Command func(c *client.Client, args ...string) ([]byte, error)

func WithLoginRetry(cmd Command) Command {
	return func(c *client.Client, args ...string) ([]byte, error) {
		r, err := cmd(c, args...)
		if errors.Is(err, ErrUnauthorized) {
			logging.LogDebug("Retrying command with login")
			if err := c.Login(); err != nil {
				return nil, fmt.Errorf("failed to login: %w", err)
			}
			return cmd(c, args...)
		}
		return r, err
	}
}

func RunCommand(c *client.Client, cmd Command, args ...string) ([]byte, error) {
	return cmd(c, args...)
}

// Example new command definition (can be anywhere in this package)
// 	func Command1(c *client.Client, args ...string) ([]byte, error) {
// 		// argument validation here
// 		rsp, err := c.Http.Get(c.Cfg.Root + "/endopint")
// 		if err != nil {
// 			return nil, fmt.Errorf("failed doing something!")
// 		}
// 		if rsp.StatusCode == http.StatusUnauthorized {
// 			return nil, ErrUnauthorized
// 		}
//		// other logic here
// 		return "result string (probably response json body)", nil
// 	}
