package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hscli/client"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
)

func Login(c *client.Client, args ...string) ([]byte, error) {
	var payload map[string]string = map[string]string{
		"username": c.Cfg.User,
		"password": c.Cfg.Password,
	}
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %w", err)
	}

	var endpoint string = c.Cfg.Root + "/login"
	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(payloadJson))
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest POST %s: %w", endpoint, err)
	}
	req.Header.Set("Content-Type", "application/json")

	rsp, err := c.Http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http.Do: %w", err)
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

func GetMembers(c *client.Client, args ...string) ([]byte, error) {
	var endpoint string = c.Cfg.Root + "/members"
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

func GetMemberByUsername(c *client.Client, args ...string) ([]byte, error) {
	if len(args) == 0 {
		return nil, NewCommandError("Missing argument to command, expected 1 got 0", nil)
	}

	var endpoint string = c.Cfg.Root + "/members/" + args[0]
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

func GetMemberProjects(c *client.Client, args ...string) ([]byte, error) {
	if len(args) == 0 {
		return nil, NewCommandError("Missing argument to command, expected 1 got 0", nil)
	}

	var endpoint string = c.Cfg.Root + "/members/" + args[0] + "/projects"
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

func GetMemberLogo(c *client.Client, args ...string) ([]byte, error) {
	if len(args) == 0 {
		return nil, NewCommandError("Missing argument to command, expected 1 got 0", nil)
	}

	var endpoint string = c.Cfg.Root + "/members/" + args[0] + "/logo"
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

func CreateMember(c *client.Client, args ...string) ([]byte, error) {
	if len(args) == 0 {
		return nil, NewCommandError("Missing argument to command, expected 1 got 0", nil)
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

func UpdateMember(c *client.Client, args ...string) ([]byte, error) {
	if len(args) != 2 {
		return nil, NewCommandError(fmt.Sprintf("Missing argument to command, expected 2 got %d", len(args)), nil)
	}

	var filePath string = args[1]
	f, err := os.Open(filePath)
	if err != nil {
		return nil, NewCommandError("Failed opening file", fmt.Errorf("os.Open %s: %w", filePath, err))
	}
	defer f.Close()

	var endpoint string = c.Cfg.Root + "/members/" + args[0]
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

func UpdateMemberLogo(c *client.Client, args ...string) ([]byte, error) {
	if len(args) != 2 {
		return nil, NewCommandError(fmt.Sprintf("Missing arugments to commands, expected 2 got %d", len(args)), nil)
	}

	var filePath string = args[1]
	f, err := os.Open(filePath)
	if err != nil {
		return nil, NewCommandError("Failed opening file", fmt.Errorf("os.Open %s: %w", filePath, err))
	}
	defer f.Close()

	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	// Detect MIME type by file extension
	ext := filepath.Ext(filePath)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream" // default if not detected
	}
	// Create a form file part with custom content type
	partHeaders := make(textproto.MIMEHeader)
	partHeaders.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "file", filepath.Base(filePath)))
	partHeaders.Set("Content-Type", mimeType)
	part, err := w.CreatePart(partHeaders)
	if err != nil {
		return nil, NewCommandError("Failed creating multipart form", fmt.Errorf("multipart.Writer.CreatePart: %w", err))
	}
	_, err = io.Copy(part, f)
	if err != nil {
		return nil, NewCommandError("Failed writing multipart form", fmt.Errorf("io.Copy: %w", err))
	}
	w.Close()
	fmt.Println()

	var endpoint string = c.Cfg.Root + "/members/" + args[0] + "/logo"
	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(body.Bytes()))
	if err != nil {
		return nil, NewCommandError("Failed creating server request", fmt.Errorf("http.NewRequest POST %s: %w", endpoint, err))
	}
	fmt.Println(w.FormDataContentType())
	req.Header.Add("Content-Type", w.FormDataContentType())
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

func DeleteMember(c *client.Client, args ...string) ([]byte, error) {
	if len(args) == 0 {
		return nil, NewCommandError("Missing argument to command, expect 1 got 0", nil)
	}

	var endpoint string = c.Cfg.Root + "/members/" + args[0]
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

func AddProject(c *client.Client, args ...string) ([]byte, error) {
	if len(args) != 3 {
		return nil, NewCommandError(fmt.Sprintf("Missing argument to command, expected 3 got %d", len(args)), nil)
	}

	var filePath string = args[2]
	f, err := os.Open(filePath)
	if err != nil {
		return nil, NewCommandError("Failed opening file", fmt.Errorf("os.Open %s: %w", filePath, err))
	}
	defer f.Close()

	var endpoint string = c.Cfg.Root + "/members/" + args[0] + "/" + args[1]
	req, err := http.NewRequest("POST", endpoint, f)
	if err != nil {
		return nil, NewCommandError("Failed creating request server", fmt.Errorf("http.NewRequest POST %s: %w", endpoint, err))
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

func GetTags(c *client.Client, args ...string) ([]byte, error) {
	if len(args) == 0 {
		return nil, NewCommandError("Missing argument to command, expected 1 got 0", nil)
	}

	var endpoint string = c.Cfg.Root + "/members/" + args[0] + "/tags"
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

func AddTag(c *client.Client, args ...string) ([]byte, error) {
	if len(args) != 2 {
		return nil, NewCommandError(fmt.Sprintf("Missing argument to command, expected 2 got %d", len(args)), nil)
	}

	var filePath string = args[1]
	f, err := os.Open(filePath)
	if err != nil {
		return nil, NewCommandError("Failed opening file", fmt.Errorf("os.Open %s: %w", filePath, err))
	}
	defer f.Close()

	var endpoint string = c.Cfg.Root + "/members/" + args[0] + "/tags"
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

func DeleteTag(c *client.Client, args ...string) ([]byte, error) {
	if len(args) != 2 {
		return nil, NewCommandError(fmt.Sprintf("Missing argument to command, expected 2 got %d", len(args)), nil)
	}

	var filePath string = args[1]
	f, err := os.Open(filePath)
	if err != nil {
		return nil, NewCommandError("Failed opening file", fmt.Errorf("os.Open %s: %w", filePath, err))
	}
	defer f.Close()

	var endpoint string = c.Cfg.Root + "/members/" + args[2] + "/tags"
	req, err := http.NewRequest("PUT", endpoint, f)
	if err != nil {
		return nil, NewCommandError("Failed creating request server", fmt.Errorf("http.NewRequest DELETE %s: %w", endpoint, err))
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
