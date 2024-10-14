package commands

import (
	"fmt"
	"hscli/client"
	"io"
	"net/http"
	"os"
)

func GetMembers(c *client.Client, args ...string) ([]byte, error) {
	rsp, err := c.Http.Get(c.Cfg.Root + "/members")
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

func GetMemberByUsername(c *client.Client, args ...string) ([]byte, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("missing arguments to command call, expected at least 1 got 0")
	}
	rsp, err := c.Http.Get(c.Cfg.Root + "/members/" + args[0])
	if err != nil {
		return nil, fmt.Errorf("failed making request: %w", err)
	}
	defer rsp.Body.Close()
	if rsp.StatusCode == http.StatusUnauthorized {
		return nil, ErrUnauthorized
	}
	if rsp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("Member with provided username does not exist!")
	}
	rspData, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response: %w", err)
	}
	return rspData, nil
}

func CreateMember(c *client.Client, args ...string) ([]byte, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("missing arguments to command call, expected at least 1 got 0")
	}
	f, err := os.Open(args[0])
	if err != nil {
		return nil, fmt.Errorf("failed opening file: %w", err)
	}
	defer f.Close()
	rsp, err := c.Http.Post(c.Cfg.Root+"/members", "application/json", f)
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
