import './style.css';
import './app.css';

import { GetEntryByName, GetEntryByValue, GetEntriesByValue, GetAllEntries, InsertKVEntryIntoDatabase } from '../wailsjs/go/database/Database';
import { InsertKVEntry } from '../wailsjs/go/session/CloudflareSession';

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
        <label for="insertEntryName">Insert:</label>
        <div style="position: relative; display: inline-block;">
            <input class="input" id="insertEntryName" type="text" placeholder="Enter name" size="40"/>
            <!-- Add SVG icon -->
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
        <div id="entryMetadata"></div>
        <span class="indent">
            <button class="btn" onclick="addMetaDataEntryField()" style="width:auto;margin-top:10px;margin-left:80px;">+ MetaData</button>
        </span>
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


function updateResults(content = '') {
    resultElement.innerHTML = content;
    clearResultsButton.style.display = content ? 'inline' : 'none';
}

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

window.addMetaDataEntryField = function () {
    const entryMetadataDiv = document.getElementById('entryMetadata');

    const newEntryDiv = document.createElement('div');
    newEntryDiv.className = 'indented metadata-entry';
    newEntryDiv.style.display = 'flex';
    newEntryDiv.style.alignItems = 'center';
    newEntryDiv.style.marginBottom = '5px';

    const newKeyInput = document.createElement('input');
    newKeyInput.className = 'input jsonKey';
    newKeyInput.type = 'text';
    newKeyInput.placeholder = 'Metadata Key';
    newKeyInput.required = true;
    newKeyInput.style.marginRight = '5px';

    const newValueInput = document.createElement('input');
    newValueInput.className = 'input jsonValue';
    newValueInput.type = 'text';
    newValueInput.placeholder = 'Metadata Value';
    newKeyInput.required = true;

    const newRemoveButton = document.createElement('button')
    newRemoveButton.className = 'btn'
    newRemoveButton.style = 'width:auto;'
    newRemoveButton.innerHTML = 'Remove'
    newRemoveButton.addEventListener("click", function () {
        const jsonDiv = newRemoveButton.parentNode
        const insertDiv = jsonDiv.parentNode
        insertDiv.removeChild(jsonDiv);
    });

    newEntryDiv.appendChild(newKeyInput);
    newEntryDiv.appendChild(newValueInput);
    newEntryDiv.appendChild(newRemoveButton)

    entryMetadataDiv.appendChild(newEntryDiv);
};

window.removeMetaDataEntryField = function () {
    const entryMetadata = document.getElementById("entryMetadata");
    entryMetadata.removeChild(entryMetadata.lastChild);
};

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
        const metadataObject = JSON.parse(entry.Metadata);
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

window.generateUUID = function () {
    const entryNameInput = document.getElementById('insertEntryName');
    entryNameInput.value = crypto.randomUUID();
}

window.insertEntry = async function () {
    const metadataEntries = document.querySelectorAll('.metadata-entry');
    const jsonObject = {};

    metadataEntries.forEach(entry => {
        const key = entry.querySelector('.jsonKey').value.trim();
        const value = entry.querySelector('.jsonValue').value.trim();
        if (key && value) {
            jsonObject[key] = value;
        }
    });

    const metadata = jsonObject;
    const value = document.getElementById("insertEntryValue").value.trim();
    const name = document.getElementById("insertEntryName").value.trim();

    try {
        const response = await InsertKVEntry(name, value, metadata);
        console.log(response);
        if (response && response.success) {
            await InsertKVEntryIntoDatabase(name, value, metadata);
            clearSuccessfulInputs()
            console.log('Entry successfully inserted into the database.');
        } else {
            console.error('Failed to insert entry into Cloudflare KV:', response.errors);
        }
    } catch (error) {
        console.error('Failed to insert entry:', error);
    }
}


window.clearResults = function () {
    updateResults();
    entryValueElement.value = '';
};

window.clearSuccessfulInputs = function () {
    document.getElementById("insertEntryValue").value = '';
    document.getElementById("insertEntryName").value = '';
    const metaDataDiv = document.getElementById("entryMetadata")
    while (metaDataDiv.firstChild) {
        metaDataDiv.removeChild(metaDataDiv.firstChild);
    }
}
