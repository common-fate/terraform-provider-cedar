package eid

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	dataSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var EIDAttrs = map[string]schema.Attribute{
	"type": schema.StringAttribute{
		Required:            true,
		MarkdownDescription: "The entity type",
	},
	"id": schema.StringAttribute{
		Required:            true,
		MarkdownDescription: "The entity ID",
	},
}

var EIDAttrsForDataSource = map[string]attr.Type{
	"type": types.StringType,
	"id":   types.StringType,
}

var EIDAttributesForDataSource = map[string]dataSchema.Attribute{
	"type": dataSchema.StringAttribute{
		Required:            true,
		MarkdownDescription: "The entity type",
	},
	"id": dataSchema.StringAttribute{
		Required:            true,
		MarkdownDescription: "The entity ID",
	},
}

type EID struct {
	Type types.String `tfsdk:"type"`
	ID   types.String `tfsdk:"id"`
}
