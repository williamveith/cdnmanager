# CDN Manager

A native desktop application for managing **Cloudflare Workers KV** through a fast local interface.

CDN Manager syncs KV data into a local **SQLite** database so entries can be browsed, searched, edited, and managed from a desktop UI instead of the Cloudflare dashboard.

Built with **Go**, **Wails**, and a bundled frontend, CDN Manager is now a **general-use desktop app** that can be packaged, shared, and installed through a standard macOS **DMG** installer.

---

## Features

- Native desktop app built with **Wails**
- Local **SQLite** cache of Cloudflare Workers KV data
- First-launch configuration flow for Cloudflare credentials
- Automatic startup sync with Cloudflare once configured
- Embedded frontend assets for standalone distribution
- Embedded database schema initialization
- CSV template generation for bulk insert workflows
- Finder integration for generated files
- Simple macOS packaging pipeline with **DMG** output
- Shareable installer with per-user local configuration

---

## How It Works

On startup, CDN Manager:

1. initializes a local application directory
2. creates a local `config.json` if one does not already exist
3. initializes a persistent local SQLite database from the embedded schema if needed
4. checks whether the Cloudflare configuration is complete
5. if configuration is incomplete, shows a setup form
6. once configured, creates a Cloudflare session at runtime
7. synchronizes Cloudflare KV data into the local database
8. launches the full desktop UI

The local database and configuration are stored outside the app bundle, so they persist across rebuilds, reinstalls, and upgrades.

---

## Installation

### macOS

1. Download the latest DMG from **Releases**
2. Open the disk image
3. Drag **CDN Manager.app** into **Applications**
4. Launch the application
5. On first launch, enter your Cloudflare credentials in the setup form

CDN Manager no longer requires users to clone the repository or rebuild the application just to use it.

---

## First Launch Setup

On first launch, CDN Manager creates a local configuration file and prompts the user for:

- Cloudflare email
- Cloudflare API key
- Cloudflare account ID
- Workers KV namespace ID
- domain

Once that information is saved, the application initializes the Cloudflare session and syncs the local SQLite cache.

---

## Local Application Data

On macOS, the application stores its persistent data in a path like:

```text
~/Library/Application Support/cdnmanager/
````

Important files include:

```text
~/Library/Application Support/cdnmanager/config.json
~/Library/Application Support/cdnmanager/cdnmanager.sqlite3
```

* `config.json` stores the user's Cloudflare configuration
* `cdnmanager.sqlite3` stores the local KV cache

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

## Development Requirements

To build from source, install:

* [Git](https://git-scm.com/downloads)
* [Go](https://go.dev/dl/)
* [Node.js](https://nodejs.org/)
* npm
* [Wails](https://wails.io/docs/gettingstarted/installation)

Install Wails with:

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

On macOS, you should also have:

```text
Xcode Command Line Tools
```

---

## Development Setup

Clone the repository:

```bash
git clone https://github.com/williamveith/cdnmanager.git
cd cdnmanager
```

Run the development server:

```bash
make test
```

Because configuration is now handled at runtime, build-time `.env` setup is no longer required for normal app usage.

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

This builds the frontend, compiles the Go backend, and packages the native macOS app bundle.

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
build → stage-dmg → dmg
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

## Project Structure

```text
cdnmanager
├── Makefile
├── README.md
├── app.go
├── config.go
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
