:us: [:jp:](./README.md)

# My First Go Server

## Go Simple Server

`simple_server/` directory

### Usage

    cd simple_server
    docker build .
    docker run -p 8080:8080 [DOCKER_IMAGE]

### Behavior

    curl -XGET -H 'Content-Type:application/json' http://localhost:8080/

response:

    {"message":"Hello World!!"}

## Exposing service

http://35.247.11.48:8080/ (available only in some limited period)

### GCP Cloud Registry

- Upload the Docker image to GCP Cloud Registry by using Cloud Build
- [GCP reference](https://cloud.google.com/cloud-build/docs/quickstart-docker)

Build using Dockerfile

    gcloud builds submit --tag gcr.io/[PROJECT_ID]/[DOCKER_IMAGE] .

### Deploy

1. Kubernetes Engine > Clusters tab: Click "Deploy container" and select the Docker image uploaded above.
2. Workloads tab: In "Exposing services" section enable "Load balancer" and set port number 8080

## Go RESTful Server

`restful_server/` directory

### Usage

    cd restful_server
    docker-compose build
    docker-compose up

### Request and Response samples

#### GET /users

    curl -XGET -H 'Content-Type:application/json' http://localhost:8080/users

response:

    [{"id":1,
    "name":"test",
    "email":"hoge@example.com",
    "created_at":"2019-05-05T01:42:23.4185993+09:00",
    "updated_at":"2019-05-05T01:42:23.4185993+09:00"},
    {"id":2,
    "name":"test2",
    "email":"fuga@example.com",
    "created_at":"2019-05-05T02:23:01.3296964+09:00",
    "updated_at":"2019-05-05T02:23:01.3296964+09:00"}]


#### GET /users/:id

    curl -XGET -H 'Content-Type:application/json' http://localhost:8080/users/1

response:

    {"id":1,
    "name":"test",
    "email":"hoge@example.com",
    "created_at":"2019-05-05T01:42:23.4185993+09:00",
    "updated_at":"2019-05-05T01:42:23.4185993+09:00"}

#### POST /users

    curl -XPOST -H 'Content-Type:application/json' http://localhost:8080/users -d '{ "name": "test3", "email": "foo@example.com" }'

response:

    {"id":3,
    "name":"test3",
    "email":"foo@example.com",
    "created_at":"2019-05-05T02:26:51.1402154+09:00",
    "updated_at":"2019-05-05T02:26:51.1402154+09:00"}


#### PUT /users/:id

    curl -XPUT -H 'Content-Type:application/json' http://localhost:8080/users/3 -d '{ "name": "test3", "email": "bar@example.com" }'

response:

    {"id":3,
    "name":"test3",
    "email":"bar@example.com",
    "created_at":"2019-05-05T02:26:51.1402154+09:00",
    "updated_at":"2019-05-05T02:30:20.1519588+09:00"}

#### DELETE /users/:id

    curl -XDELETE -H 'Content-Type:application/json' http://localhost:8080/users/1

no response, status code: 204
