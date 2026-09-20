package client

import "encoding/json"

type AdminGroupResponse struct {
	Id          *string `json:"id"`
	Name        *string `json:"name"`
	MfaRequired *bool   `json:"mfaRequired"`
	AutoAssign  *bool   `json:"autoAssign"`
	Users       *[]struct {
		Id       *string `json:"id"`
		Username *string `json:"username"`
	} `json:"users"`
	CustomClaims *[]AdminClaim `json:"customClaims"`
}

func (s *AdminService) GetGroup(id string) (*AdminGroupResponse, error) {
	var v *AdminGroupResponse
	_, err := s.client.doRequest(
		"GET", "/api/admin/group/"+id,
		nil,
		&v,
	)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (r *AdminGroupResponse) UnmarshalJSON(data []byte) error {
	type alias AdminGroupResponse
	var shadow struct {
		*alias
		MfaRequired json.RawMessage `json:"mfaRequired"`
		AutoAssign  json.RawMessage `json:"autoAssign"`
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
	setField(shadow.AutoAssign, &r.AutoAssign)
	return nil
}

func (s *AdminService) CreateGroup(request *AdminGroupCreateRequest) (*AdminGroupResponse, error) {
	var v *AdminGroupResponse
	_, err := s.client.doRequest(
		"POST", "/api/admin/group",
		&request,
		&v,
	)
	if err != nil {
		return nil, err
	}

	return v, nil
}

type AdminGroupCreateRequest struct {
	Id           string         `json:"id,omitempty"`
	Name         string         `json:"name"`
	MfaRequired  bool           `json:"mfaRequired"`
	AutoAssign   bool           `json:"autoAssign"`
	Users        []AdminUserRef `json:"users"`
	CustomClaims []AdminClaim   `json:"customClaims"`
}

type AdminUserRef struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

func (s *AdminService) DeleteGroup(id string) (bool, error) {
	_, err := s.client.doRequest(
		"DELETE", "/api/admin/group/"+id,
		nil,
		nil,
	)
	if err != nil {
		return false, err
	}

	return true, nil
}
