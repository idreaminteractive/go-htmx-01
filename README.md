# go-htmx-01


## What Are We Making?

A fancy Todo app

## Todo

- Update db schema to be notes instead of todos
- Creating Notes from Dashboard
- Setting Note as public or private
- Editing notes
- Deleting notes
- Public view on root route
- SSE with some neat stuff



## What Stack and Services are we using?

A very opinionated stack we can reuse over and over.

## Stack

### Confirmed:

- Go
- TailwindCSS
- Templ
- htmx
- Air 
- SQLC
- github.com/caarlos0/env/v9
- Echo Webserver
- Goose
- SQLite + LiteFS 

### Testing:

- ...

## Services:

### Confirmed:
- Fly.io
- Doppler
- Gitpod


### Testing:

- Honeycomb.io
- Bunny.net CDN

# Creating Migrations

`goose -dir migrations sqlite3 /litefs/potato.db create add_user_email sql`
