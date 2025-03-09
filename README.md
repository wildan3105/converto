# converto
A service to convert .shapr to various format; available via API and CLI


## Usage
TBD

### Arguments
TBD

### Options
TBD

## Local Development
### Requirements
- Go 1.21 or higher

### Run locally
```bash
TBD
```

### Build Locally
TBD

### Code Formatting
```bash
gofmt -w -s .
```

### Linting
```bash
golangci-lint run
```

### Run Tests

⚠ **WARNING** ⚠

> Test is ran against the local development database. So make sure the data is clean before running the test to ensure its accuracy
> Also ensure the server and worker is running

```bash
go test ./...
```

### Clear Test Cache
```bash
go clean -testcache
```