# rakoon back end

## Description

The API for the rakoon project, fully written in Go

The used database engine is postgresql

The libraries used are gin-gonic and GORM

Made by tom-rt

## Auth
Uses Json web tokens

A token is valid 15 minutes, then it has to be refreshed

If you haven't refreshed your token in 7 days, you have to reconnect
