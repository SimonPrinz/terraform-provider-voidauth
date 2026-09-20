package client

type AdminClaim struct {
	Claim *string `json:"claim,omitempty"`
	Value *string `json:"value,omitempty"`
}
