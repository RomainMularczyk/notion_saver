data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "ariga.io/atlas-provider-gorm",
    "load",
    "--path", "./src/models/saver",
    "--dialect", "postgres",
  ]
}

env "gorm" {
  src = data.external_schema.gorm.url
  dev = "docker://postgres/17/notion_saver_main_database"
  migration {
    dir = "file://migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \" \" }}"
    }
  }
}

env "local" {
  url = "postgres://localhost:5432/notion_saver_main_database?search_path=public&sslmode=disable"
  migration {
    dir = "file://src/database/migrations"
  }
}
