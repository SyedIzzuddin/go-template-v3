# API Documentation

## Overview

This REST API provides comprehensive user and file management capabilities with a clean, standardized response format. Built with Go and Echo framework, it offers robust features including file uploads, pagination, and proper error handling.

## Base URL

```
http://localhost:8080/api/v1
```

## Authentication

This API does not currently implement authentication, but includes JWT configuration for future implementation. All endpoints are currently public.

## Response Format

All successful responses follow this format:

```json
{
  "success": true,
  "message": "Operation successful",
  "data": {
    // Response data here
  }
}
```

Error responses:

```json
{
  "success": false,
  "message": "Error description",
  "error": {
    // Error details here
  }
}
```

## Health Check Endpoint

### Get API Health

**GET** `/health`

Returns the health status of the API and its database connection.

## User Endpoints

### Create a New User

**POST** `/users`

**Request Body:**
```json
{
  "name": "John Doe",
  "email": "john.doe@example.com"
}
```

### Get All Users

**GET** `/users`

Supports pagination via query parameters:
- `page`: The page number to retrieve (defaults to 1).
- `per_page`: The number of items per page (defaults to 10).

Example: `GET /users?page=2&per_page=20`

### Get User by ID

**GET** `/users/{id}`

### Update User

**PUT** `/users/{id}`

**Request Body:**
```json
{
  "name": "John Doe Updated"
}
```

### Delete User

**DELETE** `/users/{id}`

## File Endpoints

### Upload a File

**POST** `/files/upload`

This endpoint uses a `multipart/form-data` request.

**Form Fields:**
- `file`: The file to upload.
- `description` (optional): A description of the file.
- `category` (optional): A category for the file.

### Get All Files

**GET** `/files`

Supports pagination via query parameters:
- `page`: The page number to retrieve (defaults to 1).
- `per_page`: The number of items per page (defaults to 10).

### Get My Files

**GET** `/files/my`

_Note: This endpoint currently does not have authentication, so its behavior might not be as expected. It likely retrieves files based on a hardcoded user ID._

### Get File by ID

**GET** `/files/{id}`

### Update File Metadata

**PUT** `/files/{id}`

**Request Body:**
```json
{
  "description": "An updated description for the file.",
  "category": "updated"
}
```

### Delete File

**DELETE** `/files/{id}`

### Download a File

**GET** `/files/{id}/download`

This triggers a file download.

### Serve a File

**GET** `/files/{filename}`

Serves a file directly, which can be used for displaying images or other content in a browser.

### Static File Serving

**GET** `/uploads/*`

Direct access to uploaded files for static serving.

## Environment Setup

Before running the API, ensure you have set the `DATABASE_URL` environment variable for migrations:

```bash
export DATABASE_URL="postgres://username:password@host:port/database?sslmode=disable"
```

Or use the provided docker-compose setup which handles this automatically.

## Error Codes

- `400` - Bad Request (e.g., malformed JSON).
- `404` - Not Found (e.g., user or file with the given ID does not exist).
- `422` - Unprocessable Entity (e.g., validation errors on request body).
- `429` - Too Many Requests (if rate limit is exceeded).
- `500` - Internal Server Error.

## Rate Limiting

- **Limit:** 100 requests per minute per IP address.
- **Response:** A `429 Too Many Requests` error is returned if the limit is exceeded.
