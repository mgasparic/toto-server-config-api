data google_client_config default {}

resource google_compute_global_address geo_distributor_ip_address {
  project    = var.us_gke_gae_project_id
  name       = "geo-distributor-ip-address"
  ip_version = "IPV4"
}

resource google_compute_global_address distributor_ip_address {
  project    = var.us_gke_gae_project_id
  name       = "distributor-ip-address"
  ip_version = "IPV4"
}

resource google_compute_global_address updater_ip_address {
  project    = var.us_gke_gae_project_id
  name       = "updater-ip-address"
  ip_version = "IPV4"
}

resource google_compute_managed_ssl_certificate geo_distributor_ssl_certificate {
  project = var.us_gke_gae_project_id
  name    = "geo-distributor-ssl-certificate"
  managed {
    domains = [var.geo_distributor_domain]
  }
}

resource google_compute_managed_ssl_certificate distributor_ssl_certificate {
  project = var.us_gke_gae_project_id
  name    = "distributor-ssl-certificate"
  managed {
    domains = [var.distributor_domain]
  }
}

resource google_app_engine_application eu_gae_app {
  project     = var.eu_gae_project_id
  location_id = var.eu_location
}

resource google_app_engine_application as_gae_app {
  project     = var.as_gae_project_id
  location_id = var.as_location
}

resource google_app_engine_application us_gae_app {
  project     = var.us_gke_gae_project_id
  location_id = var.us_location
}

resource google_container_cluster eu_cluster {
  project  = var.us_gke_gae_project_id
  name     = "eu-cluster"
  location = var.eu_location

  initial_node_count = 1
  networking_mode    = "VPC_NATIVE"
  ip_allocation_policy {}
  workload_identity_config { workload_pool = "${var.us_gke_gae_project_id}.svc.id.goog" }
  release_channel {
    channel = "STABLE"
  }
}

resource google_container_cluster as_cluster {
  project  = var.us_gke_gae_project_id
  name     = "as-cluster"
  location = var.as_location

  initial_node_count = 1
  networking_mode    = "VPC_NATIVE"
  ip_allocation_policy {}
  workload_identity_config { workload_pool = "${var.us_gke_gae_project_id}.svc.id.goog" }
  release_channel {
    channel = "STABLE"
  }
}

resource google_container_cluster us_cluster {
  project  = var.us_gke_gae_project_id
  name     = "us-cluster"
  location = var.us_location

  initial_node_count = 1
  networking_mode    = "VPC_NATIVE"
  ip_allocation_policy {}
  workload_identity_config { workload_pool = "${var.us_gke_gae_project_id}.svc.id.goog" }
  release_channel {
    channel = "STABLE"
  }
}
