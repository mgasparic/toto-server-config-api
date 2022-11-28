data google_client_config default {}

data google_billing_account account {
  display_name = "My Billing Account"
  open         = true
}

resource google_project eu_gae_project {
  name            = "EU GAE"
  project_id      = "eu-gae-${var.project_suffix}"
  billing_account = data.google_billing_account.account.id
}

resource google_project_service eu_gae_cloudbuild_api {
  project                    = google_project.eu_gae_project.project_id
  service                    = "cloudbuild.googleapis.com"
  disable_dependent_services = true
}

resource google_project as_gae_project {
  name            = "AS GAE"
  project_id      = "as-gae-${var.project_suffix}"
  billing_account = data.google_billing_account.account.id
}

resource google_project_service as_gae_cloudbuild_api {
  project                    = google_project.as_gae_project.project_id
  service                    = "cloudbuild.googleapis.com"
  disable_dependent_services = true
}

resource google_project us_gke_gae_project {
  name            = "US GKE GAE"
  project_id      = "us-gke-gae-${var.project_suffix}"
  billing_account = data.google_billing_account.account.id
}

resource google_project_service us_gke_gae_cloudbuild_api {
  project                    = google_project.us_gke_gae_project.project_id
  service                    = "cloudbuild.googleapis.com"
  disable_dependent_services = true
}

resource google_project_service us_gke_gae_compute_api {
  project                    = google_project.us_gke_gae_project.project_id
  service                    = "compute.googleapis.com"
  disable_dependent_services = true
}

resource google_project_service us_gke_gae_container_api {
  project                    = google_project.us_gke_gae_project.project_id
  service                    = "container.googleapis.com"
  disable_dependent_services = true
}

resource google_project_service us_gke_gae_containerregistry_api {
  project                    = google_project.us_gke_gae_project.project_id
  service                    = "containerregistry.googleapis.com"
  disable_dependent_services = true
}

resource google_project_service us_gke_gae_gkehub_api {
  project                    = google_project.us_gke_gae_project.project_id
  service                    = "gkehub.googleapis.com"
  disable_dependent_services = true
}
