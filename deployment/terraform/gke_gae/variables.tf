variable eu_gae_project_id {
  description = "Project ID for EU GAE"
  type        = string
}

variable as_gae_project_id {
  description = "Project ID for AS GAE"
  type        = string
}

variable us_gke_gae_project_id {
  description = "Project ID for US GKE GAE"
  type        = string
}

variable eu_location {
  description = "Location for EU GKE and GAE"
  type        = string
}

variable as_location {
  description = "Location for AS GKE and GAE"
  type        = string
}

variable us_location {
  description = "Location for US GKE and GAE"
  type        = string
}

variable geo_distributor_domain {
  description = "Domain name for geo-distributor endpoint"
  type        = string
}

variable distributor_domain {
  description = "Domain name for distributor endpoint"
  type        = string
}

variable updater_domain {
  description = "Domain name for updater endpoint"
  type        = string
}
