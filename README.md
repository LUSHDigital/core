# Lush Digital - Micro Service Core (Golang) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://raw.githubusercontent.com/LUSHDigital/microservice-core-golang/master/LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/LUSHDigital/microservice-core-golang)](https://goreportcard.com/report/github.com/LUSHDigital/microservice-core-golang) [![Build Status](https://travis-ci.org/LUSHDigital/microservice-core-golang.svg?branch=master)](https://travis-ci.org/LUSHDigital/microservice-core-golang)
A set of core functionality and convenience structs for a Golang microservice

## Description
This package is intended to provide a quick and easy bootstrap of functionality that a micro service is expected
to provide. This includes an information route that could be used by a service registry, it also includes a health
check route to verify your microservice is working.

The package also contains some convenience classes to help develop microservices.

## Package Contents
* Route struct for use with HTTP routing
* Response struct to provide a standardised response format for endpoints
* JSON response formatter
* Info struct to provide meta data for your service
* Helper function to retrieve and ensure environment variables.

## Installation
Install the package as normal:

```bash
$ go get -u github.com/LUSHDigital/microservice-core-golang
```

## Documentation
* [General](https://godoc.org/github.com/LUSHDigital/microservice-core-golang)
* [Response](https://godoc.org/github.com/LUSHDigital/microservice-core-golang/response)
* [Format](https://godoc.org/github.com/LUSHDigital/microservice-core-golang/format)
* [Routing](https://godoc.org/github.com/LUSHDigital/microservice-core-golang/routing)
