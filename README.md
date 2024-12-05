# Cloudflare CDN Manager

## About

**Cloudflare CDN Manager** is a tool I’m using as a personal **URL Redirect Management Service**. It manages Cloudflare Workers KV storage through a straightforward interface and runs on Windows*, macOS, Linux*, or Docker*. Built with the Wails framework and Go, it syncs with a local SQLite3 database to keep track of entries and their metadata.

    *Haven't gotten around to creating the windows, linux or docker versions yet since Im still doing local development. If you can see that message that is still currently the case. You can easily modify wails.Run in main.go to generate a native app for linux and windows. Docker... idk havent looked into it yet. Im guessing someone has already made a base image with WAILS

This setup is also structured to make the next phase easier to implement: adding a second web worker that fetches the raw content a URL points to and wraps it in MIME-appropriate HTML. The goal is to serve the content directly to users instead of relying on cloud storage systems to display it.

## Features

- **Cross-Platform Compatibility**: Runs on Windows, macOS, Linux, and supports Docker deployments.
- **Cloudflare KV Management**: Interacts with Cloudflare Workers KV storage to manage key-value pairs with metadata.
- **User-Friendly UI**: Provides a graphical interface for managing entries, metadata, and bulk operations.
- **Local Database Control**: Utilizes SQLite3 for local storage and synchronization of entries.
- **Dynamic Search Options**: Search entries by UUID, URL (single or multiple), or retrieve all entries.
- **Predefined Metadata Fields**: Supports structured metadata with predefined fields like Name, External, MimeType, etc.
- **Bulk Insertion from File**: Allows inserting multiple entries at once using a CSV file template.
- **UUID Generation**: Automatically generate UUIDs for new entries.
- **Delete Entries**: Delete entries by UUID with synchronized removal from Cloudflare KV and the local database.
- **Interactive Data Handling**: Sortable tables, clickable cells for easy copying, and dynamic forms.
- **Clipboard Notifications**: Provides on-screen notifications when content is copied.

## Live Development

To run in live development mode, execute:

    ```bash
    wails dev
    ```

This command runs a Vite development server with hot reload capabilities for frontend changes. To develop in a browser and access Go methods, connect to the dev server at [http://localhost:34115](http://localhost:34115).

## Building

To build a production-ready, redistributable package, use:

    ```bash
    wails build
    ```

**Note**: Ensure that the `.env` file is present in the project root directory before building. The `.env` file will be embedded into the binary during the build process.

## Cloudflare Integration

- **API Documentation**: [Cloudflare API Go Docs](https://developers.cloudflare.com/api-next/go/)
- **Examples**: [Cloudflare Go Examples](https://github.com/cloudflare/cloudflare-go/blob/master/workers_kv_example_test.go)

## Configuration

### `.env` File

Create a `.env` file in the project root directory with the following variables:

    ```dotenv
    cloudflare_email=""
    cloudflare_api_key=""
    account_id=""
    namespace_id=""
    domain=""
    ```

Replace the empty strings with your actual Cloudflare account details. This file will be embedded into the application during the build process.

## Backend: Go Modules

### `main.go`

The `main.go` file initializes the application, loads environment variables, synchronizes data between Cloudflare and the local database, and sets up the Wails application.

#### **Functions**

- **`loadEmbeddedEnv()`**: Loads environment variables from an embedded `.env` file into the runtime environment.
- **`SyncFromCloudflare()`**: Synchronizes entries from Cloudflare Workers KV storage to the local SQLite database.
- **`main()`**: The entry point of the application, initializing the database, Cloudflare session, and starting the Wails application.

#### **Notes**

- The application uses an embedded `.env` file for configuration, which is read into environment variables at runtime.
- Synchronization ensures that the local database is up-to-date with Cloudflare Workers KV storage.

### `models.go`

Defines data models used in the application.

#### **Structs**

- **`Entry`**

Represents a key-value entry with associated metadata.

    ```go
    type Entry struct {
        Name     string
        Metadata Metadata
        Value    string
    }
    ```

- **`Metadata`**

Contains predefined fields for entry metadata.

    ```go
    type Metadata struct {
        Name           string `json:"name"`
        External       bool   `json:"external"`
        MimeType       string `json:"mimetype"`
        Location       string `json:"location"`
        CloudStorageID string `json:"cloud_storage_id,omitempty"`
        MD5Checksum    string `json:"md5Checksum,omitempty"`
        Description    string `json:"description,omitempty"`
    }
    ```

#### **metadata Functions**

- **Serialization Methods**

  - **`ToJSONString() (string, error)`**: Serializes the struct to a JSON string.
  - **`MetadataFromJSONString(jsonStr string) (Metadata, error)`**: Deserializes a JSON string to a `Metadata` struct.
  - **`EntryFromJSONString(jsonStr string) (Entry, error)`**: Deserializes a JSON string to an `Entry` struct.

### `session.go`

Handles interactions with the Cloudflare API.

#### **session Structs**

- **`CloudflareSession`**

Manages the session with the Cloudflare API.

    ```go
    type CloudflareSession struct {
        api          *cloudflare.API
        account_id   *cloudflare.ResourceContainer
        namespace_id string
        domain       string
    }
    ```

#### **session Functions**

- **Initialization**

  - **`NewCloudflareSession() *CloudflareSession`**: Initializes a Cloudflare session using credentials from the embedded environment variables.

- **Data Retrieval**

  - **`GetValue(key string) string`**: Retrieves the value for a specific key.
  - **`GetAllValues() []string`**: Retrieves all values from the namespace.
  - **`GetAllKeys() []cloudflare.StorageKey`**: Retrieves all keys from the namespace.
  - **`GetAllEntries() []models.Entry`**: Retrieves all entries (keys, values, and metadata).
  - **`GetAllEntriesFromKeys(storageKeys []cloudflare.StorageKey) []models.Entry`**: Retrieves entries for specific keys.
  - **`Size() (int, []cloudflare.StorageKey)`**: Returns the total number of entries and the list of keys.

- **Data Manipulation**

  - **`WriteEntry(entry models.Entry) (resp cloudflare.Response)`**: Writes a single entry to Cloudflare Workers KV.
  - **`InsertKVEntry(name string, value string, metadata string) (resp cloudflare.Response)`**: Inserts a key-value pair with metadata into Cloudflare Workers KV.
  - **`WriteEntries(entries []models.Entry)`**: Writes multiple entries to Cloudflare Workers KV.
  - **`DeleteKeyValue(key string)`**: Deletes a single key-value pair from Cloudflare Workers KV.
  - **`DeleteKeyValues(keys []string)`**: Deletes multiple key-value pairs from Cloudflare Workers KV.

- **Helper Functions**

  - **`entryToWorkersKVPairs(entry models.Entry) []*cloudflare.WorkersKVPair`**: Converts an `Entry` to a `WorkersKVPair`.
  - **`entriesToWorkersKVPairs(entries []models.Entry) []*cloudflare.WorkersKVPair`**: Converts multiple `Entry` objects.

### `database.go`

Handles interactions with the local SQLite3 database.

#### **database Structs**

- **`Database`**

Manages the connection and operations on the SQLite database.

    ```go
    type Database struct {
        dbName string
        db     *sql.DB
        lock   sync.Mutex
    }
    ```

#### **database Functions**

- **Initialization**

  - **`NewDatabase(dbName string) *Database`**: Initializes a new database connection.

- **Table Management**

  - **`CreateTable()`**: Creates the `records` table.
  - **`DropTable()`**: Drops the `records` table.

- **Data Insertion**

  - **`InsertEntry(datavalues models.Entry)`**: Inserts or replaces a single entry.
  - **`InsertKVEntryIntoDatabase(name string, value string, metadata string)`**: Inserts a key-value entry into the database with metadata.
  - **`InsertEntries(datavalues []models.Entry)`**: Inserts or replaces multiple entries in a transaction.

- **Data Retrieval**

  - **`GetEntryByName(name string) models.Entry`**: Retrieves an entry by `name`.
  - **`GetEntryByValue(value string) models.Entry`**: Retrieves an entry by `value`.
  - **`GetEntriesByValue(value string) []models.Entry`**: Retrieves entries matching a `value`.
  - **`GetAllEntries() []models.Entry`**: Retrieves all entries.

- **Data Deletion**

  - **`DeleteName(key string)`**: Deletes an entry by `name`.
  - **`DeleteNames(names []string)`**: Deletes multiple entries by `name`.
  - **`DeleteEntry(entry models.Entry)`**: Deletes a specific entry.
  - **`DeleteEntries(entries []models.Entry)`**: Deletes multiple entries.

- **Utility**

  - **`Size() int`**: Returns the total number of entries.

### **Additional Notes**

- **Synchronization**: The `SyncFromCloudflare()` function ensures that the local database mirrors the Cloudflare Workers KV storage.
- **Thread Safety**: The database operations use a mutex lock to ensure thread safety during concurrent access.

## Frontend: `main.js`

Handles the frontend logic, including user interactions, dynamic content rendering, and communication with backend Go functions.

### **Import Statements**

    ```javascript
    import './app.css';

    import { GetEntryByName, GetEntryByValue, GetEntriesByValue, GetAllEntries, InsertKVEntryIntoDatabase, DeleteName } from '../wailsjs/go/database/Database';
    import { InsertKVEntry, DeleteKeyValue } from '../wailsjs/go/session/CloudflareSession';
    ```

### **HTML Structure**

The HTML content is dynamically generated:

**Search Entry Section**: Allows users to search for entries.

    ```html
    <div class="input-box" id="search-entry">
        <label for="searchType">Search:</label>
        <select id="searchType" style="width:292px;">
            <option value="GetEntryByName">Search by UUID</option>
            <option value="GetEntryByValue">Search by URL (single)</option>
            <option value="GetEntriesByValue">Search by URL (multiple)</option>
            <option value="GetAllEntries">Get All Entries</option>
        </select>
        <input class="input" id="entryValue" type="text" autocomplete="off" placeholder="Enter search value" style="width:400px;"/>
        <button class="btn" onclick="searchEntry()">Search</button>
        <button id="clear" class="btn" onclick="clearResults()" style="display:none;">Clear</button>
    </div>
    <div class="result" id="entryResult"></div>
    ```

**Insert Entry Section**: Allows users to insert new entries manually or from a file.

    ```html
    <div class="input-box" id="insert-entry" style="margin-top:10px;">
        <label for="insertEntrySelector">Insert:</label>
        <select id="insertEntrySelector" style="width:292px;" onchange="updateInsertEntry()">
            <option value="default" selected disabled>Select Insertion Method</option>
            <option value="manual">Insert Manually</option>
            <option value="fromFile">From File</option>
            <option value="getTemplate">Download File Template</option>
        </select>
    </div>
    <div class="result" id="dynamicInsertEntry"></div>
    ```

**Delete Entry Section**: Allows users to delete entries by UUID.

    ```html
    <div class="input-box" id="delete-entry" style="margin-top:20px;">
        <label for="deleteEntryName">Delete:</label>
        <input class="input" id="deleteEntryName" type="text" placeholder="Enter UUID" size="40"/>
        <button class="btn" onclick="deleteEntry()">Delete</button>
    </div>
    ```

### **JavaScript Functions**

#### **Event Listeners**

    ```javascript
    const searchTypeElement = document.getElementById("searchType");
    const entryValueElement = document.getElementById("entryValue");
    const resultElement = document.getElementById("entryResult");
    const clearResultsButton = document.getElementById("clear");

    searchTypeElement.addEventListener('change', () => {
        if (searchTypeElement.value === "GetAllEntries") {
            entryValueElement.style.display = 'none';
            entryValueElement.value = '';
        } else {
            entryValueElement.style.display = 'inline';
        }
    });
    ```

#### **Search Functionality**

    ```javascript
    window.searchEntry = async function () {
        const value = entryValueElement.value.trim();
        const searchType = searchTypeElement.value;

        if (searchType !== "GetAllEntries" && value === "") {
            updateResults("Please enter a search value.");
            return;
        }

        try {
            let entries = [];
            switch (searchType) {
                case "GetEntryByName":
                    const entryByName = await GetEntryByName(getUUIDFromString(value));
                    if (entryByName?.Name) entries.push(entryByName);
                    break;
                case "GetEntryByValue":
                    const entryByValue = await GetEntryByValue(value);
                    if (entryByValue?.Name) entries.push(entryByValue);
                    break;
                case "GetEntriesByValue":
                    entries = await GetEntriesByValue(value);
                    break;
                case "GetAllEntries":
                    entries = await GetAllEntries();
                    break;
                default:
                    updateResults("Invalid search type.");
                    return;
            }

            if (entries.length > 0) {
                displayEntries(entries);
            } else {
                updateResults("No entries found for the provided value.");
            }
        } catch (err) {
            updateResults(`An error occurred while fetching the entries. ${err}`);
        }
    };
    ```

- **Insert Entry Functionality**

#### **Manual Insertion**

    ```javascript
    window.updateInsertEntry = function (entryMethod = undefined) {
        const selectedValue = entryMethod == undefined ? document.getElementById("insertEntrySelector").value : entryMethod;
        const dynamicInsertEntryDiv = document.getElementById("dynamicInsertEntry");
        switch (selectedValue) {
            case "manual":
                dynamicInsertEntryDiv.innerHTML = `
                <!-- HTML structure for manual entry -->
                `;
                break;
            case "fromFile":
                dynamicInsertEntryDiv.innerHTML = `
                <!-- HTML structure for file upload -->
                `;
                break;
            case "getTemplate":
                // Code to download CSV template
                break;
            default:
                // Reset the selection
        }
    };

    window.insertEntry = async function () {
        // Collect metadata and entry data
        // Validate required fields
        // Insert into Cloudflare and local database
    };
    ```

#### **Bulk Insertion from File**

    ```javascript
    window.insertEntryFromFile = async function () {
        // Read file content
        // Parse CSV data
        // Loop through entries and insert
    };
    ```

#### **Delete Entry Functionality**

    ```javascript
    window.deleteEntry = async function () {
        const uuid = document.getElementById("deleteEntryName").value.trim();
        await DeleteKeyValue(uuid);
        await DeleteName(uuid);
        clearSuccessfulDelete();
    };
    ```

#### **Metadata Management**

    ```javascript
    window.updateExternalInternalMetadataSelector = function () {
        // Show or hide additional metadata fields based on 'external' value
    };
    ```

#### **Utility Functions**

    ```javascript
    function updateResults(content = '') {
        resultElement.innerHTML = content;
        clearResultsButton.style.display = content ? 'inline' : 'none';
    }

    function displayEntries(entries) {
        // Render entries in a sortable and clickable table
        enableSorting();
        enableCopying();
        enableUUIDLinkCopying();
    }

    function enableSorting() {
        // Implement table column sorting
    }

    function enableCopying() {
        // Enable copying cell content to clipboard
    }

    function enableUUIDLinkCopying() {
        // Copy UUID as a hyperlink to the clipboard
    }

    function getUUIDFromString(stringContainingUUID) {
        // Extract UUID from a string
    }

    window.clearResults = function () {
        updateResults();
        entryValueElement.value = '';
    };

    window.clearSuccessfulDelete = function () {
        document.getElementById("deleteEntryName").value = '';
    };

    window.generateUUID = function () {
        const entryNameInput = document.getElementById('insertEntryName');
        entryNameInput.value = crypto.randomUUID();
    };
    ```

### **Features**

- **Dynamic Search Options**: Search by UUID, URL (single or multiple), or get all entries.
- **Insert Entries with Metadata**: Add new entries with predefined metadata fields.
- **Bulk Insertion from File**: Insert multiple entries at once using a CSV file template.
- **UUID Generation**: Automatically generate UUIDs for new entries.
- **Delete Entries**: Remove entries by UUID from both Cloudflare Workers KV and the local database.
- **Interactive Table**: Sortable and clickable table for displaying search results.
- **Clipboard Notifications**: On-screen notification when content is copied.
- **Dynamic Form Updates**: Forms update dynamically based on user selections.
- **Metadata Validation**: Enforces required metadata fields during entry insertion.

### **Usage**

- **Search Entries**:
  1. Select the search type.
  2. Enter the search value if required.
  3. Click "Search" to display results.
  4. Click "Clear" to reset.

- **Insert Entries Manually**:
  1. Select "Insert Manually" from the insertion method dropdown.
  2. Enter "Name" and "Value". Use the UUID generator if needed.
  3. Fill in the required metadata fields.
  4. Click "Insert" to add the entry.

- **Insert Entries from File**:
  1. Select "From File" from the insertion method dropdown.
  2. Choose the CSV file using the provided template.
  3. Click "Insert" to add entries.

- **Delete Entries**:
  1. Enter the UUID of the entry to delete.
  2. Click "Delete" to remove the entry from both Cloudflare and the local database.

- **Copy Data**:
  - Click on any cell to copy its content to the clipboard.
  - Click on a UUID to copy its hyperlink.

## Notes

- **Backend Integration**: Go functions are exposed via Wails bindings for frontend interaction.
- **Environment Configuration**: The `.env` file must be configured with valid Cloudflare credentials and will be embedded during the build process.
- **Error Handling**: The application includes error handling for database operations and clipboard interactions.
- **Styling**: Styles are defined in `app.css` to match UI requirements.
- **Thread Safety**: Database operations are thread-safe, ensuring data integrity during concurrent access.

**Note**: Ensure that the `.env` file is included in the Docker build context and properly configured before building the Docker image.
