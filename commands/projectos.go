package commands

import (
	"fmt"
	"hscli/client"
	"io"
	"net/http"
)

func GetProjects(c *client.Client, args ...string) ([]byte, error) {
	var endpoint string = c.Cfg.Root + "/projects"
	rsp, err := c.Http.Get(endpoint)
	if err != nil {
		return nil, NewCommandError("Failed requesting server", fmt.Errorf("http.Get %s: %w", endpoint, err))
	}
	defer rsp.Body.Close()
	if rsp.StatusCode == http.StatusUnauthorized {
		return nil, client.ErrUnauthorized
	}
	if rsp.StatusCode != http.StatusOK {
		return nil, NewCommandError(fmt.Sprintf("%d %s", rsp.StatusCode, http.StatusText(rsp.StatusCode)), nil)
	}
	rspData, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, NewCommandError("Failed receiving server response", fmt.Errorf("io.ReadAll: %w", err))
	}
	return rspData, nil
}

func GetProjectByID(c *client.Client, args ...string) ([]byte, error) {
	if len(args) == 0 {
		return nil, NewCommandError("Missing argument to command, expect 1 got 0", nil)
	}

	var endpoint string = c.Cfg.Root + "/projects/" + args[0]
	rsp, err := c.Http.Get(endpoint)
	if err != nil {
		return nil, NewCommandError("Failed requesting server", fmt.Errorf("http.Get %s: %w", endpoint, err))
	}
	defer rsp.Body.Close()
	if rsp.StatusCode == http.StatusUnauthorized {
		return nil, client.ErrUnauthorized
	}
	if rsp.StatusCode == http.StatusNotFound {
		return nil, NewCommandError("Member with provided ID does not exist", nil)
	}
	if rsp.StatusCode != http.StatusOK {
		return nil, NewCommandError(fmt.Sprintf("%d %s", rsp.StatusCode, http.StatusText(rsp.StatusCode)), nil)
	}
	rspData, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, NewCommandError("Failed receiving server response", fmt.Errorf("io.ReadAll: %w", err))
	}
	return rspData, nil
}
