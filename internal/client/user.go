package client

type UserService struct {
	client *Client
}

func (s *UserService) Me() (*UserMeResponse, error) {
	var v *UserMeResponse
	_, err := s.client.doRequest(
		"GET", "/api/user/me",
		nil,
		&v,
	)
	if err != nil {
		return nil, err
	}

	return v, nil
}
