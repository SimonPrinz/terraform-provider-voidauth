package provider

import (
	"net/http"

	"github.com/SimonPrinz/terraform-provider-voidauth/internal/client"
)

type VoidauthClient struct {
	client   *client.Client
	baseUrl  string
	username string
	password string
}

func NewVoidauthClient(baseUrl, username, password string) (*VoidauthClient, error) {
	httpClient := &http.Client{}

	v := &VoidauthClient{
		baseUrl:  baseUrl,
		username: username,
		password: password,
	}

	c, err := client.NewClient(baseUrl, httpClient, username, password)
	if err != nil {
		return nil, err
	}

	v.client = c
	return v, nil
}
