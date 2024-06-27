module "team1" {
    source = "./github_team_membership.tf"
    name = "Team1"

    members = ["user1", "user3", "user5", "user6", "user7", "user8", "user9", "user10"]
    maintainers = ["user2"]
} 

module "team2" {
    source = "./github_team_membership.tf"
    name = "Team2"

    members = ["user2", "user3", "user4", "user6", "user7", "user8", "user9", "user10"]
    maintainers = ["user1"]
}

module "team3" {
  source = "./github_team_membership.tf"
  name = "Team3"

  members = ["user1", "user2", "user4", "user5", "user6", "user7"]
  maintainers = ["user3"]
}

module "team4" {
    source = "./github_team_membership.tf"
    name = "Team4"

    members = ["user1", "user3", "user4", "user5", "user6", "user7", "user8", "user9", "user10"]
    maintainers = ["user2"]
}

module "team5" {
    source = "./github_team_membership.tf"
    name = "Team5"

    members = ["user2", "user3", "user4", "user6", "user7", "user8", "user9", "user10"]
    maintainers = ["user1"]
}

module "team6" {
    source = "./github_team_membership.tf"
    name = "Team6"

    members = ["user1", "user2", "user3", "user5", "user7", "user8", "user9", "user10"]
    maintainers = ["user4"]
}
