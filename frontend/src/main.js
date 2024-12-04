import './app.css';

import { GetEntryByName, GetEntryByValue, GetEntriesByValue, GetAllEntries, InsertKVEntryIntoDatabase, DeleteName } from '../wailsjs/go/database/Database';
import { InsertKVEntry, DeleteKeyValue } from '../wailsjs/go/session/CloudflareSession';

document.querySelector('#app').innerHTML = `
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
`;

document.querySelector('#app').innerHTML += `
    <div class="input-box" id="insert-entry" style="margin-top:10px;">
        <label for="insertEntrySelector">Insert:</label>
        <select id="insertEntrySelector" style="width:292px;" onchange="updateInsertEntry()">
            <option value="default" selected disabled>Select Insertion Method</option>
            <option value="manual">Insert Manually</option>
            <option value="fromFile">From File</option>
            <option value="getTemplate">Download File Template</option>
        </select>
    </div>
    <div  class="result" id="dynamicInsertEntry"></div>
`;

window.updateExternalInternalMetadataSelector = function () {
    const selectedValue = document.getElementById("externalMetadataToggle").value;
    const cloudStorageDiv = document.getElementById("cloud-storage-id-div");
    const md5ChecksumDiv =  document.getElementById("md5checksum-div");
    switch(selectedValue) {
        case "true":
            cloudStorageDiv.style.display = 'none';
            md5ChecksumDiv.style.display = 'none';            
            break;
        case "false":
        default:
            cloudStorageDiv.style.display = 'block';
            md5ChecksumDiv.style.display = 'block';
    }
}

window.updateInsertEntry = function (entryMethod = undefined) {
    const selectedValue = entryMethod == undefined ? document.getElementById("insertEntrySelector").value : entryMethod
    const dynamicInsertEntryDiv = document.getElementById("dynamicInsertEntry");
    switch (selectedValue) {
        case "manual":
            dynamicInsertEntryDiv.innerHTML = `
            <div class="input-box" id="manual-insert-entry" style="margin-top:10px;margin-left:75px;">
                <div style="position: relative; display: inline-block;">
                    <input class="input" id="insertEntryName" type="text" placeholder="Enter name" size="40"/>
                    <svg 
                        onclick="generateUUID()" 
                        xmlns="http://www.w3.org/2000/svg" 
                        width="20" 
                        height="20" 
                        viewBox="0 0 24 24" 
                        fill="none" 
                        stroke="#5007b5" 
                        stroke-width="2" 
                        stroke-linecap="round" 
                        stroke-linejoin="round" 
                        style="position: absolute; top: 50%; right: 10px; transform: translateY(-50%); cursor: pointer;">
                        <circle cx="12" cy="12" r="10"></circle>
                        <line x1="12" y1="8" x2="12" y2="16"></line>
                        <line x1="8" y1="12" x2="16" y2="12"></line>
                    </svg>
                </div>
                <input class="input" id="insertEntryValue" type="text" placeholder="Enter value" style="width:400px;"/>
                <button class="btn" onclick="insertEntry()">Insert</button>
                <div id="entryMetadata">
                    <div class="indented metadata-entry">
                        <input class="input jsonKey" type="text" value="metadata_name" readonly style="margin-right: 5px;">
                        <input class="input jsonValue" type="text" placeholder="Resource Title" required>
                    </div>
                    <div class="indented metadata-entry">
                        <input class="input jsonKey" type="text" value="metadata_external" readonly style="margin-right: 5px;">
                        <select  class="input jsonValue" id="externalMetadataToggle" style="width:422px;" required  onchange="updateExternalInternalMetadataSelector()">
                            <option value="default" selected disabled>Resource Is External</option>
                            <option value="true">True</option>
                            <option value="false">False</option>
                        </select>
                    </div>
                    <div class="indented metadata-entry">
                        <input class="input jsonKey" type="text" value="metadata_mimetype" readonly style="margin-right: 5px;">
                        <input class="input jsonValue" type="text" placeholder="Resource MimeType" required>
                    </div>
                    <div class="indented metadata-entry">
                        <input class="input jsonKey" type="text" value="metadata_location" readonly style="margin-right: 5px;">
                        <input class="input jsonValue" type="text" placeholder="Resource Location (domain or owner email)" required>
                    </div>
                    <div class="indented metadata-entry">
                        <input class="input jsonKey" type="text" value="metadata_description" readonly style="margin-right: 5px;">
                        <input class="input jsonValue" type="text" placeholder="Resource Description">
                    </div>
                    <div id="cloud-storage-id-div" class="indented metadata-entry" style="display:none;">
                        <input class="input jsonKey" type="text" value="metadata_cloud_storage_id" readonly style="margin-right:-5px;">
                        <input class="input jsonValue" type="text" placeholder="Resource Cloud Storage ID">
                    </div>
                    <div  id="md5checksum-div" class="indented metadata-entry" style="display:none;">
                        <input class="input jsonKey" type="text" value="metadata_md5Checksum" readonly style="margin-right:-5px;">
                        <input class="input jsonValue" type="text" placeholder="Resource MD5 Checksum">
                    </div>
                </div>
            </div>
        `;
            break;
        case "fromFile":
            dynamicInsertEntryDiv.innerHTML = `
            <div class="input-box" id="file-insert-entry" style="margin-top:10px;margin-left:75px;">
                <input class="input" id="insertFile" type="file" accept=".csv" style="border:0px;background-color:transparent;"/>
                <button class="btn" onclick="insertEntryFromFile()">Insert</button>
            </div>
        `;
            break;
        case "getTemplate":
            const csvContent = `name,value,metadata_name,metadata_external,metadata_mimetype,metadata_location,metadata_description,metadata_cloud_storage_id,metadata_md5Checksum`;

            const blob = new Blob([csvContent], { type: "text/csv" });
            const a = document.createElement("a");
            a.href = URL.createObjectURL(blob);
            a.download = "CDN Manager Bulk Insert Template.csv";
            a.style.display = "none";
            document.body.appendChild(a);
            a.click();
            document.body.removeChild(a);
        default:
            document.getElementById("insertEntrySelector").value = "default";
            dynamicInsertEntryDiv.innerHTML = `
              <div class="result"></div>
            `;
    }
};

document.querySelector('#app').innerHTML += `
    <div class="input-box" id="delete-entry" style="margin-top:10px;">
        <label for="deleteEntryName">Delete:</label>
        <input class="input" id="deleteEntryName" type="text" placeholder="Enter UUID" size="40"/>
        <button class="btn" onclick="deleteEntry()">Delete</button>
    </div>
`;

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
        console.error(err);
        updateResults("An error occurred while fetching the entries.");
    }
};

window.removeMetaDataEntryField = function () {
    const entryMetadata = document.getElementById("entryMetadata");
    entryMetadata.removeChild(entryMetadata.lastChild);
};

window.deleteEntry = async function () {
    const uuid = document.getElementById("deleteEntryName").value.trim();
    await DeleteKeyValue(uuid);
    await DeleteName(uuid)
    clearSuccessfulDelete();
}

window.clearResults = function () {
    updateResults();
    entryValueElement.value = '';
};

window.clearSuccessfulDelete = function () {
    document.getElementById("deleteEntryName").value = ''
}

window.generateUUID = function () {
    const entryNameInput = document.getElementById('insertEntryName');
    entryNameInput.value = crypto.randomUUID();
}

window.insertEntry = async function () {
    const metadataEntries = document.querySelectorAll('.metadata-entry');
    const metadata = {};

    metadataEntries.forEach(entry => {
        const keyInput = entry.querySelector('.jsonKey');
        const valueInput = entry.querySelector('.jsonValue');

        const key = keyInput.value.trim();
        let value;

        if (valueInput.tagName.toLowerCase() === 'select') {
            value = valueInput.options[valueInput.selectedIndex].value;
        } else {
            value = valueInput.value.trim();
        }

        if (key && value && value !== 'default') {
            metadata[key] = value;
        }
    });

    const value = document.getElementById("insertEntryValue").value.trim();
    const name = document.getElementById("insertEntryName").value.trim();

    // Validate required fields
    if (!name || !value) {
        alert("Please provide both Name and Value.");
        return;
    }

    // Check if 'metadata_external' is selected
    if (metadata['metadata_external'] === undefined || metadata['metadata_external'] === 'default') {
        alert("Please select whether the resource is external.");
        return;
    }

    try {
        // InsertKVEntry to add to cloudflare
        // InsertKVEntryIntoDatabase to add to local database
        const metadataString = JSON.stringify(metadata)
        const response = await InsertKVEntry(name, value, metadataString);
        console.log(response);
        if (response && response.success) {
            await InsertKVEntryIntoDatabase(name, value, metadataString);
            updateInsertEntry("");
            console.log('Entry successfully inserted into the database.');
        } else {
            console.error('Failed to insert entry:', response.errors);
            alert('Failed to insert entry: ' + response.errors.join(', '));
        }
    } catch (error) {
        console.error('Failed to insert entry:', error);
        alert('An error occurred while inserting the entry.');
    }
};


window.insertEntryFromFile = async function () {
    console.log("stub function insertEntryFromFile")
}

function updateResults(content = '') {
    resultElement.innerHTML = content;
    clearResultsButton.style.display = content ? 'inline' : 'none';
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

function displayEntries(entries) {
    let tableHTML = `
        <table id="resultTable" style="margin-bottom:10px;">
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
        const metadataObject = entry.Metadata;
        const metadataFormatted = Object.entries(metadataObject)
            .map(([key, value]) => `${key}: ${value}`)
            .join('\n');

        tableHTML += `
            <tr>
                <td class="hyperlink">${entry.Name}</td>
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
    enableUUIDLinkCopying();
}

function enableSorting() {
    const table = document.getElementById("resultTable");
    const headers = table.querySelectorAll(".sortable");
    let sortDirection = 1;

    headers.forEach(header => {
        header.addEventListener("click", () => {
            const column = header.dataset.column;
            const rows = Array.from(table.querySelector("tbody").rows);

            rows.sort((a, b) => {
                const aText = a.querySelector(`td:nth-child(${Array.from(headers).indexOf(header) + 1})`).textContent.trim();
                const bText = b.querySelector(`td:nth-child(${Array.from(headers).indexOf(header) + 1})`).textContent.trim();

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

function enableUUIDLinkCopying() {
    document.querySelectorAll('.hyperlink').forEach(td => {
        td.addEventListener('click', () => {
            const hyperlink = `https://cdn.williamveith.com/?id=${td.textContent.trim()}`
            navigator.clipboard.writeText(hyperlink).then(() => {
                displayClipboardMessage(`Copied: ${hyperlink}`);
            }).catch(err => {
                console.error('Error copying to clipboard:', err);
            });
        });
    });
}

function getUUIDFromString(stringContainingUUID) {
    const uuidPattern = /\b[0-9a-fA-F]{8}(?:-[0-9a-fA-F]{4}){3}-[0-9a-fA-F]{12}\b/;
    const match = stringContainingUUID.match(uuidPattern);
    return match ? match[0] : '';
}