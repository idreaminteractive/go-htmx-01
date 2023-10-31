# go-htmx-01


## What Are We Making?

A fancy Todo app

## Todo

- Form validations in a sane way
    - look @ https://dev.to/thanhphuchuynh/customizing-error-messages-in-struct-validation-using-tags-in-go-4k0j 
    - create a simple way to add in a field -> string + then can pass in errors per field
- flash sessions
- SSE with some neat stuff
- Semantic layout, styling, UX and better test coverage.
- GH Action setup for CICD



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



# Structure

What is the best approach for our structure of the application?

- db + main + http is all good at the moment
- services are good (but move to an Init to return an interface for it) (maybe look into moq or mockery during tests)
- Routes:
    - Have a main route that loads the entire thing, composed of it's sub components
    - The sub components can be auto-refreshed as necessary w/ htmx events 
        - means controllers will be able to get data for all the sub parts separately
        - similar to a remix setup, actually. 
        - If this is the case - how do we efficiently pull up data?
        - A sub route needs it's own data fetching - unless it's provided to it?
    - if we detect that it's from htmx - we can send partials 
    - The route is controller code 
        - Session checks, CSRF checks, etc 
        - Can handle sub posts if we wanted too + return items or oob updates or triggers
        - Composes the responses w/ templates
    - Templates are not colocated with controllers, due to templ generation
- don't need to do pointers all the time?


    some interesting articles on the approach:
    https://medium.com/@kyodo-tech/lindy-approach-to-web-development-htmx-and-go-809bdfdf2279

    https://www.youtube.com/watch?v=F9H6vYelYyU