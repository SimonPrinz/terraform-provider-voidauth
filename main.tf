terraform {
  required_providers {
    voidauth = {
      source = "simonprinz/voidauth"
    }
  }
}

provider "voidauth" {
  url      = "http://localhost:3000"
  username = "auth_admin"
  password = "Voidauth!23"
}

data "voidauth_user" "auth_admin" {
  id = "982a9b7c-089a-494c-af04-213d8b444608"
}

data "voidauth_group" "auth_admins" {
  id = "137be4f1-a943-4009-baf8-db540f22112d"
}

resource "voidauth_group" "test" {
  name = "test_group"
  users = [
    {
      id = data.voidauth_user.auth_admin.id
      username = data.voidauth_user.auth_admin.username
    }
  ]
}

resource "voidauth_proxyauth" "test" {
  domain = "test.com"
  max_session_length = 20
  groups = [
    voidauth_group.test.name,
    data.voidauth_group.auth_admins.name
  ]
}
