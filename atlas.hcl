
variable "DB_URL" {
  type    = string
  default = "postgres://sarthakjdev@localhost:5432/wapikit?sslmode=disable"
}

locals {
  src_files = ["file://database/schema.hcl"]
}

// Optional flag to include/exclude enterprise schema
variable "USE_ENTERPRISE_SCHEMA" {
  type    = bool
  default = false
}

// Define an environment named "local"
env "local" {
  src = var.USE_ENTERPRISE_SCHEMA ? concat(local.src_files, ["file://database/enterprise.schema.hcl"]) : local.src_files

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
