package client

import "encoding/json"

func (s *AdminService) GetProxyAuth(id string) (*AdminProxyAuthResponse, error) {
	var v *AdminProxyAuthResponse
	_, err := s.client.doRequest(
		"GET", "/api/admin/proxyauth/"+id,
		nil,
		&v,
	)
	if err != nil {
		return nil, err
	}

	return v, nil
}

type AdminProxyAuthResponse struct {
	Id               *string   `json:"id"`
	Domain           *string   `json:"domain"`
	MfaRequired      *bool     `json:"mfaRequired"`
	MaxSessionLength *int      `json:"maxSessionLength"`
	Groups           *[]string `json:"groups"`
}

func (r *AdminProxyAuthResponse) UnmarshalJSON(data []byte) error {
	type alias AdminProxyAuthResponse
	var shadow struct {
		*alias
		MfaRequired json.RawMessage `json:"mfaRequired"`
	}
	shadow.alias = (*alias)(r)
	if err := json.Unmarshal(data, &shadow); err != nil {
		return err
	}
	setField := func(raw json.RawMessage, dst **bool) {
		if len(raw) == 0 || string(raw) == "null" {
			return
		}
		if v, ok := rawBool(raw); ok {
			*dst = &v
		}
	}
	setField(shadow.MfaRequired, &r.MfaRequired)
	return nil
}

func (s *AdminService) CreateProxyAuth(request *AdminProxyAuthCreateRequest) (*AdminProxyAuthResponse, error) {
	var v *AdminProxyAuthResponse
	_, err := s.client.doRequest(
		"POST", "/api/admin/proxyauth",
		&request,
		&v,
	)
	if err != nil {
		return nil, err
	}

	return v, nil
}

type AdminProxyAuthCreateRequest struct {
	Id               string   `json:"id,omitempty"`
	Domain           string   `json:"domain"`
	MfaRequired      bool     `json:"mfaRequired"`
	MaxSessionLength int      `json:"maxSessionLength"`
	Groups           []string `json:"groups"`
}

func (s *AdminService) DeleteProxyAuth(id string) (bool, error) {
	_, err := s.client.doRequest(
		"DELETE", "/api/admin/proxyauth/"+id,
		nil,
		nil,
	)
	if err != nil {
		return false, err
	}

	return true, nil
}
