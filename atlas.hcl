
variable "DB_URL" {
  type    = string
  default = "postgres://sarthakjdev@localhost:5432/wapikit?sslmode=disable"
}

// Define an environment named "local"
env "local" {
  src = "file://database/schema.hcl"
  url = var.DB_URL
  dev = var.DB_URL
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
  migration {
    // URL where the migration directory resides.
    dir = "file://database/migrations"
  }
}
