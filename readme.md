<div align=center>
<img src="docs/logo.png" width=300 />

# Spade Tenant Services
</div>


Tenant APIs and services for the [Spade](https://github.com/data-preservation-programs/spade) - Filecoin network storage proposal management service

The Tenant services handle all aspects of Spade Storage Proposals from the Tenant side. At a high level, this includes interfaces to manage the following:
- Policies
- Storage Contracts
- Collections
- Content (i.e, sources of Filecoin Pieces)
- Replication Constraints
- Storage Provider SLA Compliance

A high level architecture diagram of the Tenant functionalities and their relationships can be found [here](https://miro.com/app/board/uXjVMwfo8TM=/?share_link_id=267050417898) 

<img src="docs/arch-diagram.png" width=500 />

# API
All functionality is exposed over a REST API. Live API docs can be viewed at https://spade.storage (Coming soon), or inside the repository [docs folder.](https://github.com/data-preservation-programs/spade-tenant/tree/main/docs)

Users are encouraged to deploy their own instance of the Spade Tenant Services, and customize it for their specific applications, however we will make a public, hosted instance available at https://apiv1.spade.storage/ (Coming soon).

# How to run Spade Tenant Services

## 0. Prerequisites
You must have a Postgres database available to connect to. The easiest way to do that is with a Docker container, such as the following 

```bash
docker run --name spade-tenant-postgres -p 5432:5432 -e POSTGRES_PASSWORD=password -d postgres:14.7
```

## 1. Configuration
The Tenant Services are configured using environment variables. The following variables may be specified:

```bash
PORT=3000 # TCP port that API runs on; optional; default 3000 
DB_URL=postgres://postgres:password@localhost:5432/ # URI to connect to Postgres DB; required
DRY_RUN=false # If true, no changes will be made to the DB, SQL will simply be logged ; optional; default false
DEBUG=false # If true, debug logs will be printed; optional; default false
```

## 2. Database Migration
The database schema is managed using Gorm and Gormigrate, however migrations are only applied if explicitly requested.

To apply migrations to the database run the following command:

`DB_ALLOW_MIRGATIONS=true ./spade-tenant-svc`

To simply output generated SQL without applying to the database, you can execute this instead
`DRY_RUN=true DB_ALLOW_MIRGATIONS=true go run ./spade-tenant-svc out.sql`

## 3. Run

The Tenant Services are distributed as a single binary. No runtime flags are necesssary. To run the binary, simply execute it from the command line

```bash
./spade-tenant-svc
```


# Developer Notes 

## Prerequisites
- Go v1.20+

## Building
To build the Tenant Services, run the following command:

`make clean all`

This will output the compiled binary at `./spade-tenant-svc`

You can install it to `/usr/local/bin` by executing `make install`

For convenience, if an `.env` file exists at the root of the repository, its values will be loaded into the environment.

## Testing
(Coming soon) Execute tests with `go test .`


## Generating Swagger Docs
First, install swag: 
`go install github.com/swaggo/swag/cmd/swag@latest`

Then, run the swag generation script
`./scripts/swag.sh`