# Made by tom-rt
# Back end project for the rakoon project
# The libraries used are gin-gonic and GORM
# The database is postgresql

## Auth
# Uses Json web tokens
# A token is valid 15 minutes, then it has to be refreshed
# If you haven't refreshed your token in 7 days, you have to reconnect