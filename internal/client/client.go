package client

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"sync"
)

type Client struct {
	BaseUrl    *url.URL
	httpClient *http.Client

	PublicService *PublicService
	UserService   *UserService
	AdminService  *AdminService

	authMu sync.Mutex
	authed bool

	username, password string
}

func NewClient(baseUrl string, httpClient *http.Client, username, password string) (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("creating cookie jar: %w", err)
	}

	if httpClient == nil {
		httpClient = &http.Client{Jar: jar}
	} else if httpClient.Jar == nil {
		httpClient.Jar = jar
	}

	parsedUrl, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}

	client := &Client{
		BaseUrl:    parsedUrl,
		httpClient: httpClient,
		username:   username,
		password:   password,
	}
	client.PublicService = &PublicService{client: client}
	client.UserService = &UserService{client: client}
	client.AdminService = &AdminService{client: client}

	return client, nil
}

func (c *Client) ensureAuthenticated() error {
	c.authMu.Lock()
	defer c.authMu.Unlock()

	if c.authed {
		return nil
	}
	if c.username == "" && c.password == "" {
		return nil
	}

	if err := c.login(); err != nil {
		return err
	}
	c.authed = true
	return nil
}

func (c *Client) resetAuth() {
	c.authMu.Lock()
	c.authed = false
	c.authMu.Unlock()
}

func (c *Client) isAuthenticated() bool {
	c.authMu.Lock()
	defer c.authMu.Unlock()
	return c.authed
}

func (c *Client) login() error {
	base := strings.TrimRight(c.BaseUrl.String(), "/")

	state := make([]byte, 16)
	if _, err := rand.Read(state); err != nil {
		return fmt.Errorf("voidauth login: generating state: %w", err)
	}

	authURL := base + "/oidc/auth?" + url.Values{
		"client_id":     {"auth_internal_client"},
		"response_type": {"none"},
		"redirect_uri":  {base},
		"scope":         {"openid"},
		"state":         {hex.EncodeToString(state)},
		"prompt":        {"login"},
	}.Encode()

	req1, err := http.NewRequest(http.MethodGet, authURL, nil)
	if err != nil {
		return fmt.Errorf("voidauth login: building auth request: %w", err)
	}
	resp1, err := c.httpClient.Do(req1)
	if err != nil {
		return fmt.Errorf("voidauth login: starting authorization: %w", err)
	}
	io.Copy(io.Discard, resp1.Body)
	resp1.Body.Close()
	if resp1.StatusCode >= 400 {
		return fmt.Errorf("voidauth login: authorization request failed: %s", resp1.Status)
	}

	creds, err := json.Marshal(map[string]string{
		"input":    c.username,
		"password": c.password,
	})
	if err != nil {
		return fmt.Errorf("voidauth login: encoding credentials: %w", err)
	}
	req2, err := http.NewRequest(http.MethodPost, base+"/api/interaction/login", bytes.NewReader(creds))
	if err != nil {
		return fmt.Errorf("voidauth login: building credential request: %w", err)
	}
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("Accept", "application/json")

	resp2, err := c.httpClient.Do(req2)
	if err != nil {
		return fmt.Errorf("voidauth login: submitting credentials: %w", err)
	}

	if resp2.StatusCode != http.StatusOK {
		resp2.Body.Close()
		return fmt.Errorf("voidauth login: credentials rejected (HTTP %d) — verify the configured username/password", resp2.StatusCode)
	}

	var loginResp struct {
		Success  *bool  `json:"success"`
		Location string `json:"location"`
	}
	decodeErr := json.NewDecoder(resp2.Body).Decode(&loginResp)
	resp2.Body.Close()
	if decodeErr != nil {
		return fmt.Errorf("voidauth login: decoding credential response: %w", decodeErr)
	}
	if loginResp.Success == nil || !*loginResp.Success {
		return fmt.Errorf("voidauth login: credentials rejected — verify the configured username/password")
	}
	if loginResp.Location == "" {
		return fmt.Errorf("voidauth login: server accepted credentials but returned no resume location")
	}

	resumeURL, err := url.Parse(loginResp.Location)
	if err != nil {
		return fmt.Errorf("voidauth login: invalid resume location %q: %w", loginResp.Location, err)
	}
	if !resumeURL.IsAbs() {
		resumeURL = c.BaseUrl.ResolveReference(resumeURL)
	}
	req3, err := http.NewRequest(http.MethodGet, resumeURL.String(), nil)
	if err != nil {
		return fmt.Errorf("voidauth login: building resume request: %w", err)
	}
	resp3, err := c.httpClient.Do(req3)
	if err != nil {
		return fmt.Errorf("voidauth login: completing authorization: %w", err)
	}
	io.Copy(io.Discard, resp3.Body)
	resp3.Body.Close()
	if resp3.StatusCode >= 400 {
		return fmt.Errorf("voidauth login: completing authorization failed: %s", resp3.Status)
	}

	return nil
}

func (c *Client) doRequest(method, path string, body interface{}, out interface{}) (*http.Response, error) {
	rel := &url.URL{Path: path}
	u := c.BaseUrl.ResolveReference(rel)

	protected := !strings.HasPrefix(path, "/api/public/")

	encodeBody := func() (*bytes.Buffer, error) {
		if body == nil {
			return nil, nil
		}
		buf := new(bytes.Buffer)
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return nil, err
		}
		return buf, nil
	}

	if protected {
		if err := c.ensureAuthenticated(); err != nil {
			return nil, err
		}
	}

	buf, err := encodeBody()
	if err != nil {
		return nil, err
	}

	resp, err := c.sendRequest(method, u, buf)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusUnauthorized && protected && c.isAuthenticated() {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()

		c.resetAuth()
		if err := c.ensureAuthenticated(); err != nil {
			return nil, err
		}

		buf, err = encodeBody()
		if err != nil {
			return nil, err
		}
		resp, err = c.sendRequest(method, u, buf)
		if err != nil {
			return nil, err
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("voidauth login: reading response body: %w", err)
		}
		bodyString := string(bodyBytes)
		return nil, fmt.Errorf("voidauth login: authentication failed: %s", bodyString)
	}

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		if err == io.EOF {
			return resp, nil
		}
		return nil, err
	}

	return resp, nil
}

func (c *Client) sendRequest(method string, u *url.URL, buf *bytes.Buffer) (*http.Response, error) {
	var bodyReader io.Reader
	if buf != nil {
		bodyReader = buf
	}
	req, err := http.NewRequest(method, u.String(), bodyReader)
	if err != nil {
		return nil, err
	}

	if buf != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	return c.httpClient.Do(req)
}
