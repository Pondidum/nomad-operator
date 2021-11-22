job "operator" {
  datacenters = ["dc1"]
  type = "service"

  vault {
    policies = [ "operator-job" ]
  }

  group "operators" {

    task "operator" {
      driver = "docker"

      config {
        image = "operator:latest"
      }

      template {
        data = <<EOF
        {{ with secret "nomad/creds/operator-job" }}
        NOMAD_TOKEN={{ .Data.secret_id  | toJSON }}
        {{ end }}
        EOF
        destination = "secrets/db.env"
        env = true
      }

      env {
        NOMAD_ADDR = "nomad.service.consul"
      }
    }
  }
}
