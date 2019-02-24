
# NYC Cab Data

This code is 100% in compliance with golint, go_vet and gofmt. Check this for more details: [![Go Report Card](https://goreportcard.com/badge/github.com/nikhil-github/api-cab-data)](https://goreportcard.com/report/github.com/nikhil-github/api-cab-data) [![Build](https://travis-ci.org/nikhil-github/api-cab-data.svg?branch=master)](https://travis-ci.org/nikhil-github/api-cab-data)


## Introduction:

API Cab Data provides http endpoint to query how many trips a particular cab made for a pick up date(using date part in pick up datetime). Endpoint accepts one or many medallions(comma separated) and returns trips made by each medallion.
Results are cached for faster access next time. Endpoint allows user to bypass cache for results.

`/trips/v1/medallion/:medallions?pickupdate=:pickupdate&bypasscache=:bypasscache` - GET

API provides a second endpoint to clear the cache entries

`/trips/v1/cache/contents` - DELETE

API Health Check

`/health` - GET

## Project Set up and Structure:

Application is designed with a simple three layered architecture.

Controller -> Service -> Data Access Objects

GO version 1.9 is used for building the app with MYSQL server the database providing cab trip data.

### Structure
Projects have been packaged based on their responsibility(SINGLE RESPONSIBILITY principle)
- handler -> Responsible for handling http requests.
- wiring -> Initialisation and Wiring up the components.
- database -> Data access layer for DB operations.
- service -> Business logic layer that interacts with database and cache
- cache -> Provides interface to Get / Set / Clear cache entries
- output -> Defines the output JSON structure

### External Packages
- github.com/gorilla/mux (http request routing and dispatching)
- github.com/muesli/cache2go (Concurrency-safe golang caching library with expiration capabilities)
- go.uber.org/zap (provides fast, structured, leveled logging)
- github.com/pkg/errors(error handling)
- github.com/stretchr/testify (unit test suite, mocking and assertion)
- gopkg.in/DATA-DOG/go-sqlmock.v1(SQL mocking library)
- github.com/jmoiron/sqlx (lib with set of extensions on go's standard database/sql library)

Dep is the dependency management tool.

### Unit Tests
- Follows data table approach
- Consistent pattern using Args/Fields/Want format
- Assertion/Mocking using testify

### Config values
- Supplied through .env to run locally
- Docker env file .env.docker

### Pre-Requisites:
- Git (just to clone the repo)
- Docker and Docker-compose

## Installation:
 Clone this repository
`https://github.com/nikhil-github/api-cab-data.git`

Note : Please run the database migration before consuming the API.

### Run Locally

`make run`

### Run in Docker

`make run-docker`

API will be listening on port 3000 , endpoint:

`http://localhost:3000/trips/v1/medallion/:medallions?pickupdate=:pickupdate&bypasscache=:bypasscache`


### Make targets

1. `make` - build the project
2. `make fmt` - format the codebase using `go fmt` and `goimports`
3. `make test` - run unit tests for the project

### Database Migration

1. `make start-db` - start a mysql database container
2. `docker exec -it mysql bash` - log in to mysql container
3. `mysql -u root -p cabtrips < ny_cab_data_cab_trip_data_full.sql` - load SQL

password : password

Import SQL takes a little while (~30 minutes) due to the large size of the SQL

Migration is required just one time unless DB volumes are removed.

### Instructions to test the API
Number of trips made by cab with medallion - 67EB082BFFE72095EAF18488BEA96050 on 31st Dec 2013

 1. Request with one medallion
 `curl http://localhost:3000/trips/v1/medallion/67EB082BFFE72095EAF18488BEA96050?pickupdate=2013-12-31&bypasscache=true`

   ```
   [
     {
        "medallion": "67EB082BFFE72095EAF18488BEA96050",
        "trips": 39
     }
   ]
   ```
  2. Request with multiple medallions

  `curl http://localhost:3000/trips/v1/medallion/67EB082BFFE72095EAF18488BEA96050,D7D598CD99978BD012A87A76A7C891B7?pickupdate=2013-12-31&bypasscache=true`

  ```
  [
    {
        "medallion": "67EB082BFFE72095EAF18488BEA96050",
        "trips": 39
    },
    {
        "medallion": "D7D598CD99978BD012A87A76A7C891B7",
        "trips": 0
    }
  ]
  ```

  3. Clear Cache

  `curl -X DELETE http://localhost:3000/trips/v1/cache/contents`

  ```
  Cache Cleared
  ```
### Assumptions:
- No requirement for distributed cache and in memory caching is allowed.
- The date format used in this API is YYYY-MM-DD.
- By Passing cache is an optional parameter and if no supplied its value is false.
- Endpoint allows maximum of 20 medallions per request.
- Secrets/Configs are supplied as env variables.
