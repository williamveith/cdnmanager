# CDN Manager

CDN Manager is a native desktop application for managing Cloudflare Workers KV through a local desktop interface.

Instead of working directly in the Cloudflare dashboard, the app syncs KV records into a local SQLite database so records can be searched, reviewed, inserted, and deleted quickly from a desktop UI.

The application is built with:

* Go
* Wails
* SQLite
* Vanilla JavaScript
* Fuse.js for approximate table search

macOS releases are packaged as a universal DMG containing a universal app bundle for both Apple Silicon (arm64) and Intel (x86_64) Macs. Release builds are Developer ID signed and Apple notarized.

---

## Features

### Cloudflare KV management

* Syncs records from a Cloudflare Workers KV namespace into a local SQLite database
* Inserts KV records with structured metadata
* Deletes KV records from Cloudflare and the local database
* Uses the configured domain to generate shareable entry links

### Local database workflow

* Maintains a local SQLite cache of KV records for fast lookups

* Supports search by:

  * all entries
  * UUID
  * URL (single result)
  * URL (multiple results)

* Supports approximate in-table search using Fuse.js on metadata fields:

  * name
  * mimetype
  * location
  * description

### Metadata support

Each KV record can include structured metadata:

* name
* external
* mimetype
* location
* cloud_storage_id
* md5Checksum
* description

### Bulk insertion support

* Downloadable CSV template for bulk inserts
* CSV-based insert workflow from the desktop UI

---

## macOS Distribution

* Universal macOS app bundle (arm64 + x86_64)
* Universal DMG installer
* Developer ID signed
* Apple notarized
* Gatekeeper-compatible distribution

---

## Developer Release Pipeline

### Source provenance

* Signed Git commits
* Signed Git tags
* Tagged GitHub releases

### Build and distribution

* Universal macOS build generation
* App bundle signing (Developer ID)
* DMG packaging
* Apple notarization
* Stapling and verification

> Note: Git tagging and GitHub release publishing are performed outside of the Makefile.

---

## Project Structure

```
.
├── Makefile
├── README.md
├── app.go
├── data
│   └── schema.sql
├── frontend
│   ├── dist
│   ├── index.html
│   ├── package-lock.json
│   ├── package.json
│   ├── package.json.md5
│   ├── src
│   │   ├── assets
│   │   ├── main.js
│   │   └── styles
│   └── wailsjs
├── go.mod
├── go.sum
├── hooks
├── main.go
├── pkg
│   ├── config
│   ├── database
│   ├── models
│   ├── reconcile
│   └── session
└── wails.json
```

---

## Key Components

### `main.go`

Application entry point. Responsible for:

* resolving app data/config paths
* initializing the config file
* initializing the SQLite database from embedded schema
* embedding frontend assets and schema
* launching the Wails desktop app

---

### `app.go`

Main application bridge exposed to the frontend. Handles:

* config checks and saves
* retrieving the configured domain
* Cloudflare session initialization
* syncing Cloudflare KV into the local database
* insert/delete actions
* generating the CSV bulk insert template

---

### `pkg/config`

Defines the application config structure and JSON load/save behavior.

---

### `pkg/session`

Wraps Cloudflare Workers KV API operations, including:

* listing keys
* reading values
* inserting entries
* deleting entries
* loading KV entries
* concurrent retrieval of KV values

---

### `pkg/reconcile`

Responsible for synchronizing Cloudflare KV state with the local SQLite database.

* Compares remote KV data against the local cache
* Determines when a full rebuild or update is required
* Ensures local database consistency with Cloudflare
* Acts as the coordination layer between `pkg/session` and `pkg/database`

This package isolates synchronization logic from both the API layer and storage layer.

---

### `pkg/database`

Local SQLite cache layer for:

* creating/dropping the records table
* inserting entries
* deleting entries
* querying records
* retrieving cached entries

---

### `pkg/models`

Shared record and metadata models with JSON serialization helpers.

---

### `frontend/src/main.js`

Handles:

* initial config flow
* syncing on startup
* search and rendering
* approximate search
* insert/delete workflows
* CSV processing
* dynamic link generation

---

## How It Works

### 1. First launch

* resolves user config directory
* creates `cdnmanager/`
* initializes `config.json` and database
* loads frontend

---

### 2. Configuration

Prompts for:

* Cloudflare API token
* Cloudflare account ID
* Cloudflare namespace ID
* domain

---

### 3. Sync

* initializes Cloudflare session
* retrieves KV state
* runs reconciliation against local DB via `pkg/reconcile`
* updates or rebuilds cache as needed

---

### 4. Search and browse

All reads occur against the local SQLite cache.

---

### 5. Insert and delete

Writes are applied to:

* Cloudflare Workers KV
* local SQLite database

---

## Configuration File

`config.json`

```json
{
  "cloudflare_api_token": "your_api_token",
  "account_id": "0123456789abcdef0123456789abcdef",
  "namespace_id": "fedcba9876543210fedcba9876543210",
  "domain": "cdn.example.com"
}
```

### Fields

* `cloudflare_api_token`
* `account_id`
* `namespace_id`
* `domain`

---

## Local Data Paths

```
cdnmanager/
├── config.json
└── cdnmanager.sqlite3
```

---

## Database Schema

Table: `records`

* `name` (primary key)
* `value`
* `metadata` (JSON text)

---

## Search Modes

* All
* By UUID
* By URL (single)
* By URL (multiple)

---

## Bulk Insert CSV Template

```
CDN Manager Bulk Insert Template.csv
```

Header:

```
name,value,metadata_name,metadata_external,metadata_mimetype,metadata_location,metadata_description,metadata_cloud_storage_id,metadata_md5Checksum
```

---

## Link Generation

```
https://your-domain.example/?id=<uuid>
```

Fallback:

```
?id=<uuid>
```

---

## Build Requirements

### Core

* Go
* Node.js / npm
* Wails v2

### macOS

* Apple Developer account
* Developer ID certificate
* Xcode CLI tools
* `notarytool` configured

---

## Development

```bash
cd frontend
npm install

wails dev
```

---

## Commands

```bash
make build
make release
make clean
```

---

## Architecture Notes

* Config centralized in `pkg/config`
* Token-based Cloudflare auth
* Reconciliation logic isolated in `pkg/reconcile`
* SQLite used as primary read layer
* Parallel KV retrieval

---

## Tech Stack

* Go
* Wails v2
* Cloudflare Go SDK
* SQLite
* Vanilla JS
* Fuse.js

---

## Author

William Veith
