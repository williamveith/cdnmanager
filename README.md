# CDN Manager

CDN Manager is a native desktop application for managing **Cloudflare Workers KV** through a local desktop interface.

Instead of working directly in the Cloudflare dashboard, the app syncs KV records into a local **SQLite** database so records can be searched, reviewed, inserted, and deleted quickly from a desktop UI.

The application is built with:

- **Go**
- **Wails**
- **SQLite**
- **Vanilla JavaScript**
- **Fuse.js** for approximate table search

macOS releases are packaged as a desktop app and distributed as a **Developer ID signed** and **Apple notarized** download.

---

## Features

### Cloudflare KV management
- Syncs records from a Cloudflare Workers KV namespace into a local SQLite database
- Inserts KV records with structured metadata
- Deletes KV records from Cloudflare and the local database
- Uses the configured domain to generate shareable entry links

### Local database workflow
- Maintains a local SQLite cache of KV records for fast lookups
- Supports search by:
  - all entries
  - UUID
  - URL (single result)
  - URL (multiple results)
- Supports approximate in-table search using Fuse.js on metadata fields:
  - name
  - mimetype
  - location
  - description

### Metadata support
Each KV record can include structured metadata:

- `name`
- `external`
- `mimetype`
- `location`
- `cloud_storage_id`
- `md5Checksum`
- `description`

### Bulk insertion support
- Downloadable CSV template for bulk inserts
- CSV-based insert workflow from the desktop UI

### Desktop-focused behavior
- Native desktop shell via Wails
- Finder reveal for downloaded bulk insert template
- Cross-platform window configuration for macOS, Windows, and Linux

### macOS distribution
- Developer ID signed
- Apple notarized
- Distributed as a DMG installer

---

## Project Structure

```text
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
│   │   │   ├── fonts
│   │   │   └── images
│   │   ├── main.js
│   │   └── styles
│   │       └── app.css
│   └── wailsjs
├── go.mod
├── go.sum
├── hooks
│   ├── postbuild.sh
│   └── prebuild.sh
├── main.go
├── pkg
│   ├── config
│   │   └── config.go
│   ├── database
│   │   └── database.go
│   ├── models
│   │   └── models.go
│   └── session
│       └── session.go
└── wails.json
````

### Key components

#### `main.go`

Application entry point. Responsible for:

* resolving app data/config paths
* initializing the config file
* initializing the SQLite database from embedded schema
* embedding frontend assets and schema
* launching the Wails desktop app

#### `app.go`

Main application bridge exposed to the frontend. Handles:

* config checks and saves
* retrieving the configured domain
* Cloudflare session initialization
* syncing Cloudflare KV into the local database
* insert/delete actions
* generating the CSV bulk insert template

#### `pkg/config`

Defines the application config structure and JSON load/save behavior.

#### `pkg/session`

Wraps Cloudflare Workers KV API operations, including:

* listing keys
* reading values
* inserting entries
* deleting entries
* loading all entries from KV
* concurrent retrieval of KV values from known storage keys

#### `pkg/database`

Local SQLite cache layer for:

* creating/dropping the records table
* inserting entries
* deleting entries
* querying records by UUID or value
* retrieving all cached entries

#### `pkg/models`

Shared record and metadata models with JSON serialization helpers.

#### `frontend/src/main.js`

Main UI logic for:

* initial config flow
* syncing on startup
* searching and rendering the table
* approximate search
* manual insert form
* CSV insert workflow
* delete workflow
* copy-to-clipboard behavior
* dynamic link generation based on configured domain

---

## How It Works

### 1. First launch

On startup, the app:

* determines the platform-specific user config directory
* creates an app folder named `cdnmanager`
* creates `config.json` if it does not exist
* creates `cdnmanager.sqlite3` if it does not exist
* initializes the frontend

### 2. Configuration

If the config is incomplete, the app presents a setup form asking for:

* Cloudflare email
* Cloudflare API key
* Cloudflare account ID
* Cloudflare namespace ID
* domain

The frontend validates these inputs before submission.

### 3. Sync

After configuration is saved, the app:

* creates a Cloudflare session
* fetches KV keys
* compares Cloudflare KV count against the local database count
* rebuilds the local table if the counts differ

### 4. Search and browse

The local SQLite cache is used to display and search records quickly.

### 5. Insert and delete

Manual inserts and CSV inserts write to:

* Cloudflare Workers KV
* the local SQLite cache

Deletes remove records from:

* Cloudflare Workers KV
* the local SQLite cache

---

## Configuration File

The application stores configuration in a JSON file named `config.json`.

Example:

```json
{
  "cloudflare_email": "you@example.com",
  "cloudflare_api_key": "your_api_token",
  "account_id": "0123456789abcdef0123456789abcdef",
  "namespace_id": "fedcba9876543210fedcba9876543210",
  "domain": "cdn.example.com"
}
```

### Config fields

* `cloudflare_email`: Cloudflare account email
* `cloudflare_api_key`: Cloudflare API token/key used by the app
* `account_id`: Cloudflare account ID
* `namespace_id`: Workers KV namespace ID
* `domain`: base domain used to generate entry links

The app normalizes the configured domain so both of these work:

* `cdn.example.com`
* `https://cdn.example.com`

---

## Local Data Paths

The app stores its files under the user's config directory in an app folder named `cdnmanager`.

Files created by the app:

* `config.json`
* `cdnmanager.sqlite3`

The exact base directory depends on the OS, based on `os.UserConfigDir()`.

---

## Database Schema

The SQLite database stores records in a table named `records` with:

* `name` as the primary key
* `value`
* `metadata` stored as JSON text

This database acts as a local cache of Cloudflare KV data.

---

## Search Modes

The UI supports four search modes:

* **All** — retrieve all cached entries
* **By UUID** — retrieve a single entry by ID
* **By URL (single)** — retrieve one entry by exact value
* **By URL (multiple)** — retrieve all entries with the same value

After results are loaded, the table also supports approximate filtering through Fuse.js.

---

## Bulk Insert CSV Template

The app can generate a CSV template named:

```text
CDN Manager Bulk Insert Template.csv
```

The file is saved to the user's `Downloads` folder and revealed in Finder.

Expected header:

```csv
name,value,metadata_name,metadata_external,metadata_mimetype,metadata_location,metadata_description,metadata_cloud_storage_id,metadata_md5Checksum
```

### Example row

```csv
635ce241-ea02-4faf-b888-295f522e7cb4,https://example.com/file.pdf,Project File,false,application/pdf,cdn.example.com owner@example.com,Example document,abc123,d41d8cd98f00b204e9800998ecf8427e
```

---

## Frontend Validation

The setup form validates:

* **Cloudflare Email** as an email input
* **Cloudflare API Key** using a token-like regex
* **Account ID** as 32 lowercase hex characters
* **Namespace ID** as 32 lowercase hex characters
* **Domain** as a valid hostname or URL-like domain string

---

## Link Generation

Entry links are built from the configured domain.

For a record ID like:

```text
635ce241-ea02-4faf-b888-295f522e7cb4
```

the generated link becomes:

```text
https://your-domain.example/?id=635ce241-ea02-4faf-b888-295f522e7cb4
```

If the domain has not yet been loaded, the frontend temporarily falls back to:

```text
?id=<uuid>
```

---

## Build Requirements

### Core requirements

* Go
* Node.js / npm
* Wails v2
* macOS, Windows, or Linux development environment

### macOS release requirements

For signed and notarized macOS builds:

* Apple Developer account
* Developer ID Application certificate
* Xcode command line tools
* `notarytool` credentials stored in a keychain profile

---

## Development

### Install frontend dependencies

```bash
cd frontend
npm install
```

### Run in development mode

```bash
wails dev
```

### Run checks

```bash
make check
```

### Build the app

```bash
make build
```

---

## Release Workflow

The project includes a `Makefile` for building and packaging release artifacts.

Release pipeline steps:

1. clean old build artifacts
2. run environment checks
3. build the Wails app
4. sign the app bundle
5. stage DMG contents
6. create the DMG
7. notarize the DMG with Apple
8. staple the notarization ticket
9. verify the notarized release artifact

### Common commands

```bash
make build
make release
make clean
```

---

## Current Architecture Notes

* Config logic is centralized in `pkg/config`
* Cloudflare session initialization uses the shared config type directly
* The frontend dynamically uses the configured domain instead of a hardcoded URL
* KV value loading from storage keys is parallelized with workers in `pkg/session`
* The local SQLite database is used as the app's fast searchable cache layer

---

## Tech Stack

* **Go**
* **Wails v2**
* **Cloudflare Go SDK**
* **SQLite**
* **Vanilla JavaScript**
* **Fuse.js**
* **HTML/CSS**

---

## Version

Current release:

```text
v2.1.0
```

This version includes:

* config package refactor
* dynamic domain-based link generation
* improved setup form validation
* concurrent KV entry loading
* signed and notarized macOS release packaging

---

## Author

**William Veith**
