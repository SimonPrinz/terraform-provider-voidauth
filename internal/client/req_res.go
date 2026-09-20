package client

import "encoding/json"

type PublicConfigResponse struct {
	Domain            *string `json:"domain,omitempty"`
	AppName           *string `json:"appName,omitempty"`
	ZxcvbnMin         *int32  `json:"zxcvbnMin,omitempty"`
	EmailActive       *bool   `json:"emailActive,omitempty"`
	EmailVerification *bool   `json:"emailVerification,omitempty"`
	Registration      *bool   `json:"registration,omitempty"`
	ContactEmail      *string `json:"contactEmail,omitempty"`    // nullable?
	DefaultRedirect   *string `json:"defaultRedirect,omitempty"` // nullable?
	MfaRequired       *bool   `json:"mfaRequired,omitempty"`
}

type PublicPasswordStrengthRequest struct {
	Password string `json:"password"`
}

type PublicPasswordStrengthResponse struct {
	CalcTime     int     `json:"calcTime"`
	Password     string  `json:"password"`
	Guesses      int     `json:"guesses"`
	GuessesLog10 float64 `json:"guessesLog10"`
	Sequence     []struct {
		Pattern             string  `json:"pattern"`
		I                   int     `json:"i"`
		J                   int     `json:"j"`
		Token               string  `json:"token"`
		MatchedWord         string  `json:"matchedWord"`
		Rank                int     `json:"rank"`
		DictionaryName      string  `json:"dictionaryName"`
		Reversed            bool    `json:"reversed"`
		L33T                bool    `json:"l33t"`
		BaseGuesses         int     `json:"baseGuesses"`
		UppercaseVariations int     `json:"uppercaseVariations"`
		L33TVariations      int     `json:"l33tVariations"`
		Guesses             int     `json:"guesses"`
		GuessesLog10        float64 `json:"guessesLog10"`
	} `json:"sequence"`
	CrackTimes struct {
		OnlineThrottlingXPerHour struct {
			Base    int    `json:"base"`
			Seconds int    `json:"seconds"`
			Display string `json:"display"`
		} `json:"onlineThrottlingXPerHour"`
		OnlineNoThrottlingXPerSecond struct {
			Base    int    `json:"base"`
			Seconds int    `json:"seconds"`
			Display string `json:"display"`
		} `json:"onlineNoThrottlingXPerSecond"`
		OfflineSlowHashingXPerSecond struct {
			Base    interface{} `json:"base"`
			Seconds float64     `json:"seconds"`
			Display string      `json:"display"`
		} `json:"offlineSlowHashingXPerSecond"`
		OfflineFastHashingXPerSecond struct {
			Base    interface{} `json:"base"`
			Seconds float64     `json:"seconds"`
			Display string      `json:"display"`
		} `json:"offlineFastHashingXPerSecond"`
	} `json:"crackTimes"`
	Score    int `json:"score"`
	Feedback struct {
		Warning     string   `json:"warning"`
		Suggestions []string `json:"suggestions"`
	} `json:"feedback"`
}

type UserMeResponse struct {
	Id                        *string      `json:"id"`
	IsAdmin                   *bool        `json:"isAdmin"`
	EmailVerified             *bool        `json:"emailVerified"`
	HasTotp                   *bool        `json:"hasTotp"`
	HasPasskeys               *bool        `json:"hasPasskeys"`
	ExpiresAt                 *interface{} `json:"expiresAt"`
	Approved                  *bool        `json:"approved"`
	Amr                       *[]string    `json:"amr"`
	CanLogin                  *bool        `json:"canLogin"`
	IsPrivilegedForTotpCreate *bool        `json:"isPrivilegedForTotpCreate"`
	IsPrivilegedForEmail      *bool        `json:"isPrivilegedForEmail"`
	IsPrivileged              *bool        `json:"isPrivileged"`
	HasEmail                  *bool        `json:"hasEmail"`
}

func (r *UserMeResponse) UnmarshalJSON(data []byte) error {
	type alias UserMeResponse
	var shadow struct {
		*alias
		EmailVerified json.RawMessage `json:"emailVerified"`
		Approved      json.RawMessage `json:"approved"`
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
	return nil
}

func rawBool(raw json.RawMessage) (bool, bool) {
	switch string(raw) {
	case "true", "1", `"1"`, `"true"`:
		return true, true
	case "false", "0", `"0"`, `"false"`:
		return false, true
	}
	var n json.Number
	if err := json.Unmarshal(raw, &n); err == nil {
		if f, err := n.Float64(); err == nil {
			return f != 0, true
		}
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		switch s {
		case "1", "true":
			return true, true
		case "0", "false":
			return false, true
		}
	}
	return false, false
}
