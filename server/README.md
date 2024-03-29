# Backend Server

Go backend server using sqlite3 db.

## Requirements

- Go >= 1.21
- make
- docker

## Configuration

Before building or running backend create `config.json` in `server` dir.

example `config.json`:

```json
{
    "server": {
        "host": "localhost",
        "port": ":3000",
        "redisHost": "localhost",
        "redisPort": "6379",
        "redisPass": "<password>"  
    },
    "mailSender": {
        "email": "<email>",
        "sendgrid_key": "<sendgrid-key>",
        "timeout": 20 
    },
    "database": {
        "file": "./database.db"
    },
    "token": {
        "secret": "mysecret",
        "timeout": 100
    },
    "account": {
        "mnemonics": "<mnemonics>",
        "network": "<grid-network>"
    },
    "version": "v1",
    "salt": "<salt>",
    "admins": [],
    "notifyAdminsIntervalHours": 6,
    "adminSSHKey": "<ssh key>"
}
```

## Build

```bash
make build
```

### Build Using Docker

```bash
docker build -t cloud4students .
```

## Run

```bash
make run
```

### Run Using Docker

```bash
docker run cloud4students
```
