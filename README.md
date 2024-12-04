# Cloudflare CDN Manager

## About

**Cloudflare CDN Manager** is a cross-platform application designed to manage a Cloudflare KV/WebWorker redirect service with an intuitive user interface. The program is capable of running on Windows, macOS, Linux, and can also be deployed via Docker for web deployments. It leverages a Wails UI with Go for making Cloudflare V4 API calls and controls a local SQLite3 database to store and manage redirect entries.

## Future Plans

- **Content Replication & Tracking**: Auto replicate content across multiple cloud storage providers, track those copies in the KV space, and have a webworker maintain the list by removing resources if they become inaccessible.
- **Wrapping Content Instead of Redirecting**: Currently this is a URL redirect service. I currently use it for QR Codes, to make them dynamic. However, the content served to users is just the raw resource. Building another webworker that wraps the content in a MimeType appropriate manner would allow the content to be served in a more dynamic way (Render HTML document content, video content, etc) without being limited to the particular cloud storage implementations.

## Features

- **Cross-Platform Compatibility**: Runs on Windows, macOS, Linux, and supports Docker deployments.
- **Cloudflare KV Management**: Interacts with Cloudflare Workers KV storage to manage key-value pairs for redirects.
- **User-Friendly UI**: Provides a graphical interface for managing redirects, entries, and metadata.
- **Local Database Control**: Utilizes SQLite3 for local storage and management of redirect entries.
- **Dynamic Search Options**: Search entries by UUID, URL (single or multiple), or retrieve all entries.
- **Metadata Support**: Allows adding custom metadata to entries.
- **Interactive Data Handling**: Sortable tables, clickable cells for easy copying, and dynamic forms.

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

Replace the empty strings with your actual Cloudflare account details.

## Backend: Go Modules

### `database.go`

Handles interactions with the local SQLite3 database.

#### **Structs - Database**

- **`Database`**

  Manages the connection and operations on the SQLite database.

  ```go
  type Database struct {
      dbName string
      db     *sql.DB
      lock   sync.Mutex
  }
  ```

#### **Functions**

- **Initialization**

  - **`NewDatabase(dbName string) *Database`**: Initializes a new database connection.

- **Table Management**

  - **`CreateTable()`**: Creates the `records` table.
  - **`DropTable()`**: Drops the `records` table.

- **Data Insertion**

  - **`InsertEntry(datavalues session.Entry)`**: Inserts or replaces a single entry.
  - **`InsertEntries(datavalues []session.Entry)`**: Inserts or replaces multiple entries in a transaction.

- **Data Retrieval**

  - **`GetEntryByName(name string) session.Entry`**: Retrieves an entry by `name`.
  - **`GetEntryByValue(value string) session.Entry`**: Retrieves an entry by `value`.
  - **`GetEntriesByValue(value string) []session.Entry`**: Retrieves entries matching a `value`.
  - **`GetAllEntries() []session.Entry`**: Retrieves all entries.

- **Data Deletion**

  - **`DeleteName(key string)`**: Deletes an entry by `name`.
  - **`DeleteNames(names []string)`**: Deletes multiple entries by `name`.
  - **`DeleteEntry(entry session.Entry)`**: Deletes a specific entry.
  - **`DeleteEntries(entries []session.Entry)`**: Deletes multiple entries.

- **Utility**

  - **`Size() int`**: Returns the total number of entries.

#### **Helper Functions**

- **`convertMetadataToString(metadata interface{}) string`**: Converts metadata to a JSON string for storage.

### `session.go`

Handles interactions with the Cloudflare API.

#### **Structs**

- **`Entry`**

  Represents a key-value entry with optional metadata.

  ```go
  type Entry struct {
      Name     string
      Metadata interface{}
      Value    string
  }
  ```

- **`CloudflareSession`**

  Manages the session with Cloudflare API.

  ```go
  type CloudflareSession struct {
      api          *cloudflare.API
      account_id   *cloudflare.ResourceContainer
      namespace_id string
      domain       string
  }
  ```

#### **Functions - Session**

- **Initialization**

  - **`NewCloudflareSession() *CloudflareSession`**: Initializes a Cloudflare session using credentials from the `.env` file.

- **Data Retrieval**

  - **`GetValue(key string) string`**: Retrieves the value for a specific key.
  - **`GetAllValues() []string`**: Retrieves all values from the namespace.
  - **`GetAllKeys() []cloudflare.StorageKey`**: Retrieves all keys from the namespace.
  - **`GetAllEntries() []Entry`**: Retrieves all entries (keys and values).
  - **`GetAllEntriesFromKeys(storageKeys []cloudflare.StorageKey) []Entry`**: Retrieves entries for specific keys.
  - **`Size() (int, []cloudflare.StorageKey)`**: Returns the total number of entries and the list of keys.

- **Data Manipulation**

  - **`WriteEntry(entry Entry)`**: Writes a single entry to Cloudflare Workers KV.
  - **`WriteEntries(entries []Entry)`**: Writes multiple entries.
  - **`DeleteKeyValue(key string)`**: Deletes a single key-value pair.
  - **`DeleteKeyValues(keys []string)`**: Deletes multiple key-value pairs.

#### **Helper Functions - Session**

- **`entryToWorkersKVPairs(entry Entry) []*cloudflare.WorkersKVPair`**: Converts an `Entry` to a `WorkersKVPair`.
- **`entriesToWorkersKVPairs(entries []Entry) []*cloudflare.WorkersKVPair`**: Converts multiple `Entry` objects.

## Frontend: `main.js`

Handles the frontend logic, including user interactions, dynamic content rendering, and communication with backend Go functions.

### **Import Statements**

```javascript
import './app.css';

import { GetEntryByName, GetEntryByValue, GetEntriesByValue, GetAllEntries } from '../wailsjs/go/database/Database';
```

### **HTML Structure**

The HTML content is dynamically generated:

- **Search Entry Section**: Allows users to search for entries.

  ```html
  <div class="input-box" id="search-entry">
      <label for="searchType">Search:</label>
      <select id="searchType">
          <option value="GetEntryByName">Search by UUID</option>
          <option value="GetEntryByValue">Search by URL (single)</option>
          <option value="GetEntriesByValue">Search by URL (multiple)</option>
          <option value="GetAllEntries">Get All Entries</option>
      </select>
      <input class="input" id="entryValue" type="text" autocomplete="off" placeholder="Enter search value" />
      <button class="btn" onclick="searchEntry()">Search</button>
      <button id="clear" class="btn" onclick="clearResults()" style="display:none;">Clear</button>
  </div>
  <div class="result" id="entryResult"></div>
  ```

- **Insert Entry Section**: Allows users to insert new entries.

  ```html
  <div class="input-box" id="insert-entry">
      <label for="entryName">Insert:</label>
      <input class="input" id="entryName" type="text" placeholder="Enter name" />
      <input class="input" id="entryValue" type="text" placeholder="Enter value" />
      <button class="btn" onclick="insertEntry()">Insert</button>
      <div id="entryMetadata"></div>
      <span class="indent">
          <button class="btn" onclick="addMetaDataEntryField()" style="width:auto;margin-top:10px;margin-left:250px;">+ MetaData</button>
          <button class="btn" onclick="removeMetaDataEntryField()" style="width:auto;">- MetaData</button>
      </span>
  </div>
  ```

### **JavaScript Functions**

- **Event Listeners**

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

- **Search Functionality**

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
                  const entryByName = await GetEntryByName(value);
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
          console.error(err);
          updateResults("An error occurred while fetching the entries.");
      }
  };
  ```

- **Insert Entry Functionality**

  *(Note: Implement the `insertEntry` function to handle the insertion of new entries.)*

- **Metadata Management**

  ```javascript
  window.addMetaDataEntryField = function () {
      const entryMetadataDiv = document.getElementById('entryMetadata');

      const newEntryDiv = document.createElement('div');
      newEntryDiv.className = 'indented';
      newEntryDiv.style.display = 'flex';
      newEntryDiv.style.alignItems = 'center';
      newEntryDiv.style.marginBottom = '5px';

      const newKeyInput = document.createElement('input');
      newKeyInput.className = 'input jsonKey';
      newKeyInput.type = 'text';
      newKeyInput.placeholder = 'Enter JSON Key';
      newKeyInput.style.marginRight = '5px';

      const newValueInput = document.createElement('input');
      newValueInput.className = 'input jsonValue';
      newValueInput.type = 'text';
      newValueInput.placeholder = 'Enter JSON Value';

      newEntryDiv.appendChild(newKeyInput);
      newEntryDiv.appendChild(newValueInput);

      entryMetadataDiv.appendChild(newEntryDiv);
  };

  window.removeMetaDataEntryField = function(){
      const entryMetadata = document.getElementById("entryMetadata");
      if (entryMetadata.lastChild) {
          entryMetadata.removeChild(entryMetadata.lastChild);
      }
  };
  ```

- **Utility Functions**

  ```javascript
  function updateResults(content = '') {
      resultElement.innerHTML = content;
      clearResultsButton.style.display = content ? 'inline' : 'none';
  }

  function displayEntries(entries) {
      let tableHTML = `
          <table id="resultTable">
              <thead>
                  <tr>
                      <th data-column="Name" class="sortable">Name</th>
                      <th data-column="Value" class="sortable">Value</th>
                      <th data-column="Metadata">Metadata</th>
                  </tr>
              </thead>
              <tbody>
      `;

      entries.forEach(entry => {
          const metadataObject = JSON.parse(entry.Metadata || '{}');
          const metadataFormatted = Object.entries(metadataObject)
              .map(([key, value]) => `${key}: ${value}`)
              .join('\n');

          tableHTML += `
              <tr>
                  <td class="clickable">${entry.Name}</td>
                  <td class="clickable">${entry.Value}</td>
                  <td class="clickable"><pre>${metadataFormatted}</pre></td>
              </tr>
          `;
      });

      tableHTML += `
              </tbody>
          </table>
      `;

      updateResults(tableHTML);

      enableSorting();
      enableCopying();
  }

  function enableSorting() {
      const table = document.getElementById("resultTable");
      const headers = table.querySelectorAll(".sortable");
      let sortDirection = 1;

      headers.forEach(header => {
          header.addEventListener("click", () => {
              const columnIndex = Array.from(headers).indexOf(header);
              const rows = Array.from(table.querySelector("tbody").rows);

              rows.sort((a, b) => {
                  const aText = a.cells[columnIndex].textContent.trim();
                  const bText = b.cells[columnIndex].textContent.trim();

                  if (!isNaN(aText) && !isNaN(bText)) {
                      return sortDirection * (parseFloat(aText) - parseFloat(bText));
                  }

                  return sortDirection * aText.localeCompare(bText);
              });

              rows.forEach(row => table.querySelector("tbody").appendChild(row));
              sortDirection *= -1;
          });
      });
  }

  function enableCopying() {
      document.querySelectorAll('.clickable').forEach(td => {
          td.addEventListener('click', () => {
              navigator.clipboard.writeText(td.textContent.trim()).then(() => {
                  displayClipboardMessage(`Copied: ${td.textContent.trim()}`);
              }).catch(err => {
                  console.error('Error copying to clipboard:', err);
              });
          });
      });
  }

  function displayClipboardMessage(message) {
      const messageElement = document.createElement('div');
      messageElement.innerText = message;
      messageElement.style.position = 'fixed';
      messageElement.style.bottom = '10px';
      messageElement.style.right = '10px';
      messageElement.style.backgroundColor = '#333';
      messageElement.style.color = '#fff';
      messageElement.style.padding = '10px';
      messageElement.style.borderRadius = '5px';
      messageElement.style.zIndex = 1000;

      document.body.appendChild(messageElement);
      setTimeout(() => {
          document.body.removeChild(messageElement);
      }, 2000);
  }

  window.clearResults = function () {
      updateResults();
      entryValueElement.value = '';
  };
  ```

### **Features**

- **Dynamic Search Options**: Search by UUID, URL (single or multiple), or get all entries.
- **Insert Entries with Metadata**: Add new entries with optional metadata fields.
- **Interactive Table**: Sortable and clickable table for displaying search results.
- **Clipboard Notifications**: On-screen notification when content is copied.
- **Dynamic Metadata Fields**: Add or remove metadata fields dynamically.

### **Usage**

- **Search Entries**:
  1. Select the search type.
  2. Enter the search value if required.
  3. Click "Search" to display results.
  4. Click "Clear" to reset.

- **Insert Entries**:
  1. Enter "Name" and "Value".
  2. Optionally add metadata fields.
  3. Click "Insert" to add the entry.

- **Copy Data**:
  - Click on any cell to copy its content to the clipboard.

## Notes

- **Backend Integration**: Ensure that Go functions are properly implemented and exposed via Wails bindings.
- **Implement `insertEntry` Function**: The `insertEntry` function in `main.js` needs to be implemented to handle entry insertion.
- **Environment Configuration**: The `.env` file must be configured with valid Cloudflare credentials.
- **Error Handling**: Implement error handling for database and clipboard interactions.
- **Styling**: Define styles in `app.css` to match UI requirements.

## Docker Deployment

To deploy the application using Docker:

1. Build the Docker image:

   ```bash
   docker build -t cloudflare-cdn-manager .
   ```

2. Run the Docker container:

   ```bash
   docker run -p 8080:8080 cloudflare-cdn-manager
   ```
