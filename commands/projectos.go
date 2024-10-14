package commands

import (
	"fmt"
	"hscli/client"
	"io"
	"net/http"
)

func GetProjects(c *client.Client, args ...string) ([]byte, error) {
	rsp, err := c.Http.Get(c.Cfg.Root + "/projects")
	if err != nil {
		return nil, fmt.Errorf("failed making request: %w", err)
	}
	defer rsp.Body.Close()
	if rsp.StatusCode == http.StatusUnauthorized {
		return nil, ErrUnauthorized
	}
	rspData, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response: %w", err)
	}
	return rspData, nil
}

func GetProjectByID(c *client.Client, args ...string) ([]byte, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("missing arguments to command call, expected at least 1 got 0")
	}
	fmt.Println(c.Cfg.Root + "/projects/" + args[0])
	rsp, err := c.Http.Get(c.Cfg.Root + "/projects/" + args[0])
	if err != nil {
		return nil, fmt.Errorf("failed making request: %w", err)
	}
	defer rsp.Body.Close()
	if rsp.StatusCode == http.StatusUnauthorized {
		return nil, ErrUnauthorized
	}
	if rsp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("Project with provided ID does not exist!")
	}
	rspData, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response: %w", err)
	}
	return rspData, nil
}
