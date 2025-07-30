resource "google_sql_database_instance" "postgres" {
  name             = var.instance_name
  database_version = "POSTGRES_15"
  region           = var.region

  settings {
    tier = "db-custom-1-3840"
  }
}

resource "google_sql_database" "default" {
  name     = "spiceroute"
  instance = google_sql_database_instance.postgres.name
}
