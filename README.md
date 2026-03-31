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

## Project Structure

```
.
├── Makefile                     # Build, package, and release automation
├── README.md                    # Project documentation
├── app.go                       # Wails application bindings and exposed methods
├── data
│   └── schema.sql               # SQLite schema used to initialize the local database
├── frontend
│   ├── dist                     # Production frontend build output
│   │   ├── assets
│   │   │   ├── IBMPlexMono-Regular.49ce58b4.woff2          # Bundled IBM Plex Mono font
│   │   │   ├── appicon.d007682b.png                        # Bundled application icon
│   │   │   ├── glyphicons-halflings-regular.fe185d11.woff2 # Bundled glyphicon font
│   │   │   ├── index.24410bbd.css                          # Compiled frontend styles
│   │   │   └── index.6078d637.js                           # Compiled frontend logic
│   │   └── index.html              # Built frontend entrypoint
│   ├── index.html                 # Frontend HTML shell used during development/build
│   ├── package-lock.json          # Locked frontend dependency versions
│   ├── package.json               # Frontend package manifest
│   ├── package.json.md5           # Integrity hash for package manifest
│   ├── src                        # Frontend source code
│   │   ├── assets
│   │   │   ├── fonts
│   │   │   │   ├── IBM_Plex_Mono
│   │   │   │   │   ├── IBMPlexMono-Regular.woff2          # Primary monospace UI font
│   │   │   │   │   └── license.txt                        # IBM Plex Mono license
│   │   │   │   └── glyphicons_halflings
│   │   │   │       ├── glyphicons-halflings-regular.woff2 # Icon font
│   │   │   │       └── license.txt                        # Glyphicon license
│   │   │   └── images
│   │   │       ├── appicon.png          # Raster app icon
│   │   │       └── appicon.svg          # Vector app icon
│   │   ├── controllers                  # Frontend behavior and UI orchestration
│   │   │   ├── appController.js         # Application bootstrap/controller coordination
│   │   │   ├── configController.js      # Setup form behavior and submission handling
│   │   │   ├── deleteController.js      # Delete flow and event handling
│   │   │   ├── insertController.js      # Manual/CSV insert workflows
│   │   │   └── searchController.js      # Search flow and result handling
│   │   ├── main.js                      # Frontend entrypoint
│   │   ├── services                     # Thin wrappers around Wails/Go APIs
│   │   │   ├── appService.js            # App-level backend service calls
│   │   │   └── dbService.js             # Database-related backend service calls
│   │   ├── state
│   │   │   └── appState.js              # Shared frontend state
│   │   ├── styles
│   │   │   └── app.css                  # Application styling
│   │   ├── utils
│   │   │   ├── clipboard.js             # Clipboard helpers
│   │   │   ├── domain.js                # Domain/link normalization helpers
│   │   │   └── uuid.js                  # UUID helper functions
│   │   └── views                        # DOM rendering and UI templates
│   │       ├── configView.js            # Initial setup/configuration view
│   │       ├── shellView.js             # Main application shell
│   │       └── tableView.js             # Search result table rendering
│   └── wailsjs                          # Auto-generated Wails JS bindings
│       ├── go
│       │   ├── database
│       │   │   ├── Database.d.ts        # TypeScript definitions for DB bindings
│       │   │   └── Database.js          # JS bindings for DB methods
│       │   ├── main
│       │   │   ├── App.d.ts             # TypeScript definitions for app bindings
│       │   │   └── App.js               # JS bindings for app methods
│       │   └── models.ts                # Shared generated model definitions
│       └── runtime
│           ├── package.json             # Wails runtime package metadata
│           ├── runtime.d.ts             # Wails runtime TypeScript definitions
│           └── runtime.js               # Wails runtime helpers
├── go.mod                         # Go module definition
├── go.sum                         # Go dependency checksums
├── hooks
│   ├── postbuild.sh               # Post-build automation hook
│   └── prebuild.sh                # Pre-build automation hook
├── main.go                        # Go application entrypoint
├── pkg                            # Internal Go packages
│   ├── config
│   │   └── config.go              # App configuration loading, validation, normalization
│   ├── database
│   │   └── database.go            # SQLite/database access layer
│   ├── models
│   │   └── models.go              # Shared Go data models
│   ├── reconcile
│   │   └── reconcile.go           # Cloudflare/database reconciliation logic
│   └── session
│       └── session.go             # Session/runtime state management
└── wails.json                     # Wails project configuration
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
