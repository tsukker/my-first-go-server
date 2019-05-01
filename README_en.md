:us: [:jp:](./README.md)

# My First Go Server

## Go Simple Server

`simple_server/` directory

### Usage

    cd simple_server
    docker build .
    docker run -p 8080:8080 [DOCKER_IMAGE]

## Service Publish

(IP address):8080 (available in the limited period)

### GCP Cloud Registry

- Upload the Docker image to GCP Cloud Registry by using Cloud Build
- [GCP reference](https://cloud.google.com/cloud-build/docs/quickstart-docker)

Build using Dockerfile

    gcloud builds submit --tag gcr.io/[PROJECT_ID]/[DOCKER_IMAGE] .

### Deploy

1. Kubernetes Engine > Clusters tab: Click "Deploy container" and select the Docker image uploaded above.
2. Workloads tab: In "Publish service" section enable "Load balancer" and set port number 8080
