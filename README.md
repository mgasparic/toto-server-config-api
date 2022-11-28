# Toto Server Config Api

The exact description of the task is provided
in https://github.com/mgasparic/toto-server-config-api/blob/main/docs/task.pdf. The task is challenging because it
requires IP geolocation, while aiming at response times under 200ms when serving up to 50 requests per second.

It is obvious that we need a solution that is distributed around the globe. As a side note, the light can travel only
60.000km in 200ms, and the requests have to go back and forth - not all with the speed of light and not all in the
straight line. In my experience, using external IP geolocation services, such as https://ipstack.com/, during the
request handling, will delay the execution by more than 500ms. In theory, the client could call such a service before
making a request, and could pass the geolocation value (i.e., in our case, the country code) to the server, but that is
unreliable. I will present the performance metrics of a service that could provide all the data in a single request.
However, the core of the solution, which is presented in this project, is based on the assumption that the geolocation
happens on the server-side.

Multi-regional load balancers on Google Cloud Infrastructure (GCI) can be deployed on top of Google Kubernetes Engine
(GKE) or virtual machines. They allow us to use a single external IP, single domain name, and single SSL certificate,
while routing the requests to the nearest servers. Unfortunately, they do not provide geolocation service. We can use
Google App Engine (GAE), which provides `X-Appengine-Country` header, but, we can only deploy one App Engine solution
per project, and it can only be located in one region.

Out of the box, for GAE, GCI does not provide any service that would behave in a similar way as multi-regional load
balancers. We can use solutions like Cloudflare or custom Nginx set up, however, the requests passing through such
reverse proxies would break the geolocation, since we would receive `X-Appengine-Country` value that belongs to the
server "in the middle", not to the client. Consequently, I decided to implement a solution that consists of 3 layers:
GKE clusters of Nginx servers behind a multi-regional load balancer that perform "301 redirects" to the closest GAE
service (instead of forwarding!), which then forwards the requests to the actual remote-config service (together with
the obtained country code), which is also running on multiple GKE clusters behind a multi-regional load balancer.

My solution doubles the distance that requests have to travel, but ensures that all services being used by a single
request are located in the GCI region that is the closest to the client. Moreover, geolocation is accurate for all the
requests that are not deliberately re-routed by a client (for example, via a virtual private network).

## Architecture

The system is based on micro-service architecture. Three services implemented in Go are:

- Updater, which provides the latest configuration table and can be used to update the table,
- Remote-config Distributor, which subscribes to the updates and provides selected SKU to the client, based on the
  rules,
- Proxy, which performs geolocation.

For the latency testing, I deployed the system to `europe-west2` (London, UK), `asia-east2` (Hong Kong SAR),
and `us-west1` (Oregon, USA) regions. To replicate the deployment, you can
run `./scripts/setup.sh <unique project suffix>`. And to destroy the deployment, you can
run `terraform -chdir=./deployment/terraform/project_permissions destroy -auto-approve`.

The architecture is summarized in the following image:

![](docs/topology.jpg)

For the optimal reliability and availability, Remote-config Distributors and Proxies use 3 replicas in every region.
It is important to note that a higher density of Nginx servers, GAE applications, and Distributors, across other
regions, would lower the latency (e.g., requests from Sao Paolo were pretty slow, it would make sense to add another
instance in S. America). However, there is no need to use a lot of replicas for Updater; I used only one instance. If
Updater service goes down, the latest config tables will be served without interruptions. And it is also important to
note that the Event Listener pattern that I implemented will update the cached tables almost instantly. The update can
be executed with a similar command as this:

```
curl -X POST -H "Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c3IiOiJzeXMuYWRtaW5AemlwcG8uYXBwcy5jb20ifQ.mfv6Md2uWsIyP7ykVv072Tk3MiGnYk2W13wGhtZFVHMGlytlgFT9w-g8dorAsbf7a2p0432e81MY8PnjaIdpVKyUownzfylZ-2T6TGID8pqbik5lJy0tQUcOtWOfsUbGlk1SzSXGa1oZaq9VTguafET7ofpqY3tFP_BSKtzeqvHGi1--oN5gta8VpHxDIq2Ebuwipp2MWXjcxX5T_stA3zpEDgflKHWh38XvJUSC03HBb_8ImCASiIvnAqLjRj86j4avLGaMLOvMs08Q0txwqiXR0rBio09P7o-NRCclK7u0Y71HbznWtJzu3O1VdEp6k-Glr2l59kZWyAuyPK-crg" \
http://updater.sharedtensors.com/configuration -d '[{"package":"com.softinit.iquitos.mainapp","country_code":"ZZ","percentile_min":0,"percentile_max":100,"main_sku":"completely_new"}]'
```

## Endpoints

### Updater

| Path           | Method          | Data                              | Success Response    |
|:---------------|:----------------|:----------------------------------|:--------------------|
| /configuration | POST            | Authorization<br>New Config Table | NA                  |
| /subscription  | GET / websocket | Authorization                     | Active Config Table |

### Distributor

| Path                        | Method | Data | Success Response |
|:----------------------------|:-------|:-----|:-----------------|
| /offering/:package/:country | GET    | NA   | SKU              |

### Proxy

| Path               | Method | Data | Success Response |
|:-------------------|:-------|:-----|:-----------------|
| /offering/:package | GET    | NA   | SKU              |

## Latency Tests

The exact results of the latency tests are presented
in https://github.com/mgasparic/toto-server-config-api/tree/main/docs/test-results. The code used to generate the
results and the executable are available in https://github.com/mgasparic/toto-server-config-api/tree/main/test/stress.

I tested the latency from virtual machines hired on Amazon Web Services (AWS). I used `us-west-2` (Oregon, USA),
`eu-central-1` (Frankfurt, Germany), `sa-east-1` (Sao Paulo, Brazil), `ap-southeast-2` (Sydney, Australia),
and `ap-northeast-1` (Tokyo, Japan). I executed the following script, three-times:

```
#!/bin/bash
set -e

export LOAD=''
export MIN='10'
export MAX='100'
export STEP='10'
export METHOD='GET'
export OUTPUT='load-<region>.<attempt nr.>.test.csv'
export URL='https://geo-distributor.sharedtensors.com/offering/com.softinit.iquitos.mainapp'
./executable
```

The averages and proportions of requests returned in less than 200ms are presented in the following table:

| Region \ Latencies | Average | Requests in <=200ms |
|:-------------------|:--------|:--------------------|
| us-west-2          | 102.0ms | 96%                 |
| eu-central-1       | 114.5ms | 97%                 |
| sa-east-1          | 453.3ms | 0%                  |
| ap-southeast-2     | 372.3ms | 0%                  |
| ap-northeast-1     | 159.3ms | 85%                 |

As already mentioned, the upper results require geolocation and consequently require a redirect. If we were able to call
the service with the following script, the time needed to execute the requests would be much shorter.

```
#!/bin/bash
set -e

export LOAD=''
export MIN='10'
export MAX='100'
export STEP='10'
export METHOD='GET'
export OUTPUT='nogeo-load-<region>.<attempt nr.>.test.csv'
export URL='https://distributor.sharedtensors.com/offering/com.softinit.iquitos.mainapp/ZZ'
./executable
```

Results:

| Region \ Latencies | Average | Requests in <=200ms |
|:-------------------|:--------|:--------------------|
| us-west-2          | 40.8ms  | 100%                |
| eu-central-1       | 45.2ms  | 100%                |
| sa-east-1          | 218.4ms | 62%                 |
| ap-southeast-2     | 182.0ms | 85%                 |
| ap-northeast-1     | 75.6ms  | 99%                 |
