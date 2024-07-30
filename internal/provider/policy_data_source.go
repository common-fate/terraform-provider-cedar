package provider

import (
	"context"
	"strings"

	"github.com/common-fate/terraform-provider-cedar/pkg/cedarpolicy"
	"github.com/common-fate/terraform-provider-cedar/pkg/eid"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ datasource.DataSource = &PolicyDataSource{}

type PolicyDataSource struct{}

func NewPolicyDataSource() datasource.DataSource {
	return &PolicyDataSource{}
}

type PolicyDataSourceModel struct {
	Policies []cedarpolicy.Policy `tfsdk:"policy"`
	Text     types.String         `tfsdk:"text"`
}

func (d *PolicyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policyset"
}

func (d *PolicyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Generates a Cedar Policy Set for use with resources that expect Cedar policies.",
		MarkdownDescription: `Generates a Cedar Policy Set for use with resources that expect Cedar policies.

You may specify multiple Cedar policies in a PolicySet by providing multiple 'policy' blocks.

Each policy block must contain an effect (either 'permit' or 'forbid') and must contain a principal, action, and resource clause for the policy scope.

For the principal clause, you can provide 'principal', 'principal_in', 'principal_is', or 'any_principal'.
For the action clause, you can provide 'action', 'action_in', or 'any_action'.
For the resource clause, you can provide 'resource', 'resource_in', 'resource_is', or 'any_resource'.

You may also optionally provide one or more 'when' and 'unless' conditions as blocks.
`,
		Attributes: map[string]schema.Attribute{
			"text": schema.StringAttribute{
				MarkdownDescription: "The Cedar PolicySet, rendered as a string.",
				Computed:            true,
			},
		},
		Blocks: map[string]schema.Block{
			"policy": schema.ListNestedBlock{
				MarkdownDescription: "a list of policies to be included in the PolicySet",

				NestedObject: schema.NestedBlockObject{

					Attributes: map[string]schema.Attribute{
						"effect": schema.StringAttribute{
							MarkdownDescription: "Must be either 'permit' or 'forbid'.",
							Required:            true,
						},
						"annotation": schema.SingleNestedAttribute{
							MarkdownDescription: "Additional application-specific metadata attached to Cedar policies.",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "The name of the @decorator, eg. @advice()",
									Required:    true,
								},
								"value": schema.StringAttribute{
									Description: "The value of the @decorator, eg. @advice(value)",
									Required:    true,
								},
							},
						},
						"any_principal": schema.BoolAttribute{
							MarkdownDescription: "Specifies the principal component of the policy scope. Matches all principals. Equivalent to writing 'principal'",
							Optional:            true,
						},
						"principal": schema.ObjectAttribute{
							MarkdownDescription: "Specifies the principal component of the policy scope. Equivalent to writing 'principal =='",
							Optional:            true,
							AttributeTypes:      eid.EIDAttrsForDataSource,
						},
						"principal_is": schema.StringAttribute{
							MarkdownDescription: "Specifies the principal component of the policy scope. Equivalent to writing 'principal in'",
							Optional:            true,
						},
						"principal_in": schema.ListAttribute{
							MarkdownDescription: "Specifies the principal component of the policy scope. Equivalent to writing 'principal =='",
							Optional:            true,
							ElementType: basetypes.ObjectType{
								AttrTypes: eid.EIDAttrsForDataSource,
							},
						},

						"any_action": schema.BoolAttribute{
							MarkdownDescription: "Specifies the action component of the policy scope. Matches all actions. Equivalent to writing 'action'",
							Optional:            true,
						},
						"action": schema.ObjectAttribute{
							MarkdownDescription: "Specifies the action component of the policy scope. Equivalent to writing 'action =='",
							Optional:            true,
							AttributeTypes:      eid.EIDAttrsForDataSource,
						},
						"action_in": schema.ListAttribute{
							MarkdownDescription: "Specifies the action component of the policy scope. Equivalent to writing 'action in'",
							Optional:            true,
							ElementType: basetypes.ObjectType{
								AttrTypes: eid.EIDAttrsForDataSource,
							},
						},

						"any_resource": schema.BoolAttribute{
							MarkdownDescription: "Specifies the resource component of the policy scope. Matches all resources. Equivalent to writing 'resource'",
							Optional:            true,
						},
						"resource": schema.ObjectAttribute{
							MarkdownDescription: "Specifies the resource component of the policy scope. Equivalent to writing 'resource =='",
							Optional:            true,
							AttributeTypes:      eid.EIDAttrsForDataSource,
						},
						"resource_is": schema.StringAttribute{
							MarkdownDescription: "Specifies the resource component of the policy scope. Equivalent to writing 'resource is'",
							Optional:            true,
						},
						"resource_in": schema.ListAttribute{
							MarkdownDescription: "Specifies the resource component of the policy scope. Equivalent to writing 'resource in'",
							Optional:            true,
							ElementType: basetypes.ObjectType{
								AttrTypes: eid.EIDAttrsForDataSource,
							},
						},
					},
					Blocks: map[string]schema.Block{
						"when": schema.ListNestedBlock{
							MarkdownDescription: "Defines additional conditions under which the policy applies. The 'when' block must evaluate to true, otherwise the policy does not apply.",

							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"text": schema.StringAttribute{
										MarkdownDescription: "when can be used with the text attribute to define the when clause in plain-text.",
										Required:            true,
									},
								},
							},
						},

						"unless": schema.ListNestedBlock{
							MarkdownDescription: "Defines additional conditions under which the policy applies. The 'when' block must evaluate to true, otherwise the policy does not apply.",

							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"text": schema.StringAttribute{
										MarkdownDescription: "unless can be used with the text attribute to define the when clause in plain-text.",
										Required:            true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *PolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PolicyDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	renderedPolicies := make([]string, len(data.Policies))
	for i, policy := range data.Policies {

		var principalClauses, actionClauses, resourceClauses int

		if policy.Principal != nil {
			principalClauses++
		}

		if policy.PrincipalIn != nil {
			principalClauses++
		}

		if policy.PrincipalIs.ValueString() != "" {
			principalClauses++
		}

		if policy.AnyPrincipal.ValueBool() {
			principalClauses++
		}

		if principalClauses == 0 {
			resp.Diagnostics.AddError(
				"Unable to Create data source: Cedar PolicySet",
				"a principal clause must be specified, one of: principal, principal_in, principal_is, any_principal",
			)
		}

		if principalClauses > 1 {
			resp.Diagnostics.AddError(
				"Unable to Create data source: Cedar PolicySet",
				"only one principal clause must be specified, one of: principal, principal_in, principal_is, any_principal",
			)
		}

		if policy.Action != nil {
			actionClauses++
		}

		if policy.ActionIn != nil {
			actionClauses++
		}

		if policy.AnyAction.ValueBool() {
			actionClauses++
		}

		if actionClauses == 0 {
			resp.Diagnostics.AddError(
				"Unable to Create data source: Cedar PolicySet",
				"a action clause must be specified, one of: action, action_in, action_is, any_action",
			)
		}

		if actionClauses > 1 {
			resp.Diagnostics.AddError(
				"Unable to Create data source: Cedar PolicySet",
				"only one action clause must be specified, one of: action, action_in, action_is, any_action",
			)
		}

		if policy.Resource != nil {
			resourceClauses++
		}

		if policy.ResourceIn != nil {
			resourceClauses++
		}

		if policy.ResourceIs.ValueString() != "" {
			resourceClauses++
		}

		if policy.AnyResource.ValueBool() {
			resourceClauses++
		}

		if resourceClauses == 0 {
			resp.Diagnostics.AddError(
				"Unable to Create data source: Cedar PolicySet",
				"a resource clause must be specified, one of: resource, resource_in, resource_is, any_resource",
			)
		}

		if resourceClauses > 1 {
			resp.Diagnostics.AddError(
				"Unable to Create data source: Cedar PolicySet",
				"only one resource clause must be specified, one of: resource, resource_in, resource_is, any_resource",
			)
		}

		currentPolicy, err := policy.RenderString()
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Create data source: Cedar PolicySet",
				"An unexpected error occurred while parsing the policy. "+
					"Please report this issue to the provider developers.\n\n"+
					"JSON Error: "+err.Error(),
			)

			return
		}
		renderedPolicies[i] = currentPolicy
	}

	allPolicies := strings.Join(renderedPolicies, "\n\n") + "\n"

	data.Text = types.StringValue(allPolicies)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
