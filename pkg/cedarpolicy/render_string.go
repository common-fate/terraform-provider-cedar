package cedarpolicy

import (
	"errors"
	"fmt"
	"strings"
)

// RenderString renders a text-based representation of the policy.
func (p Policy) RenderString() (string, error) {
	var output []string

	// create annotations in the format
	// @name("value")
	for _, anno := range p.Annotations {
		name := anno.Name.ValueString()
		value := anno.Value.ValueString()

		if name == "" {
			return "", errors.New("policy annotation 'name' field must be specified")
		}
		if value == "" {
			return "", errors.New("policy annotation 'value' field must be specified")
		}

		line := fmt.Sprintf(`@%s("%s")`, name, value)
		output = append(output, line)
	}

	effect := p.Effect.ValueString()

	if effect != "permit" && effect != "forbid" {
		return "", fmt.Errorf("effect must be either 'permit' or 'forbid', got %q", effect)
	}

	// permit (
	output = append(output, fmt.Sprintf("%s (", effect))

	if p.Principal != nil {
		// principal == <entity>,

		principalType := p.Principal.Type.ValueString()
		principalID := p.Principal.ID.ValueString()

		if principalType == "" {
			return "", errors.New("principal type must be specified")
		}

		if principalID == "" {
			return "", errors.New("principal ID must be specified")
		}

		line := fmt.Sprintf("\tprincipal == %s::%q,", principalType, principalID)
		output = append(output, line)
	} else if p.PrincipalIn != nil {
		// principal == [<entity>, <entity>],

		entities := make([]string, len(*p.PrincipalIn))

		for i, ent := range *p.PrincipalIn {
			typ := ent.Type.ValueString()
			id := ent.ID.ValueString()

			if typ == "" {
				return "", fmt.Errorf("principal_in entry %v: type must be specified", i)
			}
			if id == "" {
				return "", fmt.Errorf("principal_in entry %v: ID must be specified", i)
			}
			entities[i] = fmt.Sprintf(`%s::"%s"`, typ, id)
		}

		line := fmt.Sprintf("\tprincipal in [%s],", strings.Join(entities, ", "))
		output = append(output, line)
	} else if p.PrincipalIs.ValueString() != "" {
		// principal is <entity type>,
		line := fmt.Sprintf("\tprincipal is %s,", p.PrincipalIs.ValueString())
		output = append(output, line)
	} else if p.AnyPrincipal.ValueBool() {
		// principal,
		output = append(output, "\tprincipal,")
	}

	if p.Action != nil {
		// action == <entity>,

		typ := p.Action.Type.ValueString()
		id := p.Action.ID.ValueString()

		if typ == "" {
			return "", errors.New("action type must be specified")
		}

		if id == "" {
			return "", errors.New("action ID must be specified")
		}

		line := fmt.Sprintf("\taction == %s::%q,", typ, id)
		output = append(output, line)
	} else if p.ActionIn != nil {
		// action == [<entity>, <entity>],

		entities := make([]string, len(*p.ActionIn))

		for i, ent := range *p.ActionIn {
			typ := ent.Type.ValueString()
			id := ent.ID.ValueString()

			if typ == "" {
				return "", fmt.Errorf("action_in entry %v: type must be specified", i)
			}
			if id == "" {
				return "", fmt.Errorf("action_in entry %v: ID must be specified", i)
			}
			entities[i] = fmt.Sprintf(`%s::"%s"`, typ, id)
		}

		line := fmt.Sprintf("\taction in [%s],", strings.Join(entities, ", "))
		output = append(output, line)
	} else if p.AnyAction.ValueBool() {
		// action,
		output = append(output, "\taction,")
	}

	if p.Resource != nil {
		// resource == <entity>,

		resourceType := p.Resource.Type.ValueString()
		resourceID := p.Resource.ID.ValueString()

		if resourceType == "" {
			return "", errors.New("resource type must be specified")
		}

		if resourceID == "" {
			return "", errors.New("resource ID must be specified")
		}

		line := fmt.Sprintf("\tresource == %s::%q", resourceType, resourceID)
		output = append(output, line)
	} else if p.ResourceIn != nil {
		// resource == [<entity>, <entity>],

		entities := make([]string, len(*p.ResourceIn))

		for i, ent := range *p.ResourceIn {
			typ := ent.Type.ValueString()
			id := ent.ID.ValueString()

			if typ == "" {
				return "", fmt.Errorf("resource_in entry %v: type must be specified", i)
			}
			if id == "" {
				return "", fmt.Errorf("resource_in entry %v: ID must be specified", i)
			}
			entities[i] = fmt.Sprintf(`%s::"%s"`, typ, id)
		}

		line := fmt.Sprintf("\tresource in [%s]", strings.Join(entities, ", "))
		output = append(output, line)
	} else if p.ResourceIs.ValueString() != "" {
		// resource is <entity type>,
		line := fmt.Sprintf("\tresource is %s", p.ResourceIs.ValueString())
		output = append(output, line)
	} else if p.AnyResource.ValueBool() {
		// resource,
		output = append(output, "\tresource")
	}

	// end of policy scope section
	output = append(output, ")")

	// when conditions
	for i, when := range p.When {
		text := when.Text.ValueString()
		if text == "" {
			return "", fmt.Errorf("when condition index %v: 'text' must be specified", i)
		}

		output = append(output, "when {")
		output = append(output, fmt.Sprintf("\t%s", text))
		output = append(output, "}")
	}

	// unless conditions
	for i, unless := range p.Unless {
		text := unless.Text.ValueString()
		if text == "" {
			return "", fmt.Errorf("unless condition index %v: 'text' must be specified", i)
		}

		output = append(output, "unless {")
		output = append(output, fmt.Sprintf("\t%s", text))
		output = append(output, "}")
	}

	// add semicolon to the final line
	output[len(output)-1] += ";"

	return strings.Join(output, "\n"), nil
}
