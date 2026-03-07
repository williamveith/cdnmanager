# CDN Manager

A native desktop application for managing **Cloudflare Workers KV** through a fast local interface.

CDN Manager syncs KV data into a local **SQLite** database so entries can be browsed, searched, edited, and managed from a desktop UI instead of the Cloudflare dashboard.

Built with **Go**, **Wails**, and a bundled frontend, the app currently targets **macOS** as its primary release platform.

---

## Features

- Native desktop app built with **Wails**
- Local **SQLite** cache of Cloudflare Workers KV data
- Automatic startup sync with Cloudflare
- Embedded frontend assets for standalone distribution
- Embedded database schema initialization
- CSV template generation for bulk insert workflows
- Finder integration for generated files
- Simple macOS packaging pipeline with **DMG** output

---

## How It Works

On startup, CDN Manager:

1. loads embedded configuration from `.env`
2. loads the embedded SQLite schema
3. initializes a persistent local database
4. creates a Cloudflare session
5. compares local data against Cloudflare KV
6. refreshes the local database if the remote state has changed
7. launches the desktop UI

The local database is stored outside the app bundle so it persists across rebuilds and upgrades.

---

## Requirements

To build from source, install:

- [Git](https://git-scm.com/downloads)
- [Go](https://go.dev/dl/)
- [Node.js](https://nodejs.org/)
- npm
- [Wails](https://wails.io/docs/gettingstarted/installation)

Install Wails with:

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
````

On macOS, you should also have:

```text
Xcode Command Line Tools
```

---

## Setup

Clone the repository:

```bash
git clone https://github.com/williamveith/cdnmanager.git
cd cdnmanager
```

The application requires a `.env` file containing your Cloudflare credentials.

Run:

```bash
make
```

If `.env` does not exist, the Makefile will create it and stop so you can fill in the required values.

Generated `.env` template:

```env
cloudflare_email=""
cloudflare_api_key=""
account_id=""
namespace_id=""
domain=""
```

After filling in the values, run the build again.

---

## Development

Run the Wails development server:

```bash
make test
```

---

## Build

Build the application with:

```bash
make
```

or:

```bash
make build
```

This runs dependency checks, validates `.env`, builds the frontend, compiles the Go backend, and packages the native app.

Output:

```text
build/bin/cdnmanager.app
```

---

## Package as a DMG

Create a distributable macOS disk image with:

```bash
make dmg
```

Output:

```text
build/bin/CDN-Manager-v1.1.1.dmg
```

---

## Release Build

Run the full release pipeline with:

```bash
make release
```

This executes:

```text
check → build → stage-dmg → dmg
```

Artifacts produced:

```text
build/bin/cdnmanager.app
build/bin/CDN-Manager-v1.1.1.dmg
```

---

## Available Make Commands

Build the app:

```bash
make
```

Run the development server:

```bash
make test
```

Create the DMG:

```bash
make dmg
```

Run the full release build:

```bash
make release
```

Launch the built app:

```bash
make run
```

Clean build artifacts:

```bash
make clean
```

---

## Database Location

On macOS, the persistent database is stored at a path like:

```text
~/Library/Application Support/cdnmanager/cdnmanager.sqlite3
```

---

## Bulk Insert Template

The application can generate a CSV template for bulk insertion.

Generated file:

```text
CDN Manager Bulk Insert Template.csv
```

It is written to the user's `Downloads` folder and revealed in Finder.

Template columns:

```csv
name,value,metadata_name,metadata_external,metadata_mimetype,metadata_location,metadata_description,metadata_cloud_storage_id,metadata_md5Checksum
```

---

## Installation

### macOS

1. Download the latest DMG from Releases
2. Open the disk image
3. Drag **CDN Manager.app** into **Applications**
4. Launch the application

---

## Project Structure

```text
cdnmanager
├── Makefile
├── README.md
├── app.go
├── build
│   ├── bin
│   ├── darwin
│   └── dmg
├── data
│   └── schema.sql
├── frontend
│   ├── dist
│   ├── src
│   └── wailsjs
├── hooks
│   ├── postbuild.sh
│   └── prebuild.sh
├── main.go
├── pkg
│   ├── database
│   ├── models
│   └── session
├── go.mod
├── go.sum
└── wails.json
```

---

## Tech Stack

* Go
* Wails
* SQLite
* Cloudflare Workers KV
* Node.js
* npm

---

## Author

William Veith

```

A strong GitHub README usually also includes one screenshot near the top.
```
