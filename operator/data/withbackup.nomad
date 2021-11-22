job "withbackup" {
  datacenters = ["dc1"]

  meta {
    auto-backup = true
    backup-schedule = "@daily"
    backup-target-db = "postgres"
  }

  group "cache" {
    network {
      port "db" {
        to = 6379
      }
    }

    task "redis" {
      driver = "docker"

      config {
        image = "redis:3.2"
        ports = ["db"]
      }

      resources {
        cpu    = 500
        memory = 256
      }
    }
  }
}
