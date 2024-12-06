
variable "DB_URL" {
  type = string
}

locals {
  src_files        = ["file://internal/database/schema.hcl"]
  enterprise_files = ["file://database/enterprise.schema.hcl"]
}

// Optional flag to include/exclude enterprise schema
variable "USE_ENTERPRISE_SCHEMA" {
  type    = bool
  default = false
}

// Define an environment named "local"
env "global" {
  src = var.USE_ENTERPRISE_SCHEMA ? concat(local.src_files, local.enterprise_files) : local.src_files

  url = var.DB_URL
  dev = var.DB_URL

  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
  migration {
    // URL where the migration directory resides.
    dir = "file://internal/database/migrations"
  }
}
