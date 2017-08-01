# Audition

[![Travis](https://img.shields.io/travis/arbourd/audition.svg)](https://travis-ci.org/arbourd/audition) 

A demonstration with Golang, Hyperapp and Terraform

- [Development](#development)
    - [Getting Started](#getting-started)
      - [Build locally](#build-locally)
      - [Build with Docker](#build-with-docker)
      - [Accessing API and Client](#accessing-api-and-client)
    - [Deploy](#deploy)
      - [Variables](#variables)
      - [Provision Instance](#provision-instance)
      - [SSH Access](#ssh-access)
- [Implementation Architecture](#implementation-architecture)
- [Sequence Diagram](#sequence-diagram)
- [REST API Documentation](#rest-api-documentation)
    - [Introduction](#introduction)
      - [Authentication](#authentication)
      - [Responses](#responses)
      - [Errors](#errors)
    - [Messages](#messages)
      - [Retrieve all Messages](#retrieve-all-messages)
      - [Retrieve an existing Message](#retrieve-an-existing-message)
      - [Create a Message](#create-a-message)
      - [Delete a Message](#delete-a-message)

## Development

### Getting Started

#### Build locally

Both Go (1.8 recommended) and Node.js (8.2 recommended) with [yarn](https://yarnpkg.com) must be installed

First, install Node.js dependencies with yarn
```bash
$ cd client && yarn
```

Build the client bundle with webpack
```bash
$ yarn dist
```

Building and run the binary
```bash
$ go build -o app .
$ ./app
> Audition is running...
```

Or just run the app (make sure to exclude `*_test.go`)
```bash
$ go run $(ls *.go | grep -v _test.go)
> Audition is running...
```

#### Build with Docker

With `docker-compose`
```bash
$ docker-compose up --build
```

With `docker build/run`
```bash
$ docker build . -t arbourd/audition
$ docker run --name audition --restart unless-stopped -p 80:8080 -v $(pwd)/db:/db arbourd/audition
```

#### Accessing API and Client

Once running either locally or with Docker, the API at [localhost:8080/api](http://localhost:8080/api) and the client is available at [localhost:8080](http://localhost:8080).

### Deploy

> Deploy to Digital Ocean with Terraform on CoreOS

You will need to:
- Install [Terraform](https://github.com/hashicorp/terraform)
- Get Digital Ocean [access token](https://cloud.digitalocean.com/settings/api/tokens)
- Add (or already have) a [SSH key](https://cloud.digitalocean.com/settings/security) to your Digital Ocean security settings
  - Get the SSH key ID(s) from [Digital Ocean API](https://developers.digitalocean.com/documentation/v2/#list-all-keys) or with [doctl](https://github.com/digitalocean/doctl): `doctl compute ssh-key-list`
  - The path of a private key from a corresponding ID above

#### Variables

| Variable           | Description                                       | Type   |
| ------------------ | ------------------------------------------------- | ------ |
| TF_VAR_do_token    | Digital Ocean access token                        | string |
| TF_VAR_ssh_key_ids | ID(s) of SSH keys on Digital Ocean                | list   | 
| TF_VAR_private_key | Path of private key (defaults to `~/.ssh/id_rsa`) | string |

You can make additional changes to the Terraform configuration by simply editing `deploy/main.tf`

#### Provision Instance

From inside the `deploy` folder, initialize Terraform
```bash
$ cd deploy && terraform init
```

Create a compute instance and deploy app
```bash
$ TF_VAR_do_token="accesstoken" TF_VAR_ssh_key_ids="[1234, 4567]" TF_VAR_private_key="~/.ssh/deploy"\
  terraform apply
```

Destroy compute instance
```bash
$ TF_VAR_do_token="accesstoken" TF_VAR_ssh_key_ids="[1234, 4567]" TF_VAR_private_key="~/.ssh/deploy" \
  terraform destroy
```

These variables can also be exported
```bash
$ export TF_VAR_do_token="accesstoken"
$ export TF_VAR_ssh_key_ids="[1234, 4567]"
$ export TF_VAR_private_key="~/.ssh/deploy"
$ terraform apply
```

#### SSH Access

To access the box after the fact, simply SSH into using the public IP and the `core` user
```bash
$ ssh core@audition.ip
```

Starting service
```bash
$ sudo systemctl start audition.service
```

Restarting service
```bash
$ sudo systemctl restart audition.service
```

Starting the container manually without systemd service
```bash
$ docker run -d --name audition --restart unless-stopped -p 80:8080 -v $(pwd)/db:/db arbourd/audition
```

## Implementation Architecture

```
                               ┌────────────────────────────────────────────────────────────────────┐
                               │                                                                    │
                      Serves   │                       Server          ┌────────────────────────┐   │
         ┌────────────GET /────┼─────────────────┐                     │                        │   │
         │                     │                 │                     │   api.go               │   │
         ▼                     │                 │                     │                        │   │
┌─────────────────┐            │                 │                     │   APIService           │   │
│                 │            │                 │                     │   <struct>             │   │
│     Client      │            │                 │                     │   - db Database        │   │
│                 │            │                 │                     │                        │   │
│   index.html    │            │                 │                ┌───▶│   ┌─────────────────┐  │   │
│   bundle.js     │            │     ┌───────────────────────┐    │    │   │     Message     │  │   │
│                 │            │     │                       │    │    │   │ - list all      │  │   │
│        ▲        │            │     │  main.go              │    │    │   │ - get one       │  │   │
│        │        │────────┐   │     │                       │    │    │   │ - create one    │  │   │
│     Webpack     │        │   │     │  - Inits DB           │    │    │   │ - destroy one   │  │   │
│        │        │        │   │     │  - Registers API      │Register │   └─────────────────┘  │   │
│        │        │        │   │     │    handlers           │Handlers └────────────────────────┘   │
│                 │  Calls │   │     │  - Serves requests    │    │                 ▲               │
│   index.html    │   /api └───┼────▶│    w/ mux             │────┘                 │               │
│   index.js      │            │     │                       │                      │               │
│                 │            │     │                       │                      │               │
└─────────────────┘            │     │                       │         ┌────────────────────────┐   │
                               │     │                       │         │                        │   │
                               │     │                       │         │   db.go                │   │
                               │     │                       │         │                        │   │
                               │     └───────────────────────┘         │   Database             │   │
                               │                 │                     │   <interface>          │   │
                               │                 │                     │   - Accesses bolt.DB   │   │
                               │                 │                     │   - CRUD opts for msg  │   │
                               │        .─────.  └─Init DB────────────▶│   ┌─────────────────┐  │◀─┐│
                               │      ,'       `.                      │   │     Message     │  │  ││
                               │     ;  Bolt DB  :                     │   │ - id            │  │  ││
                               │     :messages.db;                     │   │ - message       │  │  ││
                               │      ╲         ╱                      │   │ - isPalindrome  │  │  ││
                               │       `.     ,'                       │   │ - createdAt     │  │  ││
                               │         `───'                         │   └─────────────────┘  │  ││
                               │           ▲                           └────────────────────────┘  ││
                               └───────────┼───────────────────────────────────────────────────────┼┘
                                           └───────────────────JSON serialization──────────────────┘ 
```                                                                                                     


## Sequence Diagram

```
                    ┌─────────────────────┐              ┌─────────────────────┐                ┌─────────────────────┐
                    │                     │              │                     │                │                     │
                    │       Client        │              │       Handler       │                │      Database       │
                    │                     │              │                     │                │                     │
                    └─────────────────────┘              └─────────────────────┘                └─────────────────────┘
                               │                                    │                                      │           
                                                                  ┌───┐                                                
                             ┌─┴─┐         HTTP Request           │   │                                    │           
──────Create a message─────▶ │   │ ─────POST /api/messages────▶   │   │                                                
                             │   │                                │   │                                    │           
                             │   │                                │   │                                                
                             │   │                                │   │                                  ┌─┴─┐         
                             │   │                                │   │   ──────Store in Database────▶   │   │         
                             │   │                                │   │                                  │   │         
                             │   │                                │   │                                  │   │         
                             │   │          HTTP Response         │   │               Message            │   │         
                             │   │  ◀────────201 CREATED────────  │   │    ◀────Serialized as JSON─────  │   │         
                             └─┬─┘                                └─┬─┘                                  └─┬─┘         
                                                                                                                       
                               │                                    │                                      │                
```
## REST API Documentation

### Introduction

#### Authentication

This API has no authentication or authorization

#### Responses

Successfull responses return JSON either as an object or an array
```json
[
  {
    "name": "example",
    "value": 0
  }
]
```

#### Errors

Errors are returned as JSON with an error attribute and a message. Example:
```json
{
  "error": "Not Found",
  "message": "Could not find an object with ID: 0"
}
```

### Messages

Attributes of a message

| Name         | Type    | Description                                         |
| ------------ |---------| --------------------------------------------------- |
| id           | integer | A unique ID used to identify the message            |
| message      | string  | The value of the message itself                     |
| isPalindrome | boolean | Whether or not the message is a palindrome          |
| createdAt    | string  | A time value of when the message object was created |

#### Retrieve all Messages

> `GET /api/messages`

cURL example
```bash
$ curl -X GET -H "Content-Type: application/json" "http://audition.arbr.ca/api/messages" 
```

Response headers
```
Content-Type: application/json
Status: 200
```

Response body
```json
[
  {
    "id": 1,
    "message": "hello",
    "isPalindrome": false,
    "createdAt": "2017-07-30T23:42:20Z"
  },
  {
    "id": 2,
    "message": "anna",
    "isPalindrome": true,
    "createdAt": "2017-07-30T23:42:24Z"
  }
]
```

#### Retrieve an existing Message

> `GET /api/messages/id`

cURL example
```bash
$ curl -X GET -H "Content-Type: application/json" "http://audition.arbr.ca/api/messages/1" 
```

Response headers
```
Content-Type: application/json
Status: 200
```

Response body
```json
{
  "id": 2,
  "message": "anna",
  "isPalindrome": true,
  "createdAt": "2017-07-30T23:42:24Z"
}
```

#### Create a Message

> `POST /api/messages`

Request params
```json
{
  "message": "The quick brown fox!"
}
```

cURL example
```bash
$ curl -X POST -H "Content-Type: application/json" -d '{"message":"The quick brown fox!"}' \
  "http://audition.arbr.ca/api/messages" 
```

Response headers
```
Content-Type: application/json
Status: 201
```

Response body
```json
{
  "id": 3,
  "message": "third",
  "isPalindrome": false,
  "createdAt": "2017-07-30T23:46:20Z"
}
```

#### Delete a Message

> `DELETE /api/messages/id`

cURL example
```bash
$ curl -X DELETE -H "Content-Type: application/json" "http://audition.arbr.ca/api/messages/1" 
```

Response headers
```
Content-Type: application/json
Status: 204
```
