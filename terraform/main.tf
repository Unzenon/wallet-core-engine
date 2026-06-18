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

resource "docker_image" "migrate_image" {
  name         = "migrate/migrate"
  keep_locally = true
}

resource "docker_container" "migration" {
  name  = "tf_migration"
  image = docker_image.migrate_image.image_id

  volumes {
    host_path      = abspath("${path.module}/../migrations")
    container_path = "/migrations"
  }

  command = [
    "-path",
    "/migrations",
    "-database",
    "postgres://${var.db_user}:${var.db_password}@tf_postgres_db:5432/${var.db_name}?sslmode=disable",
    "up"
  ]

  networks_advanced {
    name = docker_network.wallet_net.name
  }

  depends_on = [docker_container.db]
}