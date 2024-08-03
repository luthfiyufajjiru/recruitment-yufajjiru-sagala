# Requirement
- Go 1.22
- Git

# Architecture
This is recruitment assignment but this repository is build from zero by the author: luthfiyufajjiru@gmail.com including it is architecture.

# Getting Started
You need sqlite file `default.db` or anything (adjust it on env later)
After that migrate the db using this command from root directory:
```
go run ./migration/main.go
```

You need to place `.env` file in the root directory with value:
```
ADDR="localhost:8080"
DSN_DEFAULT="default.db"
HARD_DELETE=""
``` 
Run this command in the root directory for the first time
```
go install &&
go mod tidy &&
go run main.go
```

Or just `go run main.go` after first run.