
# API NYC Cab Data


## Introduction:

API Cab Data provides http endpoint to query how many trips a particular cab made for a pick up date(using date part in pick up datetime. Endpoint accepts one or many medallions(comma separated) and returns trips made by each medallion.
Results are cached for faster access next time. Endpoint allows use to bypass cache for results.

`/trips/medallion/:medallions?pickupdate=:pickupdate&bypasscache=:bypasscache` - GET

API provides a second endpoint to clear the cache entires

`/trips/cache/contents` - DELETE

API Health Check

`/health`

## Project Set up and Structure:

Project follows a three layered architecture where the handlers receives the requests and
delegate it to service layer. Service layer interacts with DB and cache APIs to fetch
the required data.

# Dep package tool have been used for dependency management.

# Structure
Projects have been packaged based on their responsibility(SINGLE RESPONSIBILITY principle)
- handler -> Responsible for handling http requests.
- wiring -> Initialisation and Wiring up the components.
- database -> Data access layer for DB operations.
- service -> Business logic layer that interacts with database and cache
- cache -> Provides interface to Get / Set / Clear cache entries
- output -> Defines the output JSON structure

# External Packages

- github.com/gorilla/mux (http request routing and dispatching)
- github.com/muesli/cache2go (Concurrency-safe golang caching library with expiration capabilities)
- go.uber.org/zap (provides fast, structured, leveled logging)
- github.com/pkg/errors(error handling)
- github.com/stretchr/testify (unit test suite, mocking and assertion)
- gopkg.in/DATA-DOG/go-sqlmock.v1(SQL mocking library)
- github.com/jmoiron/sqlx (supporting named queries in SQL)



## Assumptions:
- No requirement for distributed cache and in memory caching is allowed.
- The date format used in this API is YYYY-MM-DD.
- By Passing cache is an optional parameter and if no supplied its value is false.

## Pre-Requisites:
- Git (just to clone the repo)
- Docker and Docker-compose
`docker --version`
`docker-compose --version`


## Installation:
 Clone this repository
`https://github.com/nikhil-github/api-cab-data.git`

# Run Locally

`make run`

# Run in Docker

`make run-docker`

API will be listening on port 3000 , endpoints:

`http://localhost:3000/trips/medallion/:medallions?pickupdate=:pickupdate&bypasscache=:bypasscache`

Note : Please run the database migration before consuming the API.

### Make targets

`make` - build the project

`make fmt` - format the codebase using `go fmt` and `goimports`

`make test` - run unit tests for the project

### Important: This will download and launch two docker images, please be patient.


### Database Migration

1. `make start db` - start a mysql database container
2. `docker exec -it mysql bash` - log in to mysql container
3. `mysql -u root -p cabtrips < ny_cab_data_cab_trip_data_full.sql`

password : password

Import SQL takes a little while (~30 minutes) due to the large size of the SQL

### Tests

Number of trips made by cab with medallion - 67EB082BFFE72095EAF18488BEA96050 on 31st Dec 2013

- http://localhost:3000/trips/medallion/67EB082BFFE72095EAF18488BEA96050?pickupdate=2013-12-31&bypasscache=true

   ```
   [
     {
        "medallion": "67EB082BFFE72095EAF18488BEA96050",
        "trips": 39
     }
   ]
   ```
