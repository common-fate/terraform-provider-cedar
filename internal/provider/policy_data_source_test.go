package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestPolicyDataSource_Simple(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_5_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				data "cedar_policyset" "test" {
					policy {
						effect = "permit"
						any_principal = true
						any_action = true
						any_resource = true
					}
				}

				output "test" {
					value = data.cedar_policyset.test.text
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckOutput("test", "permit (\n\tprincipal,\n\taction,\n\tresource\n);\n"),
				),
			},
		},
	})
}

func TestPolicyDataSource_Annotations(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_5_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				data "cedar_policyset" "test" {
					policy {
						annotation {
							name = "advice"
							value = "test"
						}

						effect = "permit"
						any_principal = true
						any_action = true
						any_resource = true
					}
				}

				output "test" {
					value = data.cedar_policyset.test.text
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckOutput("test", "@advice(\"test\")\npermit (\n\tprincipal,\n\taction,\n\tresource\n);\n"),
				),
			},
		},
	})
}
