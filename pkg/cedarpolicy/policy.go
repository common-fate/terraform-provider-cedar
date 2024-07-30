package cedarpolicy

import (
	"github.com/common-fate/terraform-provider-cedar/pkg/eid"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Condition struct {
	Text types.String `tfsdk:"text"`
}

type Annotation struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

type Policy struct {
	Effect      types.String `tfsdk:"effect"`
	Annotations []Annotation `tfsdk:"annotation"`

	AnyPrincipal types.Bool   `tfsdk:"any_principal"`
	Principal    *eid.EID     `tfsdk:"principal"`
	PrincipalIn  *[]eid.EID   `tfsdk:"principal_in"`
	PrincipalIs  types.String `tfsdk:"principal_is"`

	AnyAction types.Bool `tfsdk:"any_action"`
	Action    *eid.EID   `tfsdk:"action"`
	ActionIn  *[]eid.EID `tfsdk:"action_in"`

	AnyResource types.Bool   `tfsdk:"any_resource"`
	Resource    *eid.EID     `tfsdk:"resource"`
	ResourceIn  *[]eid.EID   `tfsdk:"resource_in"`
	ResourceIs  types.String `tfsdk:"resource_is"`

	When   []Condition `tfsdk:"when"`
	Unless []Condition `tfsdk:"unless"`
}
