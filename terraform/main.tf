terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "~> 3.0.1"
    }
  }
}

provider "docker" {}

resource "docker_network" "wallet_net" {
  name = "tf_wallet_network"
}

resource "docker_image" "postgres" {
  name         = "postgres:16-alpine"
  keep_locally = true
}

resource "docker_container" "db" {
  name  = "tf_postgres_db"
  image = docker_image.postgres.image_id
  
  env = [
    "POSTGRES_USER=${var.db_user}",
    "POSTGRES_PASSWORD=${var.db_password}",
    "POSTGRES_DB=${var.db_name}"
  ]
  
  ports {
    internal = 5432
    external = 5432
  }
  
  networks_advanced {
    name = docker_network.wallet_net.name
  }
}