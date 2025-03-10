# converto
A service to convert a `.shapr` file to various format via API


# Usage
## API Documentation

### File Upload and Conversion
<details>
<summary>POST /api/v1/conversions</summary>

**Description:**
Uploads a .shapr file and initiates conversion to specified format. Adds the job to the queue and returns a conversion ID.

**Request Type:** `multipart/form-data`

### Request Fields
| Field Name    | Type   | Description                                              | Required |
|---------------|---------|----------------------------------------------------------|-----------|
| `file`        | file    | The .shapr file to convert                                | Yes       |
| `target_format` | string  | Desired output format (`.step`, `.iges`, `.stl`, `.obj`)  | Yes       |

### Example Response
```json
{
    "id": "67cf6e74dcb672239857517a",
    "status": "pending",
    "message": "Conversion created successfully"
}
```
</details>

### List All Conversions
<details>
<summary>GET /api/v1/conversions</summary>

**Description:**
Retrieves all conversion jobs with their status, progress, and file URLs. Supports optional pagination and filtering by status.

### Query Parameters
| Parameter | Type | Description                                            | Required |
|-----------|-------|--------------------------------------------------------|-----------|
| `status`  | string | Filter by status (`pending`, `in_progress`, `completed`, `failed`) | No        |
| `page`    | int    | Page number for pagination                              | No        |
| `limit`   | int    | Number of results per page                              | No        |

### Example Response
```json
{
    "page": 1,
    "limit": 10,
    "data": [
        {
            "id": "67cf6e74dcb672239857517a",
            "status": "completed",
            "progress": 100,
            "original_file_path": "/home/wildan/original/6bb07b15-a056-4756-bf1d-03ba1f50dff1/one.shapr",
            "converted_file_path": "/home/wildan/converted/6bb07b15-a056-4756-bf1d-03ba1f50dff1/one.iges"
        }
    ]
}
```
</details>

### Get Conversion by ID
<details>
<summary>GET /api/v1/conversions/{conversion_id}</summary>

**Description:**
Retrieves the status and progress of a specific conversion.

### Example Response
```json
{
    "id": "67cf6e74dcb672239857517a",
    "status": "completed",
    "progress": 100,
    "original_file_path": "/home/wildan/original/6bb07b15-a056-4756-bf1d-03ba1f50dff1/one.shapr",
    "converted_file_path": "/home/wildan/converted/6bb07b15-a056-4756-bf1d-03ba1f50dff1/one.iges"
}
```
</details>

### Download Original File
<details>
<summary>GET /api/v1/conversions/{conversion_id}/files?type=original</summary>

**Description:**
Allows users to download the original uploaded .shapr file.

**Example Request:**
```http
GET /api/v1/conversions/12345/files?type=original
```

**Response:**
Returns the original file as raw data in the response body.
</details>

### Download Converted File
<details>
<summary>GET /api/v1/conversions/{conversion_id}/files?type=converted</summary>

**Description:**
Allows users to download the converted file if the conversion is completed.

**Example Request:**
```http
GET /api/v1/conversions/12345/files?type=converted
```

**Response:**
Returns the converted file as raw data in the response body.
</details>

## Local Development
### Requirements
- Go `v1.21` or higher
- Docker with `docker-compose` to provision dependencies

### Run locally
1. Provision dependencies: `docker-compose up -d`.
2. Copy the value in [.env.example](./.env.example) to your own `.env` and adjust the value accordingly.
3. Run the application
```bash
# run server
go run main.go server

# run worker 
go run main.go worker
```

**Note**
Default port is `3000`

### Build Locally
```bash
# first, build the app
go build -o app .

# and then run the server and worker in different terminal
./app server
./app worker
```

Alternatively, if you want to start both server and worker in a single command, you could use [forego](https://github.com/ddollar/forego):
```bash
forego start # assuming you have installed and configured forego in your system
```

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