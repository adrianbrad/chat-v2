# A Chat system implemented in GO

[![Go Report Card](https://goreportcard.com/badge/github.com/adrianbrad/chat-v2)](https://codecov.io/gh/adrianbrad/chat-v2)
[![codecov](https://codecov.io/gh/adrianbrad/chat-v2/branch/master/graph/badge.svg)](https://codecov.io/gh/adrianbrad/chat-v2)

- [A Chat system implemented in GO](#a-chat-system-implemented-in-go)
  - [Installation](#installation)
    - [Clone](#clone)
    - [Go get](#go-get)
  - [Usage](#usage)
    - [Prerequisites](#prerequisites)
    - [Run](#run)

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

```
make run
```