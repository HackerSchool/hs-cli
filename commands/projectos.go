package commands

import (
	"fmt"
	"hscli/client"
	"io"
	"net/http"
	"os"
)

func GetProjects(c *client.Client, args ...string) ([]byte, error) {
	var endpoint string = c.Cfg.Root + "/projects"
	rsp, err := c.Http.Get(endpoint)
	if err != nil {
		return nil, NewCommandError("Failed requesting server", fmt.Errorf("http.Get %s: %w", endpoint, err))
	}
	defer rsp.Body.Close()

	rspData, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, NewCommandError("Failed receiving server response", fmt.Errorf("io.ReadAll: %w", err))
	}
	if rsp.StatusCode == http.StatusUnauthorized {
		return rspData, NewCommandError(string(rspData), client.ErrUnauthorized)
	}
	if rsp.StatusCode != http.StatusOK {
		return nil, NewCommandError(string(rspData), nil)
	}
	return rspData, nil
}

func GetProjectByID(c *client.Client, args ...string) ([]byte, error) {
	if len(args) == 0 {
		return nil, NewCommandError("Missing argument to command, expecteded 1 got 0", nil)
	}

	var endpoint string = c.Cfg.Root + "/projects/" + args[0]
	rsp, err := c.Http.Get(endpoint)
	if err != nil {
		return nil, NewCommandError("Failed requesting server", fmt.Errorf("http.Get %s: %w", endpoint, err))
	}
	defer rsp.Body.Close()

	rspData, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, NewCommandError("Failed receiving server response", fmt.Errorf("io.ReadAll: %w", err))
	}
	if rsp.StatusCode == http.StatusUnauthorized {
		return rspData, NewCommandError(string(rspData), client.ErrUnauthorized)
	}
	if rsp.StatusCode != http.StatusOK {
		return nil, NewCommandError(string(rspData), nil)
	}
	return rspData, nil
}

func GetProjectMembers(c *client.Client, args ...string) ([]byte, error) {
	if len(args) == 0 {
		return nil, NewCommandError("Missing argument to command, expecteded 1 got 0", nil)
	}

	var endpoint string = c.Cfg.Root + "/projects/" + args[0] + "/members"
	rsp, err := c.Http.Get(endpoint)
	if err != nil {
		return nil, NewCommandError("Failed requesting server", fmt.Errorf("http.Get %s: %w", endpoint, err))
	}
	defer rsp.Body.Close()

	rspData, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, NewCommandError("Failed receiving server response", fmt.Errorf("io.ReadAll: %w", err))
	}
	if rsp.StatusCode == http.StatusUnauthorized {
		return rspData, NewCommandError(string(rspData), client.ErrUnauthorized)
	}
	if rsp.StatusCode != http.StatusOK {
		return nil, NewCommandError(string(rspData), nil)
	}
	return rspData, nil
}

func CreateProject(c *client.Client, args ...string) ([]byte, error) {
	if len(args) == 0 {
		return nil, NewCommandError("Missing argument to command, expecteded 1 got 0", nil)
	}
	var filePath string = args[0]
	f, err := os.Open(filePath)
	if err != nil {
		return nil, NewCommandError("Failed opening file", fmt.Errorf("os.Open %s: %w", filePath, err))
	}
	defer f.Close()

	var endpoint string = c.Cfg.Root + "/projects"
	rsp, err := c.Http.Post(endpoint, "application/json", f)
	if err != nil {
		return nil, NewCommandError("Failed requesting server", fmt.Errorf("http.Post %s: %w", endpoint, err))
	}
	defer rsp.Body.Close()

	rspData, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, NewCommandError("Failed receiving server response", fmt.Errorf("io.ReadAll: %w", err))
	}
	if rsp.StatusCode == http.StatusUnauthorized {
		return rspData, NewCommandError(string(rspData), client.ErrUnauthorized)
	}
	if rsp.StatusCode != http.StatusOK && rsp.StatusCode != http.StatusCreated {
		return nil, NewCommandError(string(rspData), nil)
	}
	return rspData, nil
}

func UpdateProject(c *client.Client, args ...string) ([]byte, error) {
	if len(args) != 2 {
		return nil, NewCommandError(fmt.Sprintf("Missing argument to command, expecteded 2 got %d", len(args)), nil)
	}

	var filePath string = args[1]
	f, err := os.Open(filePath)
	if err != nil {
		return nil, NewCommandError("Failed opening file", fmt.Errorf("os.Open %s: %w", filePath, err))
	}
	defer f.Close()

	var endpoint string = c.Cfg.Root + "/projects/" + args[0]
	req, err := http.NewRequest("PUT", endpoint, f)
	if err != nil {
		return nil, NewCommandError("Failed creating request server", fmt.Errorf("http.NewRequest PUT %s: %w", endpoint, err))
	}
	req.Header.Add("Content-Type", "application/json")
	rsp, err := c.Http.Do(req)
	if err != nil {
		return nil, NewCommandError("Failed requesting server", fmt.Errorf("http.Client.Do %s: %w", endpoint, err))
	}
	defer rsp.Body.Close()

	rspData, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, NewCommandError("Failed receiving server response", fmt.Errorf("io.ReadAll: %w", err))
	}
	if rsp.StatusCode == http.StatusUnauthorized {
		return rspData, NewCommandError(string(rspData), client.ErrUnauthorized)
	}
	if rsp.StatusCode != http.StatusOK && rsp.StatusCode != http.StatusCreated {
		return nil, NewCommandError(string(rspData), nil)
	}
	return rspData, nil
}

func DeleteProject(c *client.Client, args ...string) ([]byte, error) {
	if len(args) == 0 {
		return nil, NewCommandError("Missing argument to command, expecteded 1 got 0", nil)
	}

	var endpoint string = c.Cfg.Root + "/projects/" + args[0]
	req, err := http.NewRequest("DELETE", endpoint, nil)
	if err != nil {
		return nil, NewCommandError("Failed creating request server", fmt.Errorf("http.NewRequest DELETE %s: %w", endpoint, err))
	}
	rsp, err := c.Http.Do(req)
	if err != nil {
		return nil, NewCommandError("Failed requesting server", fmt.Errorf("http.Client.Do %s: %w", endpoint, err))
	}
	defer rsp.Body.Close()

	rspData, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, NewCommandError("Failed receiving server response", fmt.Errorf("io.ReadAll: %w", err))
	}
	if rsp.StatusCode == http.StatusUnauthorized {
		return rspData, NewCommandError(string(rspData), client.ErrUnauthorized)
	}
	if rsp.StatusCode != http.StatusOK {
		return nil, NewCommandError(string(rspData), nil)
	}
	return rspData, nil
}

func GetProjectLogo(c *client.Client, args ...string) ([]byte, error) {
	if len(args) == 0 {
		return nil, NewCommandError("Missing argument to command, expecteded 1 got 0", nil)
	}

	var endpoint string = c.Cfg.Root + "/projects/" + args[0] + "/logo"
	rsp, err := c.Http.Get(endpoint)
	if err != nil {
		return nil, NewCommandError("Failed requesting server", fmt.Errorf("http.Get %s: %w", endpoint, err))
	}
	defer rsp.Body.Close()

	rspData, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, NewCommandError("Failed receiving server response", fmt.Errorf("io.ReadAll: %w", err))
	}
	if rsp.StatusCode == http.StatusUnauthorized {
		return rspData, NewCommandError(string(rspData), client.ErrUnauthorized)
	}
	if rsp.StatusCode != http.StatusOK {
		return nil, NewCommandError(string(rspData), nil)
	}
	return rspData, nil
}
