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

	rspData, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, NewCommandError("Failed receiving server response", fmt.Errorf("io.ReadAll: %w", err))
	}
	if rsp.StatusCode != http.StatusOK {
		return nil, NewCommandError(fmt.Sprintf("%d %s\n%s", rsp.StatusCode, http.StatusText(rsp.StatusCode), string(rspData)), nil)
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

	rspData, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, NewCommandError("Failed receiving server response", fmt.Errorf("io.ReadAll: %w", err))
	}
	if rsp.StatusCode != http.StatusOK {
		return nil, NewCommandError(fmt.Sprintf("%d %s\n%s", rsp.StatusCode, http.StatusText(rsp.StatusCode), string(rspData)), nil)
	}
	return rspData, nil
}

func GetProjectMembers(c *client.Client, args ...string) ([]byte, error) {
	if len(args) == 0 {
		return nil, NewCommandError("Missing argument to command, expect 1 got 0", nil)
	}

	var endpoint string = c.Cfg.Root + "/projects/" + args[0] + "/members"
	rsp, err := c.Http.Get(endpoint)
	if err != nil {
		return nil, NewCommandError("Failed requesting server", fmt.Errorf("http.Get %s: %w", endpoint, err))
	}
	defer rsp.Body.Close()
	if rsp.StatusCode == http.StatusUnauthorized {
		return nil, client.ErrUnauthorized
	}

	rspData, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, NewCommandError("Failed receiving server response", fmt.Errorf("io.ReadAll: %w", err))
	}
	if rsp.StatusCode != http.StatusOK {
		return nil, NewCommandError(fmt.Sprintf("%d %s\n%s", rsp.StatusCode, http.StatusText(rsp.StatusCode), string(rspData)), nil)
	}
	return rspData, nil
}

func CreateProject(c *client.Client, args ...string) ([]byte, error) {
	// TODO
	if len(args) == 0 {
		return nil, NewCommandError("Missing argument to command, expect 1 got 0", nil)
	}
	return []byte{}, nil
}

func UpdateProject(c *client.Client, args ...string) ([]byte, error) {
	// TODO
	if len(args) == 0 {
		return nil, NewCommandError(fmt.Sprintf("Missing argument to command, expect 2 got %d", len(args)), nil)
	}
	return []byte{}, nil
}

func DeleteProject(c *client.Client, args ...string) ([]byte, error) {
	// TODO
	if len(args) == 0 {
		return nil, NewCommandError(fmt.Sprintf("Missing argument to command, expect 2 got %d", len(args)), nil)
	}
	return []byte{}, nil
}
