# Converto
Convert `.shapr` files to various formats via API.

---
## 📼 Demo
![demo](./docs/demovideo.webm)

## 📑 API Documentation

### 📂 File Upload & Conversion
<details>
<summary><code>POST /api/v1/conversions</code></summary>

**Description:** Uploads a `.shapr` file and initiates conversion to a specified format. Returns a conversion ID.

**Request Type:** `multipart/form-data`

#### 🔍 Request Fields
| Field Name      | Type   | Description                                            | Required |
|-----------------|---------|--------------------------------------------------------|-----------|
| `file`          | file    | The `.shapr` file to convert                            | ✅ Yes    |
| `target_format` | string  | Output format (`.step`, `.iges`, `.stl`, `.obj`)       | ✅ Yes    |

#### 📥 Example Response
```json
{
    "id": "67cf6e74dcb672239857517a",
    "status": "pending",
    "message": "Conversion created successfully"
}
```
</details>

### 📜 List All Conversions
<details>
<summary><code>GET /api/v1/conversions</code></summary>

**Description:** Retrieves all conversion jobs with status, progress, and file URLs. Supports pagination and status filtering.

#### 🔍 Query Parameters
| Parameter | Type | Description                                           | Required |
|-----------|-------|-------------------------------------------------------|-----------|
| `status`  | string | Filter by status (`pending`, `in_progress`, `completed`, `failed`) | ❌ No     |
| `page`    | int    | Page number for pagination                             | ❌ No     |
| `limit`   | int    | Number of results per page                             | ❌ No     |

#### 📥 Example Response
```json
{
    "page": 1,
    "limit": 10,
    "data": [
        {
            "id": "67cf6e74dcb672239857517a",
            "status": "completed",
            "progress": 100,
            "original_file_path": "/path/to/original.shapr",
            "converted_file_path": "/path/to/converted.iges"
        }
    ]
}
```
</details>

### 📌 Get Conversion by ID
<details>
<summary><code>GET /api/v1/conversions/{conversion_id}</code></summary>

**Description:** Retrieves the status and progress of a specific conversion.

#### 📥 Example Response
```json
{
    "id": "67cf6e74dcb672239857517a",
    "status": "completed",
    "progress": 100,
    "original_file_path": "/path/to/original.shapr",
    "converted_file_path": "/path/to/converted.iges"
}
```
</details>

### 📤 Download Original File
<details>
<summary><code>GET /api/v1/conversions/{conversion_id}/files?type=original</code></summary>

**Description:** Downloads the original uploaded `.shapr` file.

#### 📥 Example Request
```http
GET /api/v1/conversions/12345/files?type=original
```

**Response:** Returns the original file as raw data.
</details>

### 📤 Download Converted File
<details>
<summary><code>GET /api/v1/conversions/{conversion_id}/files?type=converted</code></summary>

**Description:** Downloads the converted file if the conversion is completed.

#### 📥 Example Request
```http
GET /api/v1/conversions/12345/files?type=converted
```

**Response:** Returns the converted file as raw data.
</details>

---
## 🚀 Local Development

### 🔧 Prerequisites
- Go `v1.21+`
- Docker & `docker-compose`

### 🛠️ Setup
1. **Start Dependencies:**
```bash
docker-compose up -d
```

2. **Configure Environment:**
Copy `.env.example` to `.env` and update values as needed.

3. **Run Application:**
```bash
# Start server
go run main.go server

# Start worker
go run main.go worker
```

### 📦 Build & Run Binary
```bash
# Build the app
go build -o app .

# Run server & worker in separate terminals
./app server
./app worker
```

**🧠 Pro Tip:** Run both in a single terminal using [forego](https://github.com/ddollar/forego):
```bash
forego start
```

### 🧹 Code Quality
```bash
# Formatting
gofmt -w -s .

# Linting
golangci-lint run
```

### ✅ Run Tests

⚠️ **Caution**

The tests run against the local development database and will delete all data within it. Additionally, any files inside [`BASE_DIRECTORY`](.env.example#L10) will be permanently removed. Ensure that no important files are stored there before running the tests.

```bash
# Ensure both server and worker are running

go test ./...
```

### 🧽 Clear Test Cache
```bash
go clean -testcache
```

---

### 📍 Default Port
- Server runs on: `http://localhost:3000`

### 🤖 Miscellaneous
For detailed architecture, deployment strategy, and future improvements, please take a look [here](./docs/README.md)