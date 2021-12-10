<img width="100" align="left" src="logo.png">

# Brigade-cron-gateway
 This is a work-in-progress Brigade 2 Cron job compatible gateway.

## Instructions

### Pre-requisites (Following https://quickstart.brigade.sh)
Docker Desktop Installed

Helm 3.7+ installed

Kind Installed

Brig CLI installed 

Brigade V2 intalled on Kubernetes Kind Cluster using Helm Chart

### Create Kind cluster called brigade
```console
kind create cluster --name brigade
```
### Open port-forward with api server pod running on the kind cluster
```console
kubectl --namespace brigade port-forward service/brigade-apiserver 8443:443 &>/dev/null &
```
### Login to Brigade API Server from brig CLI (ps: Using Root password is not recommended for Prodcution)
```console
export APISERVER_ROOT_PASSWORD=$(kubectl get secret --namespace brigade brigade-apiserver --output jsonpath='{.data.root-user-password}' | base64 --decode)
```
```console
brig login --insecure --server https://localhost:8443 --root --password "${APISERVER_ROOT_PASSWORD}"
```

### Create brigade Service Account (make sure you keep generated api token in a safe place)
```console
brig service-account create --id cron --description cron
```
### Give permission to Service Account to create events
```console
brig role grant EVENT_CREATOR --service-account cron --source cronproject
```
### Build Docker Image locally
```console
make hack-build
```

initial test image published here: jorgearteiro/brigade-cron-gateway:edge it creates event with Source: "cronsource", Type: "cron"

### Create Secret with API Server Token
```console
kubectl create secret generic brigade-api-server-token --namespace default\
--from-literal=apitoken=e01f2b82a1d042889396889ad741e9f2E.......
```
### Load Docker image built on kind cluster
```console
kind load docker-image brigade-cron-gateway:edge --name brigade
```
### Testing on Local Kind cluster
```console
kubectl run crongateway --image=brigade-cron-gateway:edge --restart=Never \
--env API_ADDRESS="https://brigade-apiserver.brigade.svc.cluster.local" \
--env API_TOKEN="e01f2b82a1d0428.........................." \
--env API_IGNORE_CERT_WARNINGS=true
```
