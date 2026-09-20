package client

import "encoding/json"

type AdminUserResponse struct {
	Id            *string `json:"id"`
	Username      *string `json:"username"`
	Email         *string `json:"email"`
	Name          *string `json:"name"`
	EmailVerified *bool   `json:"emailVerified"`
	Approved      *bool   `json:"approved"`
	MfaRequired   *bool   `json:"mfaRequired"`
	// ToDo: expiresAt
	HasPassword *bool `json:"hasPassword"`
	HasEmail    *bool `json:"hasEmail"`
	IsAdmin     *bool `json:"isAdmin"`
	Groups      *[]struct {
		Id           *string       `json:"id"`
		Name         *string       `json:"name"`
		CustomClaims *[]AdminClaim `json:"customClaims"`
	} `json:"groups"`
	CustomClaims *[]AdminClaim `json:"customClaims"`
	HasTotp      *bool         `json:"hasTotp"`
	HasPasskeys  *bool         `json:"hasPasskeys"`
	HasMfaGroup  *bool         `json:"hasMfaGroup"`
}

func (s *AdminService) GetUser(id string) (*AdminUserResponse, error) {
	var v *AdminUserResponse
	_, err := s.client.doRequest(
		"GET", "/api/admin/user/"+id,
		nil,
		&v,
	)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (r *AdminUserResponse) UnmarshalJSON(data []byte) error {
	type alias AdminUserResponse
	var shadow struct {
		*alias
		EmailVerified json.RawMessage `json:"emailVerified"`
		Approved      json.RawMessage `json:"approved"`
		MfaRequired   json.RawMessage `json:"mfaRequired"`
		HasPassword   json.RawMessage `json:"hasPassword"`
		HasEmail      json.RawMessage `json:"hasEmail"`
		IsAdmin       json.RawMessage `json:"isAdmin"`
		HasTotp       json.RawMessage `json:"hasTotp"`
		HasPasskeys   json.RawMessage `json:"hasPasskeys"`
		HasMfaGroup   json.RawMessage `json:"hasMfaGroup"`
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
	setField(shadow.EmailVerified, &r.EmailVerified)
	setField(shadow.Approved, &r.Approved)
	setField(shadow.MfaRequired, &r.MfaRequired)
	setField(shadow.HasPassword, &r.HasPassword)
	setField(shadow.HasEmail, &r.HasEmail)
	setField(shadow.IsAdmin, &r.IsAdmin)
	setField(shadow.HasTotp, &r.HasTotp)
	setField(shadow.HasPasskeys, &r.HasPasskeys)
	setField(shadow.HasMfaGroup, &r.HasMfaGroup)
	return nil
}
