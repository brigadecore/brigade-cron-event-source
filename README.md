# Brigade Cron Event Source

<img width="100" align="left" src="logo.png">

This Brigade V2 Cron Event Source creates one Kubernetes cron job for each "cronEvents:" parameter list on the Event Source Helm installation chart values file. 

Each cron job will trigger a Brigade v2 event, based on cronjob schedule, Source, Type, Qualifiers, Labels and Payload provided on the values file.

<br clear="left"/>

## Why a Cron Event Source?

Before getting started with this event source, do ask yourself if you truly need
it. If your only interest is in executing a _simple_ task on a particular
schedule, Kubernetes, by itself, gives you everything you need. In such a case,
we don't recommend over-complicating things by involving Brigade!

If, however, you wish to execute a more complex workflow on a particular
schedule or if your use case is well-served by other Brigade features, _do_
consider utilizing this event source.

A non-exhaustive set of reasons you might wish to use this gateway includes:

* Your workflow is complex and involves multiple containers that need to execute
  concurrently or in serial.
* You also have _other_ events you need to handle.
* You care about capturing and saving logs from your automated task.

## Pre-requisites (Following https://quickstart.brigade.sh)
Helm 3.7+ installed (Installation and Development)

Brig CLI installed (Installation and Development)

Brigade V2 installed (Installation and development(recommended Kind cluster)

Docker Desktop Installed (Development)

Kind Installed (Development)

## Installation

### Create brigade Service Account (make sure you keep generated api token in a safe place). Also used in development.
```console
brig service-account create --id cron --description cron
```
### Give permission to Service Account to create events. Also used in development.
```console
brig role grant EVENT_CREATOR --service-account cron --source cronproject
```
### Create Secret with API Server Token
```console
kubectl create secret generic brigade-api-server-token --namespace cron-event-source\
--from-literal=apitoken=e01f2b82a1d042889396889ad741e9f2E.......
```
### Helm Install
```Console
helm install cron-event-source --version 0.1.0 --create-namespace --namespace cron-event-source --wait --timeout 300 \
./charts/brigade-cron-event-source/ -f ./charts/brigade-cron-event-source/values.yaml 
```
### Helm Uninstall
```console
helm uninstall cron-event-source
````

## Development Instructions
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
### Build Docker Image locally and push to local Kind Cluster
```console
export DOCKER_REGISTRY=jorgearteiro
export VERSION=0.1.0
make hack-build-no-cache
make hack-load-image
```

Image published on docker hub: jorgearteiro/brigade-cron-event-source:0.1.0
### Testing on Local Kind cluster
```console
kubectl run cron-event-source --image=brigade-cron-event-source:edge --restart=Never --namespace cron-event-source \
--env API_ADDRESS="https://brigade-apiserver.brigade.svc.cluster.local" \
--env API_TOKEN="[<Enter you API token here>]" \
--env API_IGNORE_CERT_WARNINGS=true \
--env BRIGADE_SOURCE=cronsource \
--env BRIGADE_TYPE=cron
```

### Local VS Code debug
Create an .env file. This file is included on .gitignore. Make sure port-forward is connected with Kind Cluster
```console
API_ADDRESS="https://localhost:8444"
API_TOKEN="[<Enter you API token here>]"
API_IGNORE_CERT_WARNINGS=true
BRIGADE_SOURCE=cronsource
BRIGADE_TYPE=cron
```
Create .vscode/launch.json file
```json
{
    // Use IntelliSense to learn about possible attributes.
    // Hover to view descriptions of existing attributes.
    // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
    "version": "0.2.0",
    "configurations": [
        {
            "name": "main",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${fileDirname}",
            "envFile": "${workspaceFolder}/.env"
        }
    ]
}
```