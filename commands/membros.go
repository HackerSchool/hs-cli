package commands

import (
	"fmt"
	"hscli/client"
	"io"
	"net/http"
	"os"
)

func GetMembers(c *client.Client, args ...string) ([]byte, error) {
	var endpoint string = c.Cfg.Root + "/members"
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

func GetMemberByUsername(c *client.Client, args ...string) ([]byte, error) {
	if len(args) == 0 {
		return nil, NewCommandError("Missing argument to command, expect 1 got 0", nil)
	}

	var endpoint string = c.Cfg.Root + "/members/" + args[0]
	rsp, err := c.Http.Get(endpoint)
	if err != nil {
		return nil, NewCommandError("Failed requesting server", fmt.Errorf("http.Get %s: %w", endpoint, err))
	}
	defer rsp.Body.Close()
	if rsp.StatusCode == http.StatusUnauthorized {
		return nil, client.ErrUnauthorized
	}
	if rsp.StatusCode == http.StatusNotFound {
		return nil, NewCommandError("Member with provided username does not exist", nil)
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

func CreateMember(c *client.Client, args ...string) ([]byte, error) {
	if len(args) == 0 {
		return nil, NewCommandError("Missing argument to command, expect 1 got 0", nil)
	}

	var filePath string = args[0]
	f, err := os.Open(filePath)
	if err != nil {
		return nil, NewCommandError("Failed opening file", fmt.Errorf("os.Open %s: %w", filePath, err))
	}
	defer f.Close()

	var endpoint string = c.Cfg.Root + "/members"
	rsp, err := c.Http.Post(endpoint, "application/json", f)
	if err != nil {
		return nil, NewCommandError("Failed requesting server", fmt.Errorf("http.Post %s: %w", endpoint, err))
	}
	defer rsp.Body.Close()
	if rsp.StatusCode == http.StatusUnauthorized {
		return nil, client.ErrUnauthorized
	}
	if rsp.StatusCode != http.StatusOK && rsp.StatusCode != http.StatusCreated {
		return nil, NewCommandError(fmt.Sprintf("%d %s", rsp.StatusCode, http.StatusText(rsp.StatusCode)), nil)
	}
	rspData, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, NewCommandError("Failed receiving server response", fmt.Errorf("io.ReadAll: %w", err))
	}
	return rspData, nil
}

func UpdateMember(c *client.Client, args ...string) ([]byte, error) {
	if len(args) != 2 {
		return nil, NewCommandError("Missing argument to command, expect 1 got 0", nil)
	}

	var filePath string = args[1]
	f, err := os.Open(filePath)
	if err != nil {
		return nil, NewCommandError("Failed opening file", fmt.Errorf("os.Open %s: %w", filePath, err))
	}
	defer f.Close()

	var endpoint string = c.Cfg.Root + "/members/" + args[0]
	rsp, err := c.Http.Post(endpoint, "application/json", f)
	if err != nil {
		return nil, NewCommandError("Failed requesting server", fmt.Errorf("http.Post %s: %w", endpoint, err))
	}
	defer rsp.Body.Close()
	if rsp.StatusCode == http.StatusUnauthorized {
		return nil, client.ErrUnauthorized
	}
	if rsp.StatusCode != http.StatusOK && rsp.StatusCode != http.StatusCreated {
		return nil, NewCommandError(fmt.Sprintf("%d %s", rsp.StatusCode, http.StatusText(rsp.StatusCode)), nil)
	}
	rspData, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, NewCommandError("Failed receiving server response", fmt.Errorf("io.ReadAll: %w", err))
	}
	return rspData, nil
}
