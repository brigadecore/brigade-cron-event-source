# Brigade-cron-gateway
 This is a work-in-progress Brigade 2 Cron job compatible gateway.

## Instructions

### Pre-requisites (Following https://quickstart.brigade.sh)
Docker Desktop Installed

Kind Installed with Cluster created

Brigade V2 intalled


### Create brigade Service Account
brig service-account create --id cron --description cron

### Give permission to Service Account to create events
brig role grant EVENT_CREATOR --service-account cron --source cronproject

### Build Docker Image locally
make hack-build

initial test image published here: jorgearteiro/brigade-cron-gateway:edge it creates event with Source: "cronsource", Type: "cron"

### Create Secret with API Server Token
kubectl create secret generic brigade-api-server-token --namespace default\
--from-literal=apitoken=e01f2b82a1d042889396889ad741e9f2E.......

### Load Docker image built on kind cluster
kind load docker-image brigade-cron-gateway:edge --name <Name of your Kind Cluster>

### Testing on Local Kind cluster
kubectl run crongateway --image=brigade-cron-gateway:edge --restart=Never \
--env API_ADDRESS="https://brigade-apiserver.brigade.svc.cluster.local" \
--env API_TOKEN="e01f2b82a1d0428.........................." \
--env API_IGNORE_CERT_WARNINGS=true

