# go-htmx-01

Things to look at:

- Sass -> postcss without npm


## What Are We Making?

A fancy Todo app

- User login w/ passwordless login
  - Table + SQL setup
  - Login form
  - Login post
    - CSRF validateion
  - Parse login stuff
  - validate
  - return error
  - on valid, generate passcode w/ timeout + send email
  - reply with success + check your email message to user
  - add email link route
  - on emmail link click, check for existence of key attached to email + signing and things
  - if success, set cookie + redirect
  - if error, error
  - add a dashboard oproected by session check to redirect to
  - verify the user has access when hitting the dash
- CSS (tailwind or sth else?)
- Form entry + validation
- DB migrations, seeding + sqlc
- Websocket w/ real time updates
https://tailwindcss.com/blog/automatic-class-sorting-with-prettier

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