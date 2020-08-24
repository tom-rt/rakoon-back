# User service seed

## Description

This is a user handling service, fully written in Go.

The used database engine is postgresql.

The libraries used are exclusively gin-gonic and SQLX.

You can initialize the database model with the files located in the psql folder.

## Auth

It uses json web tokens. I don't use a lib, I implemented everything.

## Testing

It is end to end tested, you can lauch them with a local or distant database:

source .env-test; go test ./tests
