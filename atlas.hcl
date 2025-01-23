
variable "DB_URL" {
  type = string
}

locals {
  src_files        = ["file://internal/database/schema.hcl"]
  enterprise_files = ["file://.enterprise/internal/database/cloud.schema.hcl"]
}
env "global" {
  src = local.src_files

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

env "managed_cloud" {
  src = concat(local.src_files, local.enterprise_files)

  url = var.DB_URL
  dev = var.DB_URL

  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
  migration {
    // URL where the migration directory resides.
    dir = "file://.enterprise/internal/database/migrations"
  }
}
