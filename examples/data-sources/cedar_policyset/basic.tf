# Basic example
data "cedar_policyset" "example" {
  policy {
    principal = {
      type = "User"
      id   = "alice"
    }
    action = {
      type = "Action"
      id   = "Read"
    }
    any_resource = true

    when {
      text = "resource.is_public"
    }
  }
}

output "policy_text" {
  // renders the following policy:
  //
  // permit (
  //  principal == User::"alice",
  //  action == Action::"Read",
  //  resource
  // );
  value = data.cedar_policyset.example.text
}

# With when/unless conditions
data "cedar_policyset" "with_conditions" {
  policy {
    principal = {
      type = "User"
      id   = "alice"
    }
    action = {
      type = "Action"
      id   = "Read"
    }
    any_resource = true

    when {
      text = "resource.is_public"
    }
  }
}

output "policy_text" {
  // renders the following policy:
  //
  // permit (
  //  principal == User::"alice",
  //  action == Action::"Read",
  //  resource
  // );
  value = data.cedar_policyset.example.text
}
