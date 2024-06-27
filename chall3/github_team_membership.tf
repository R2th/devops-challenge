variable "maintainers" {
  type        = set(string)
  default     = []
}

variable "members" {
  type        = set(string)
  default     = []
}

variable "description" {
    type = string
    default = ""
}

resource "github_team" "team" {
    name = var.name
    description = var.description
    privacy = var.privacy
    parent_team_id = var.parent_team_id
}

locals {
    maintainers = { for i in var.maintainers : lower(i) => { role = "maintainer", username = i}}
    members = { for i in setsubtract(var.members, var.maintainers) : lower(i) => { role = "member", username = i } }

    memberships = merge(local.maintainers, local.members)
}

resource "github_team_membership" "team_membership" {
    for_each = local.memberships

    team_id = try(github_team.team[0].id, null)

    username = each.value.username
    role = each.value.role
}