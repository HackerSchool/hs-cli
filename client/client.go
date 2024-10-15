package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"hscli/config"
	"net/http"
	"time"

	"go.nhat.io/cookiejar"
	"golang.org/x/net/publicsuffix"
)

const (
	ProgramName    = "hs-cli"
	ProgramVersion = "0.0.1"
)

var ErrUnauthorized = errors.New("Unauthorized!")

type Client struct {
	Http *http.Client
	Cfg  *config.Config
}

func NewClient() *Client {
	return &Client{
		Http: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse // don't follow redirects
			},
			Transport: WithUARoundTripper{
				r: WithLoggingRoundTripper{
					r: http.DefaultTransport,
				},
			},
			Timeout: 100 * time.Second, // default, TODO should be passed as CLI arg
		},
		Cfg: &config.Config{},
	}
}

func (c *Client) SetupJar() {
	c.Http.Jar = cookiejar.NewPersistentJar(
		cookiejar.WithFilePath(c.Cfg.CookieJarPath),
		cookiejar.WithAutoSync(true),
		cookiejar.WithPublicSuffixList(publicsuffix.List),
	)
}

// Logs in to the API saving the session cookie in the Jar
func (c *Client) Login() error {
	payload := map[string]string{
		"username": c.Cfg.User,
		"password": c.Cfg.Password,
	}
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}
	var endpoint string = c.Cfg.Root + "/login"
	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(payloadJson))
	if err != nil {
		return fmt.Errorf("http.NewRequest POST %s: %w", endpoint, err)
	}
	req.Header.Set("Content-Type", "application/json")

	rsp, err := c.Http.Do(req)
	if err != nil {
		return fmt.Errorf("http.Do: %w", err)
	}
	rsp.Body.Close() // nothing relevant here
	if rsp.StatusCode == http.StatusUnauthorized {
		return ErrUnauthorized
	}
	if rsp.StatusCode != http.StatusOK {
		return fmt.Errorf("%d %s", rsp.StatusCode, http.StatusText(rsp.StatusCode))
	}
	return nil
}
