# triple-s

A simplified version of the **Simple Storage Service (S3)** for storing and retrieving files ("objects") inside named containers called **buckets**. Built using **Go**’s standard library, `triple-s` supports creating, listing, and deleting buckets, as well as uploading, retrieving, and deleting objects via a **RESTful API** that returns **XML** responses.

## Table of Contents

1. [Overview](#overview)  
2. [Features](#features)  
   - [Bucket Management](#bucket-management)  
     - [Create Bucket](#create-bucket)  
     - [List Buckets](#list-buckets)  
     - [Delete Bucket](#delete-bucket)  
   - [Object Operations](#object-operations)  
     - [Upload Object](#upload-object)  
     - [Retrieve Object](#retrieve-object)  
     - [Delete Object](#delete-object)  
3. [Directory Structure](#directory-structure)  
4. [Metadata Storage](#metadata-storage)  
   - [Bucket Metadata](#bucket-metadata)  
   - [Object Metadata](#object-metadata)  
5. [Installation and Usage](#installation-and-usage)  
   - [Build](#build)  
   - [Command-Line Arguments](#command-line-arguments)  
   - [Help](#help)  
6. [Error Handling](#error-handling)  
7. [Examples](#examples)  
   - [Create Bucket Example](#create-bucket-example)  
   - [List Buckets Example](#list-buckets-example)  
   - [Delete Bucket Example](#delete-bucket-example)  
   - [Upload Object Example](#upload-object-example)  
   - [Retrieve Object Example](#retrieve-object-example)  
   - [Delete Object Example](#delete-object-example)  
8. [Author](#author)  

---

## Overview

`triple-s` is a learning project that mirrors the core concepts of Amazon S3:
1. **Buckets** serve as containers for your data.
2. **Objects** are your actual files.
3. **XML-based** responses align with S3-like APIs.

---

## Features

### Bucket Management

1. **Create a Bucket** (`PUT /{BucketName}`)  
   - Validates and ensures a **unique** bucket name that follows S3 naming conventions.  
   - Responds with **`200 OK`** if successful, or appropriate error codes for conflicts or invalid names.

2. **List All Buckets** (`GET /`)  
   - Returns an **XML** list of all existing buckets.  
   - Responds with **`200 OK`** and the list in XML.

3. **Delete a Bucket** (`DELETE /{BucketName}`)  
   - Deletes an **empty** bucket.  
   - Responds with **`204 No Content`** on success.  
   - Returns **`404 Not Found`** if the bucket doesn’t exist, or **`409 Conflict`** if it’s not empty.

### Object Operations

1. **Upload Object** (`PUT /{BucketName}/{ObjectKey}`)  
   - Overwrites if the object key already exists.  
   - Accepts the file in the **request body** and stores it under `data/{BucketName}/`.  
   - Metadata is appended or updated in `objects.csv`.  
   - Responds with **`200 OK`** on success.

2. **Retrieve Object** (`GET /{BucketName}/{ObjectKey}`)  
   - Checks if both the bucket and object exist.  
   - Sends the file as a binary stream with appropriate `Content-Type`.  
   - Responds with **`404 Not Found`** if missing.

3. **Delete Object** (`DELETE /{BucketName}/{ObjectKey}`)  
   - Removes the file and its metadata.  
   - Responds with **`204 No Content`** on success, or **`404 Not Found`** if missing.

---

## Directory Structure

```
.
├── triple-s               # Compiled binary after build
├── data/                  # Base directory for storage
│   ├── <bucket-name>/     # Directory named after each bucket
│   │   ├── objects.csv    # CSV metadata for objects in this bucket
│   │   └── <ObjectKey>    # Actual stored object file
├── buckets.csv            # (Optional) Global CSV storing bucket metadata
└── ...
```

- **`data/`**: Root storage folder (configurable with `--dir` argument).
- **`buckets.csv`** (optional approach): Each bucket name and metadata (creation time, etc.).
- **`objects.csv`**: Stores object metadata like `ObjectKey`, `Size`, `ContentType`, `LastModified`.

---

## Metadata Storage

### Bucket Metadata

A CSV (for example, `buckets.csv`) tracks all buckets:
```
Name,CreationTime,LastModifiedTime,Status
my-bucket,2025-01-01T10:00:00Z,2025-01-01T10:00:00Z,active
```
- **Name**: The unique bucket name.  
- **CreationTime**: Timestamp of creation.  
- **LastModifiedTime**: Timestamp of last update.  
- **Status**: Could be `active` or `deleted` (implementation-defined).

### Object Metadata

Each bucket has its own `objects.csv`:
```
ObjectKey,Size,ContentType,LastModified
sunset.png,2048,image/png,2025-01-01T12:00:00Z
```
- **ObjectKey**: The file name or path-like key.  
- **Size**: Size in bytes.  
- **ContentType**: MIME type (optional, but recommended).  
- **LastModified**: Timestamp of last upload or overwrite.

---

## Installation and Usage

### Build

1. Clone or download this repository.
2. In the project root, run:
   ```sh
   go build -o triple-s .
   ```

### Command-Line Arguments

- `--port N` : The port to run the HTTP server. Defaults to `8080` if not specified.
- `--dir S`  : Path to the root data directory. Defaults to `./data` if not specified.

**Example**:
```sh
./triple-s --port 9000 --dir /tmp/triple-s-data
```
Starts the server on port `9000`, storing data in `/tmp/triple-s-data`.

### Help

To see usage information, run:
```sh
./triple-s --help
```

Example output:
```
Simple Storage Service.

Usage:
    triple-s [--port <N>] [--dir <S>]
    triple-s --help

Options:
  --help     Show this screen.
  --port N   Port number
  --dir S    Path to the directory
```

---

## Error Handling

- **400 Bad Request**: Invalid bucket name, object key, or missing parameters.
- **404 Not Found**: Nonexistent bucket or object key.
- **409 Conflict**: Attempting to create a duplicate bucket or delete a non-empty bucket.
- **200 OK / 204 No Content**: Success statuses for bucket/object creation, updates, and deletions.
- **500 Internal Server Error**: Internal unexpected issues (I/O errors, CSV parsing failures, etc.).

All error and success responses (where applicable) should be wrapped in **XML** to conform to the general S3-like specification.

---

## Examples

### Create Bucket Example

**Request**:
```
PUT /my-bucket HTTP/1.1
Host: localhost:8080
```

**Response** (Success):
```xml
HTTP/1.1 200 OK
Content-Type: application/xml

<CreateBucketResult>
  <BucketName>my-bucket</BucketName>
  <Location>/my-bucket</Location>
</CreateBucketResult>
```

### List Buckets Example

**Request**:
```
GET / HTTP/1.1
Host: localhost:8080
```

**Response**:
```xml
HTTP/1.1 200 OK
Content-Type: application/xml

<ListAllMyBucketsResult>
  <Buckets>
    <Bucket>
      <Name>my-bucket</Name>
      <CreationDate>2025-01-01T10:00:00Z</CreationDate>
    </Bucket>
    <Bucket>
      <Name>test-bucket</Name>
      <CreationDate>2025-01-02T14:00:00Z</CreationDate>
    </Bucket>
  </Buckets>
</ListAllMyBucketsResult>
```

### Delete Bucket Example

**Request**:
```
DELETE /my-bucket HTTP/1.1
Host: localhost:8080
```

**Response** (if bucket is empty):
```
HTTP/1.1 204 No Content
```

### Upload Object Example

**Request**:
```
PUT /photos/sunset.png HTTP/1.1
Host: localhost:8080
Content-Type: image/png
Content-Length: 2048

<binary data of the image>
```

**Response**:
```
HTTP/1.1 200 OK
Content-Type: application/xml

<PutObjectResult>
  <ETag>"abc123..."</ETag>
</PutObjectResult>
```

### Retrieve Object Example

**Request**:
```
GET /photos/sunset.png HTTP/1.1
Host: localhost:8080
```

**Response** (success):
```
HTTP/1.1 200 OK
Content-Type: image/png
Content-Length: 2048

<binary data of sunset.png>
```

### Delete Object Example

**Request**:
```
DELETE /photos/sunset.png HTTP/1.1
Host: localhost:8080
```

**Response**:
```
HTTP/1.1 204 No Content
```
