import './style.css';
import './app.css';

import { GetEntryByName, GetEntryByValue, GetEntriesByValue, GetAllEntries } from '../wailsjs/go/database/Database';

document.querySelector('#app').innerHTML = `
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
`;

document.querySelector('#app').innerHTML += `
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


window.clearResults = function () {
    updateResults();
    entryValueElement.value = '';
};
