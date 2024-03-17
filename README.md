# Lenslocked

Project for Jon Calhoun's Web Development with Go course.

## Development

### Serve Locally

Pre-requisites:

- Docker
- Go installed to run locally

```bash
# from the root folder
$ go run main.go
```

### Live Reload

This project uses `modd` for live reaload.

```bash
# from the project root folder
$ modd
```

### Testing

There are no automated tests. In order to make it work, you'll have to:

1. Run Docker
2. Create a user
    - `username`: admin@user.com
    - `password`: admin

### Connecting to the Database

```bash
# Ensure docker is running, then:
$ docker exec -it lenslocked_db /usr/bin/psql -U baloo -d lenslocked
```

