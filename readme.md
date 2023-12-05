# [WIP] Spade Tenant Service


## Generating Schema
`DB_ALLOW_MIRGATIONS=true go run . > out.sql`


## Generating Swagger
First, install swag: 
`go install github.com/swaggo/swag/cmd/swag@latest`

Then, run the swag generation script
`./scripts/swag.sh`