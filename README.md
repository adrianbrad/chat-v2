# A Chat system implemented in GO

[![Go Report Card](https://goreportcard.com/badge/github.com/adrianbrad/chat-v2)](https://goreportcard.com/report/github.com/adrianbrad/chat-v2)
[![codecov](https://codecov.io/gh/adrianbrad/chat-v2/branch/master/graph/badge.svg)](https://codecov.io/gh/adrianbrad/chat-v2)

- [A Chat system implemented in GO](#a-chat-system-implemented-in-go)
  - [Installation](#installation)
    - [Clone](#clone)
    - [Go get](#go-get)
  - [Usage](#usage)
    - [Prerequisites](#prerequisites)
    - [Run](#run)
      - [Directly](#directly)
      - [Make](#make)
    - [Exposed API endpoints](#exposed-api-endpoints)
      - [/users](#users)
      - [/rooms](#rooms)
      - [/chat](#chat)
      - [/auth](#auth)
      - [/client/main.wasm](#clientmainwasm)
    - [Backend Authentication](#backend-authentication)

---

## Installation

### Clone

- Clone this repo to your local machine using 

```
git clone https://github.com/adrianbrad/chat-v2
```

### Go get

- Go get this repo to your local machine and into your $GOPATH/src/github.com/adrianbrad folder using 

```
go get https://github.com/adrianbrad/chat-v2
```

---

## Usage

### Prerequisites

- A PostgreSQL database and valid credentials to connect passed to a database config file in the /configs folder.

- Set the `BASEDIR` variable in the Makefile to the project root and the `DATABASE_CONFIG_FILE` variable to the config file containing your valid db credentials.

### Run

#### Directly

```
go run ./cmd/chat-database -b={project-path} -d={database-config-file} -a={application-config-file}
```

#### Make

```
make run
```

### Exposed API endpoints

#### /users

- Allowed methods: `GET, POST, PUT, DELETE`
- Requiers [Backend Authentication](#backend-authentication)

#### /rooms

- Requiers [Backend Authentication](#backend-authentication)

#### /chat

#### /auth

#### /client/main.wasm

### Backend Authentication

- When making calls to backend endpoints that require [Backend Authentication](#backend-authentication), you have to add an `Authorization` header containg a HMAC-SHA256 hash, created from the current time in epoch and the secret key. The current time in epoch should be added to the `Date` header.