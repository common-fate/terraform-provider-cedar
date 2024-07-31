package cedarpolicy

import (
	"testing"

	"github.com/common-fate/terraform-provider-cedar/pkg/eid"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestRenderString(t *testing.T) {

	tests := []struct {
		name   string
		policy Policy
		want   string
	}{

		{
			name: "allow_all",
			policy: Policy{
				Effect:       types.StringValue("permit"),
				AnyPrincipal: types.BoolValue(true),
				AnyAction:    types.BoolValue(true),
				AnyResource:  types.BoolValue(true),
			},
			want: `permit (
	principal,
	action,
	resource
);`,
		},
		{
			name: "allow_all_with_advice",
			policy: Policy{
				Effect: types.StringValue("permit"),
				Annotations: []Annotation{
					{
						Name:  types.StringValue("advice"),
						Value: types.StringValue("test"),
					},
				},
				AnyPrincipal: types.BoolValue(true),
				AnyAction:    types.BoolValue(true),
				AnyResource:  types.BoolValue(true),
			},
			want: `@advice("test")
permit (
	principal,
	action,
	resource
);`,
		},
		{
			name: "principal_action_resource_equals",
			policy: Policy{
				Effect: types.StringValue("permit"),
				Principal: &eid.EID{
					Type: types.StringValue("CF::User"),
					ID:   types.StringValue("user1"),
				},
				Action: &eid.EID{
					Type: types.StringValue("Action::Access"),
					ID:   types.StringValue("Request"),
				},
				Resource: &eid.EID{
					Type: types.StringValue("Test::Vault"),
					ID:   types.StringValue("test1"),
				},
			},
			want: `permit (
	principal == CF::User::"user1",
	action == Action::Access::"Request",
	resource == Test::Vault::"test1"
);`,
		},
		{
			name: "when",
			policy: Policy{
				Effect: types.StringValue("permit"),
				Principal: &eid.EID{
					Type: types.StringValue("CF::User"),
					ID:   types.StringValue("user1"),
				},
				Action: &eid.EID{
					Type: types.StringValue("Action::Access"),
					ID:   types.StringValue("Request"),
				},
				Resource: &eid.EID{
					Type: types.StringValue("Test::Vault"),
					ID:   types.StringValue("test1"),
				},
				When: []Condition{
					{Text: types.StringValue("true")},
				},
			},
			want: `permit (
	principal == CF::User::"user1",
	action == Action::Access::"Request",
	resource == Test::Vault::"test1"
)
when {
	true
};`,
		},

		{
			name: "unless",
			policy: Policy{
				Effect: types.StringValue("permit"),
				Principal: &eid.EID{
					Type: types.StringValue("CF::User"),
					ID:   types.StringValue("user1"),
				},
				Action: &eid.EID{
					Type: types.StringValue("Action::Access"),
					ID:   types.StringValue("Request"),
				},
				Resource: &eid.EID{
					Type: types.StringValue("Test::Vault"),
					ID:   types.StringValue("test1"),
				},

				Unless: []Condition{
					{Text: types.StringValue("true")},
				},
			},
			want: `permit (
	principal == CF::User::"user1",
	action == Action::Access::"Request",
	resource == Test::Vault::"test1"
)
unless {
	true
};`,
		},

		{
			name: "principal_action_resource_in",
			policy: Policy{
				Effect: types.StringValue("permit"),
				PrincipalIn: &eid.EID{
					Type: types.StringValue("CF::User"),
					ID:   types.StringValue("user1"),
				},

				ActionIn: &[]eid.EID{
					{
						Type: types.StringValue("Action::Access"),
						ID:   types.StringValue("Request"),
					},
				},

				ResourceIn: &eid.EID{
					Type: types.StringValue("Test::Vault"),
					ID:   types.StringValue("test1"),
				},
			},
			want: `permit (
	principal in CF::User::"user1",
	action in [Action::Access::"Request"],
	resource in Test::Vault::"test1"
);`,
		},
		{
			name: "in_condition_multiple_values",
			policy: Policy{
				Effect: types.StringValue("permit"),
				PrincipalIn: &eid.EID{
					Type: types.StringValue("CF::User"),
					ID:   types.StringValue("user1"),
				},

				ActionIn: &[]eid.EID{
					{
						Type: types.StringValue("Action::Access"),
						ID:   types.StringValue("Request"),
					},
					{
						Type: types.StringValue("Action::Access"),
						ID:   types.StringValue("Close"),
					},
				},

				ResourceIn: &eid.EID{
					Type: types.StringValue("Test::Vault"),
					ID:   types.StringValue("test1"),
				},
			},
			want: `permit (
	principal in CF::User::"user1",
	action in [Action::Access::"Request", Action::Access::"Close"],
	resource in Test::Vault::"test1"
);`,
		},
		{
			name: "principal_resource_is",
			policy: Policy{
				Effect:      types.StringValue("permit"),
				PrincipalIs: types.StringValue("CF::User"),

				Action: &eid.EID{
					Type: types.StringValue("Action::Access"),
					ID:   types.StringValue("Request"),
				},

				ResourceIs: types.StringValue("Test::Vault"),
			},
			want: `permit (
	principal is CF::User,
	action == Action::Access::"Request",
	resource is Test::Vault
);`,
		},
		{
			name: "test having multiple when conditions",
			policy: Policy{
				Effect: types.StringValue("permit"),
				Principal: &eid.EID{
					Type: types.StringValue("CF::User"),
					ID:   types.StringValue("user1"),
				},
				Action: &eid.EID{
					Type: types.StringValue("Action::Access"),
					ID:   types.StringValue("Request"),
				},
				Resource: &eid.EID{
					Type: types.StringValue("Test::Vault"),
					ID:   types.StringValue("test1"),
				},
				When: []Condition{
					{
						Text: types.StringValue("true"),
					},
					{
						Text: types.StringValue("test2"),
					},
					{
						Text: types.StringValue("test3"),
					},
				},
			},
			want: `permit (
	principal == CF::User::"user1",
	action == Action::Access::"Request",
	resource == Test::Vault::"test1"
)
when {
	true
}
when {
	test2
}
when {
	test3
};`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := tt.policy.RenderString()
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
