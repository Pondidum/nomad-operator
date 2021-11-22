job "[[ .JobID ]]" {
  datacenters = ["dc1"]

  type = "batch"

  periodic {
    cron = "[[ .Schedule ]]"
    prohibit_overlap= true
  }

  group "backup" {
    count = 1

    task "backup" {
      driver = "docker"

      config {
        image = "alpine:latest"
        command = "echo"
        args = [ "backing up [[ .SourceJobID ]]'s [[ .TargetDB ]] database" ]
      }

      // vault {
      //   policies = [ "database-backup", "s3-backup-writer" ]
      // }

      // template {
      //   data = <<EOF
      //   {{ with secret "database/creds/backup" }}
      //   PGUSER={{ .Data.username | toJSON }}
      //   PGPASSWORD={{ .Data.password | toJSON }}
      //   {{ end }}

      //   {{ with secret "secret/data/s3/backup" }}
      //   AWS_ACCESS_KEY={{ .Data.data.AWS_ACCESS_KEY | toJSON }}
      //   AWS_SECRET_KEY={{ .Data.data.AWS_SECRET_KEY | toJSON }}
      //   {{ end }}
      //   EOF
      //   destination = "secrets/db.env"
      //   env = true
      // }

      env {
        PGHOST     = "postgres.service.consul"
        PGDATABASE = "[[ .TargetDB ]]"
        AWS_REGION = "eu-west-1"
      }
    }
  }
}
