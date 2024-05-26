data "external_schema" "gorm" {
  program = ["go", "run", "./cmd/migrate/main.go"]
}

variable "database_url" {
  type = string
  default = "postgres://sarthakjdev@localhost:5432/wapikit?sslmode=disable"
}

env "gorm" {
  src = data.external_schema.gorm.url
  dev = var.database_url  
  url = var.database_url
  migration {
    dir = "file://database/migrations" 
  }

  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}