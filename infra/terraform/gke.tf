resource "google_container_cluster" "primary" {
  name                     = var.cluster_name
  location                 = var.region
  remove_default_node_pool = true
  initial_node_count       = 1
}

resource "google_container_node_pool" "primary_nodes" {
  cluster    = google_container_cluster.primary.name
  location   = var.region
  node_count = 3
  node_config {
    machine_type = "e2-standard-4"
    oauth_scopes = ["https://www.googleapis.com/auth/cloud-platform"]
  }
}
