#!/bin/bash
set -e

cd $(dirname $0)

export TF_VAR_project_suffix=$1

export TF_VAR_eu_location="europe-west2"
export TF_VAR_as_location="asia-east2"
export TF_VAR_us_location="us-west1"

export TF_VAR_geo_distributor_domain="geo-distributor.sharedtensors.com."
export TF_VAR_distributor_domain="distributor.sharedtensors.com."
export TF_VAR_updater_domain="updater.sharedtensors.com."

gcloud auth login
gcloud auth configure-docker

# infrastructure deployment
TERRAFORM_PROJECTS_PERMISSIONS_DIRECTORY=../deployment/terraform/projects_permissions
terraform -chdir=$TERRAFORM_PROJECTS_PERMISSIONS_DIRECTORY init
terraform -chdir=$TERRAFORM_PROJECTS_PERMISSIONS_DIRECTORY apply -auto-approve
export TF_VAR_eu_gae_project_id=$(terraform -chdir=$TERRAFORM_PROJECTS_PERMISSIONS_DIRECTORY output -raw eu_gae_project_id)
export TF_VAR_as_gae_project_id=$(terraform -chdir=$TERRAFORM_PROJECTS_PERMISSIONS_DIRECTORY output -raw as_gae_project_id)
export TF_VAR_us_gke_gae_project_id=$(terraform -chdir=$TERRAFORM_PROJECTS_PERMISSIONS_DIRECTORY output -raw us_gke_gae_project_id)
sleep 30 # wait a bit for projects to get ready and permissions to take effect
TERRAFORM_GKE_GAE_DIRECTORY=../deployment/terraform/gke_gae
terraform -chdir=$TERRAFORM_GKE_GAE_DIRECTORY init
terraform -chdir=$TERRAFORM_GKE_GAE_DIRECTORY apply -auto-approve
GEO_DISTRIBUTOR_IP=$(terraform -chdir=$TERRAFORM_GKE_GAE_DIRECTORY output -raw geo_distributor_ip)
DISTRIBUTOR_IP=$(terraform -chdir=$TERRAFORM_GKE_GAE_DIRECTORY output -raw distributor_ip)
UPDATER_IP=$(terraform -chdir=$TERRAFORM_GKE_GAE_DIRECTORY output -raw updater_ip)
UPDATER_IP_NAME=$(terraform -chdir=$TERRAFORM_GKE_GAE_DIRECTORY output -raw updater_ip_name)
GEO_DISTRIBUTOR_CERT_NAME=$(terraform -chdir=$TERRAFORM_GKE_GAE_DIRECTORY output -raw geo_distributor_cert_name)
DISTRIBUTOR_CERT_NAME=$(terraform -chdir=$TERRAFORM_GKE_GAE_DIRECTORY output -raw distributor_cert_name)
EU_GAE_HOSTNAME=$(terraform -chdir=$TERRAFORM_GKE_GAE_DIRECTORY output -raw eu_gae_hostname)
AS_GAE_HOSTNAME=$(terraform -chdir=$TERRAFORM_GKE_GAE_DIRECTORY output -raw as_gae_hostname)
US_GAE_HOSTNAME=$(terraform -chdir=$TERRAFORM_GKE_GAE_DIRECTORY output -raw us_gae_hostname)
EU_GKE_NAME=$(terraform -chdir=$TERRAFORM_GKE_GAE_DIRECTORY output -raw eu_gke_name)
AS_GKE_NAME=$(terraform -chdir=$TERRAFORM_GKE_GAE_DIRECTORY output -raw as_gke_name)
US_GKE_NAME=$(terraform -chdir=$TERRAFORM_GKE_GAE_DIRECTORY output -raw us_gke_name)
echo "IPs:"
echo "geo-distributor: $GEO_DISTRIBUTOR_IP"
echo "distributor: $DISTRIBUTOR_IP"
echo "updater: $UPDATER_IP"
sleep 30 # wait a bit for GAE to get ready, you can connect domains to IPs in the meantime

# builds
DISTRIBUTOR_TAG=gcr.io/$TF_VAR_us_gke_gae_project_id/zippo.apps/distributor:1.0.0
docker build --build-arg service=distributor --tag $DISTRIBUTOR_TAG --file ../build/package/service/Dockerfile ..
docker push $DISTRIBUTOR_TAG
UPDATER_TAG=gcr.io/$TF_VAR_us_gke_gae_project_id/zippo.apps/updater:1.0.0
docker build --build-arg service=updater --tag $UPDATER_TAG --file ../build/package/service/Dockerfile ..
docker push $UPDATER_TAG
cp ../api/nginx.temp.conf ../api/nginx-eu.conf
sed -i -e "s/{{.AppUrl}}/$EU_GAE_HOSTNAME/g" ../api/nginx-eu.conf
EU_PROXY_TAG=gcr.io/$TF_VAR_us_gke_gae_project_id/zippo.apps/eu-proxy:1.0.0
docker build --build-arg region=eu --tag $EU_PROXY_TAG --file ../build/package/nginx/Dockerfile ../api
docker push $EU_PROXY_TAG
cp ../api/nginx.temp.conf ../api/nginx-as.conf
sed -i -e "s/{{.AppUrl}}/$AS_GAE_HOSTNAME/g" ../api/nginx-as.conf
AS_PROXY_TAG=gcr.io/$TF_VAR_us_gke_gae_project_id/zippo.apps/as-proxy:1.0.0
docker build --build-arg region=as --tag $AS_PROXY_TAG --file ../build/package/nginx/Dockerfile ../api
docker push $AS_PROXY_TAG
cp ../api/nginx.temp.conf ../api/nginx-us.conf
sed -i -e "s/{{.AppUrl}}/$US_GAE_HOSTNAME/g" ../api/nginx-us.conf
US_PROXY_TAG=gcr.io/$TF_VAR_us_gke_gae_project_id/zippo.apps/us-proxy:1.0.0
docker build --build-arg region=us --tag $US_PROXY_TAG --file ../build/package/nginx/Dockerfile ../api
docker push $US_PROXY_TAG

# application deployment
gcloud container fleet memberships register gke-eu --gke-cluster $TF_VAR_eu_location/$EU_GKE_NAME --enable-workload-identity --project=$TF_VAR_us_gke_gae_project_id
gcloud container fleet memberships register gke-as --gke-cluster $TF_VAR_as_location/$AS_GKE_NAME --enable-workload-identity --project=$TF_VAR_us_gke_gae_project_id
gcloud container fleet memberships register gke-us --gke-cluster $TF_VAR_us_location/$US_GKE_NAME --enable-workload-identity --project=$TF_VAR_us_gke_gae_project_id
gcloud container fleet ingress enable --config-membership=gke-us --project=$TF_VAR_us_gke_gae_project_id
gcloud container clusters get-credentials $EU_GKE_NAME --zone=$TF_VAR_eu_location --project=$TF_VAR_us_gke_gae_project_id
helm install -f ../deployment/helm/static-values/remote-config.yaml remote-config ../deployment/helm/deployment --set instance.replicas=3 --set docker.image=$DISTRIBUTOR_TAG
helm install -f ../deployment/helm/static-values/proxy.yaml proxy ../deployment/helm/deployment --set instance.replicas=3 --set docker.image=$EU_PROXY_TAG
gcloud container clusters get-credentials $AS_GKE_NAME --zone=$TF_VAR_as_location --project=$TF_VAR_us_gke_gae_project_id
helm install -f ../deployment/helm/static-values/remote-config.yaml remote-config ../deployment/helm/deployment --set instance.replicas=3 --set docker.image=$DISTRIBUTOR_TAG
helm install -f ../deployment/helm/static-values/proxy.yaml proxy ../deployment/helm/deployment --set instance.replicas=3 --set docker.image=$AS_PROXY_TAG
gcloud container clusters get-credentials $US_GKE_NAME --zone=$TF_VAR_us_location --project=$TF_VAR_us_gke_gae_project_id
kubectl create configmap configs --from-file=table=../docs/configurationTable.example.json
kubectl create secret generic jwts --from-file=public=../docs/public.example.pem
helm install persistent-storage ../deployment/helm/persistent-storage
helm install -f ../deployment/helm/static-values/updater.yaml updater ../deployment/helm/deployment --set instance.replicas=1 --set docker.image=$UPDATER_TAG
helm install -f ../deployment/helm/static-values/updater.yaml individual-service ../deployment/helm/individual-service --set static.name=$UPDATER_IP_NAME
helm install -f ../deployment/helm/static-values/remote-config.yaml remote-config ../deployment/helm/deployment --set instance.replicas=3 --set docker.image=$DISTRIBUTOR_TAG
helm install -f ../deployment/helm/static-values/proxy.yaml proxy ../deployment/helm/deployment --set instance.replicas=3 --set docker.image=$US_PROXY_TAG
helm install -f ../deployment/helm/static-values/remote-config.yaml remote-config-mci-service ../deployment/helm/mci-service --set static.ip=$DISTRIBUTOR_IP --set certificate.name=$DISTRIBUTOR_CERT_NAME
helm install -f ../deployment/helm/static-values/proxy.yaml proxy-mci-service ../deployment/helm/mci-service --set static.ip=$GEO_DISTRIBUTOR_IP --set certificate.name=$GEO_DISTRIBUTOR_CERT_NAME
cp ../deployment/gae/app.template.yaml ../app.yaml
sed -i -e "s/{{.DistributorUrl}}/$TF_VAR_distributor_domain/g" ../app.yaml
gcloud config set project $TF_VAR_eu_gae_project_id
gcloud app deploy .. --quiet
gcloud config set project $TF_VAR_as_gae_project_id
gcloud app deploy .. --quiet
gcloud config set project $TF_VAR_us_gke_gae_project_id
gcloud app deploy .. --quiet
