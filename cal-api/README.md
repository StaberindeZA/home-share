# Cal API

Calendar API

## Architecture

- SQLite DB
- Go HTTP server with REST API

## Database

### Tables

- User
  - id
  - name
  - createdAt
  - updatedAt
- Calendar (add in v2)
  - id
  - entries
- Entry
  - id
  - userId (fk)
  - start (datetime)
  - end (datetime)

## API

### Create Entry

- Path: `POST /v1/entry`
- Body
  - userId (Todo: Unsafe, fetch from auth cookie)
  - start_date_time
  - end_date_time
  - timezone

### Delete Entry

- Path: `DELETE /v1/entry/{id}`

### List Entries

- Path: `GET /v1/entry`
- Query Param
  - start_date_time (default: start of current day)
  - end_date_time (default: end of current day)
  - timezone (default: EST)
