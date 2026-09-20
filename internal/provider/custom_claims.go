package provider

import (
	"github.com/SimonPrinz/terraform-provider-voidauth/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CustomClaimModel struct {
	Claim types.String `tfsdk:"claim"`
	Value types.String `tfsdk:"value"`
}

func claimsToModels(claims *[]client.AdminClaim) []CustomClaimModel {
	if claims == nil {
		return nil
	}
	out := make([]CustomClaimModel, 0, len(*claims))
	for _, c := range *claims {
		m := CustomClaimModel{}
		if c.Claim != nil {
			m.Claim = types.StringValue(*c.Claim)
		}
		if c.Value != nil {
			m.Value = types.StringValue(*c.Value)
		}
		out = append(out, m)
	}
	return out
}

func modelsToClaims(models []CustomClaimModel) []client.AdminClaim {
	out := make([]client.AdminClaim, 0, len(models))
	for _, m := range models {
		if m.Claim.IsNull() || m.Claim.IsUnknown() {
			continue
		}
		claim := m.Claim.ValueString()
		value := m.Value.ValueString()
		out = append(out, client.AdminClaim{Claim: &claim, Value: &value})
	}
	return out
}
