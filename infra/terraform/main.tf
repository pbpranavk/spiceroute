provider "google" {
  project = var.project_id
  region  = var.region
}

module "gke" {
  source       = "./gke"
  project_id   = var.project_id
  region       = var.region
  cluster_name = "spiceroute-cluster"
}

module "cloudsql" {
  source        = "./cloudsql"
  project_id    = var.project_id
  region        = var.region
  instance_name = "spiceroute-pg"
}

module "artifact_registry" {
  source     = "./artifact_registry"
  project_id = var.project_id
  location   = var.region
  repo_name  = "spiceroute"
}

module "pubsub" {
  source     = "./pubsub"
  project_id = var.project_id
}
