output geo_distributor_ip {
  value = google_compute_global_address.geo_distributor_ip_address.address
}

output distributor_ip {
  value = google_compute_global_address.distributor_ip_address.address
}

output updater_ip {
  value = google_compute_global_address.updater_ip_address.address
}

output updater_ip_name {
  value = google_compute_global_address.updater_ip_address.name
}

output geo_distributor_cert_name {
  value = google_compute_managed_ssl_certificate.geo_distributor_ssl_certificate.name
}

output distributor_cert_name {
  value = google_compute_managed_ssl_certificate.distributor_ssl_certificate.name
}

output eu_gae_hostname {
  value = google_app_engine_application.eu_gae_app.default_hostname
}

output as_gae_hostname {
  value = google_app_engine_application.as_gae_app.default_hostname
}

output us_gae_hostname {
  value = google_app_engine_application.us_gae_app.default_hostname
}

output eu_gke_name {
  value = google_container_cluster.eu_cluster.name
}

output as_gke_name {
  value = google_container_cluster.as_cluster.name
}

output us_gke_name {
  value = google_container_cluster.us_cluster.name
}
