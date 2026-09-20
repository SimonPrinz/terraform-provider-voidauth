package client

type PublicService struct {
	client *Client
}

func (s *PublicService) Config() (*PublicConfigResponse, error) {
	var v *PublicConfigResponse
	_, err := s.client.doRequest(
		"GET", "/api/public/config",
		nil,
		&v,
	)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (s *PublicService) PasswordStrength(password string) (*PublicPasswordStrengthResponse, error) {
	var v *PublicPasswordStrengthResponse
	_, err := s.client.doRequest(
		"POST", "/api/public/passwordStrength",
		&PublicPasswordStrengthRequest{
			Password: password,
		},
		&v,
	)
	if err != nil {
		return nil, err
	}

	return v, nil
}
