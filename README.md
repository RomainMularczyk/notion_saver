# NotionSaver

## Database configuration

### Provision pgadmin connection

We can provision pgadmin connection by adding a servers.json file in the following directory:
```
./config/pgadmin4/servers.json
```

The content of the file should be:
```json
{
  "Servers": {
    "1": {
      "Name": "NotionSaver",
      "Group": "Servers",
      "Host": "localhost",
      "Port": 5432,
      "MaintenanceDB": "postgres",
      "Username": "postgres",
      "SSLMode": "prefer",
      "DBName": "notion_saver",
      "Labels": "sslmode=prefer"
    }
  }
}
```

### Handling migrations

In order to create a new migration from the models defined in `./src/models`, we can run the following command:
```bash
atlas migrate diff --env gorm --dir "file://src/database/migrations"
```

This results in a new migration file created in `./src/database/migrations`.

To apply the migrations to the database, we can run the following command:
```bash
atlas migrate up --env gorm --dir "file://src/database/migrations"
```

## Database seeding

We can seed the database by running the following command:
```bash
atlas migrate apply --dir "file://src/database/migrations" --url "postgres://user:password@localhost:5432/database_name?sslmode=prefer"
```
