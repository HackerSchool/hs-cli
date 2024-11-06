package client

import (
	"errors"
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
