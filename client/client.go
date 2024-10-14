package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hscli/config"
	"net/http"
	"time"

	"go.nhat.io/cookiejar"
	"golang.org/x/net/publicsuffix"
)

var ProgramName = "hs-cli"
var ProgramVersion = "0.0.1"

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
		return fmt.Errorf("failed encoding payload: %w", err)
	}
	req, err := http.NewRequest("POST", c.Cfg.Root+"/login", bytes.NewReader(payloadJson))
	if err != nil {
		return fmt.Errorf("failed creating POST request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	rsp, err := c.Http.Do(req)
	if err != nil {
		return fmt.Errorf("failed making POST request: %w", err)
	}
	rsp.Body.Close() // nothing relevant here
	if rsp.StatusCode != 200 {
		return fmt.Errorf("response with non 200 code %d", rsp.StatusCode)
	}
	return nil
}
