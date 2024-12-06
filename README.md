# Cloudflare CDN Manager

## About

**Cloudflare CDN Manager** is a tool I’m using as a personal **URL Redirect Management Service**. It manages Cloudflare Workers KV storage through a straightforward interface and runs on Windows*, macOS, Linux*, or Docker*. Built with the Wails framework and Go, it syncs with a local SQLite3 database to keep track of entries and their metadata.

> *Haven't gotten around to creating the Windows, Linux, or Docker versions yet since I'm still doing local development. If you can see that message, that is still currently the case. You can easily modify `wails.Run` in `main.go` to generate a native app for Linux and Windows. Docker... idk haven't looked into it yet. I'm guessing someone has already made a base image with WAILS.*

This setup is also structured to make the next phase easier to implement: adding a second web worker that fetches the raw content a URL points to and wraps it in MIME-appropriate HTML. The goal is to serve the content directly to users instead of relying on cloud storage systems to display it.

---

## Build App

1. **Install Build Tools**:
   - [git](https://git-scm.com/downloads)
   - [Go](https://go.dev/dl/)
   - [NPM](https://nodejs.org/en/download/package-manager)
   - [Wails](https://wails.io/docs/gettingstarted/installation#installing-wails)

   git is used to clone the source code for this application from a remote repository. The application itself (fetching, pushing, transforming data) is written in Go. Wails creates a webkit frontend using a JavaScript framework managed by NPM and generates JavaScript bindings for the Go application so the front end can communicate with the application.

2. **Clone Repo**

   ```bash
   git clone https://github.com/williamveith/cdnmanager.git
   cd cdnmanager
   mv template.env .env;
   ```

   Clones the remote repository containing the application source code, enters the root directory of the repository, and creates your .env file.

3. **Add Cloudflare Credentials**
   Fill out the `.env` file with your Cloudflare account details. This file will be embedded into the application during the build process, allowing the application to make API calls to Cloudflare.

4. **Build Application**:

    ```bash
    make build
    ```

   This will clean the build directory, compile the project, and generate the binary.

5. **Run the Application**:
    After building, the binary will be available in the `build/bin` directory
