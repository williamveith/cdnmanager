import './styles/app.css';

import {
    IsConfigured,
    SetupAndSync,
    SyncFromCloudflare,
    InsertKVEntry,
    DeleteKeyValue,
    GenerateCSV,
    ShowAlert,
    GetDomain
} from '../wailsjs/go/main/App';

import {
    GetEntryByName,
    GetEntryByValue,
    GetEntriesByValue,
    GetAllEntries,
    InsertKVEntryIntoDatabase,
    DeleteName
} from '../wailsjs/go/database/Database';

import Fuse from 'fuse.js';
import Papa from "papaparse";

let fuse;
let appDomain = "";
window.cachedEntries = [];

const appRoot = document.querySelector('#app');

const fuseOptions = {
    keys: ['Metadata.name', 'Metadata.mimetype', 'Metadata.location', 'Metadata.description'],
    threshold: 0.3,
    includeScore: true
};

function initializeFuse(data) {
    fuse = new Fuse(data, fuseOptions);
}

window.addEventListener("DOMContentLoaded", async () => {
    await initializeApp();
});

function normalizeDomain(domain) {
    const trimmed = (domain ?? '').trim();
    if (!trimmed) return '';

    if (/^https?:\/\//i.test(trimmed)) {
        return trimmed.replace(/\/+$/, '');
    }

    return `https://${trimmed.replace(/\/+$/, '')}`;
}

function buildEntryLink(id) {
    if (!appDomain) {
        // refresh domain asynchronously but don't block UI
        GetDomain().then(domain => {
            appDomain = normalizeDomain(domain);
        }).catch(() => { });
        return `?id=${id}`;
    }

    return `${appDomain}/?id=${id}`;
}

async function initializeApp() {
    try {
        const configured = await IsConfigured();

        if (!configured) {
            renderConfigForm();
            return;
        }

        appDomain = normalizeDomain(await GetDomain());

        await SyncFromCloudflare();
        renderMainApp();
    } catch (err) {
        renderConfigForm();
        ShowAlert(`Startup failed: ${err}`);
    }
}

function renderConfigForm() {
    appRoot.innerHTML = `
    <form id="config-form">
        <div class="section" id="config-form-section">
            <div style="font-size:24px;font-weight:bold;margin-bottom:10px;">CDN Manager Setup</div>
            <div style="margin-bottom:16px;">Enter your Cloudflare configuration to initialize the application.</div>

            <div class="section">
                <input class="input"
                    id="config-cloudflare-api-token"
                    type="password"
                    required
                    pattern="[A-Za-z0-9_-]{30,}"
                    title="Cloudflare API Token must contain only letters, numbers, underscores, or hyphens"
                    spellcheck="false"
                    placeholder="Cloudflare API Token"
                    style="width:500px;" />
            </div>

            <div class="section">
                <input class="input"
                    id="config-account-id"
                    type="text"
                    required
                    pattern="[a-f0-9]{32}"
                    title="Account ID must be 32 lowercase hexadecimal characters"
                    spellcheck="false"
                    placeholder="Account ID"
                    style="width:500px;" />
            </div>

            <div class="section">
                <input class="input"
                    id="config-namespace-id"
                    type="text"
                    required
                    pattern="[a-f0-9]{32}"
                    title="Namespace ID must be 32 lowercase hexadecimal characters"
                    spellcheck="false"
                    placeholder="Namespace ID"
                    style="width:500px;" />
            </div>

            <div class="section">
                <input class="input"
                    id="config-domain"
                    type="text"
                    required
                    pattern="(https?:\\/\\/)?([a-zA-Z0-9-]+\\.)+[a-zA-Z]{2,}"
                    title="Enter a valid domain such as cdn.example.com"
                    spellcheck="false"
                    placeholder="Domain"
                    style="width:500px;" />
            </div>

            <div class="section" style="width:auto">
                <button class="btn" id="save-config-button" type="submit">Save & Sync</button>
            </div>

            <div class="result section" id="config-status"></div>
        </div>
    </form>
    `;

    document.getElementById("config-form").addEventListener("submit", function (e) {
        e.preventDefault();
        submitConfigForm();
    });
}

async function submitConfigForm() {
    const cfg = {
        cloudflare_api_token: document.getElementById("config-cloudflare-api-token").value.trim(),
        account_id: document.getElementById("config-account-id").value.trim(),
        namespace_id: document.getElementById("config-namespace-id").value.trim(),
        domain: document.getElementById("config-domain").value.trim()
    };

    if (
        !cfg.cloudflare_api_token ||
        !cfg.account_id ||
        !cfg.namespace_id ||
        !cfg.domain
    ) {
        ShowAlert("All configuration fields are required.");
        return;
    }

    try {
        document.getElementById("config-status").innerHTML = "Saving configuration and syncing Cloudflare data...";
        await SetupAndSync(cfg);
        appDomain = normalizeDomain(await GetDomain());
        renderMainApp();
        ShowAlert("Configuration saved and database synced.");
    } catch (err) {
        document.getElementById("config-status").innerHTML = "";
        ShowAlert(`Failed to save configuration or sync database. ${err}`);
    }
}

function renderMainApp() {
    appRoot.innerHTML = `
        <div id="search-entry" class="section">
            <label for="searchType">Search:</label>
            <select id="searchType" style="width:292px;">
                <option value="GetAllEntries">All</option>
                <option value="GetEntryByName">By UUID</option>
                <option value="GetEntryByValue">By URL (single)</option>
                <option value="GetEntriesByValue">By URL (multiple)</option>
            </select>
            <input class="input" id="entryValue" type="text" spellcheck="false" autocomplete="off" placeholder="Enter search value" style="width:400px;display:none;"/>
            <button class="btn" id="search-button">Search</button>
            <button id="clear" class="btn" style="display:none;">Clear</button>
        </div>
        <div class="result section" id="entryResult"></div>

        <div id="insert-entry" class="section">
            <label for="insertEntrySelector">Insert:</label>
            <select id="insertEntrySelector" style="width:292px;">
                <option value="default" selected disabled>Select Insertion Method</option>
                <option value="manual">Insert Manually</option>
                <option value="fromFile">From File</option>
                <option value="getBulkInsertTemplate">Download File Template</option>
            </select>
            <button id="clear-insert" class="btn" style="display:none;">Clear</button>
        </div>
        <div class="result" id="dynamicInsertEntry"></div>

        <div id="delete-entry" class="section">
            <label for="deleteEntryName">Delete:</label>
            <input class="input" id="deleteEntryName" type="text" spellcheck="false" placeholder="Enter UUID" size="40"/>
            <button class="btn" id="delete-button">Delete</button>
        </div>
    `;

    const searchTypeElement = document.getElementById("searchType");
    const entryValueElement = document.getElementById("entryValue");
    const clearResultsButton = document.getElementById("clear");
    const insertEntrySelector = document.getElementById("insertEntrySelector");
    const clearInsertButton = document.getElementById("clear-insert");
    const deleteEntryName = document.getElementById("deleteEntryName");

    searchTypeElement.addEventListener('change', () => {
        if (searchTypeElement.value === "GetAllEntries") {
            entryValueElement.style.display = 'none';
            entryValueElement.value = '';
        } else {
            entryValueElement.style.display = 'inline';
        }
    });

    document.getElementById("search-button").addEventListener("click", searchEntry);
    entryValueElement.addEventListener("keydown", searchEntry);
    clearResultsButton.addEventListener("click", clearResults);

    insertEntrySelector.addEventListener("change", () => updateInsertEntry());
    clearInsertButton.addEventListener("click", () => updateInsertEntry(""));

    document.getElementById("delete-button").addEventListener("click", deleteEntry);
    deleteEntryName.addEventListener("keydown", deleteEntry);
}

window.updateExternalInternalMetadataSelector = function () {
    const selectedValue = document.getElementById("externalMetadataToggle").value;
    const cloudStorageDiv = document.getElementById("cloud-storage-id-div");
    const md5ChecksumDiv = document.getElementById("md5checksum-div");

    switch (selectedValue) {
        case "true":
            cloudStorageDiv.style.display = 'none';
            md5ChecksumDiv.style.display = 'none';
            break;
        case "false":
        default:
            cloudStorageDiv.style.display = 'block';
            md5ChecksumDiv.style.display = 'block';
    }
};

window.updateInsertEntry = function (entryMethod = undefined) {
    const selectedValue = entryMethod === undefined
        ? document.getElementById("insertEntrySelector").value
        : entryMethod;

    const dynamicInsertEntryDiv = document.getElementById("dynamicInsertEntry");

    switch (selectedValue) {
        case "manual":
            dynamicInsertEntryDiv.innerHTML = `
                <div class="section" id="manual-insert-entry">
                    <div style="position: relative; display: inline-block;">
                        <input class="input" id="insertEntryName" type="text" spellcheck="false" placeholder="Enter name" size="40"/>
                        <svg
                            id="generate-uuid-button"
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
                    <input class="input" id="insertEntryValue" type="text" spellcheck="false" placeholder="Enter value" style="width:400px;"/>
                    <button class="btn" id="insert-entry-button">Insert</button>
                    <div id="entryMetadata" class="section">
                        <div class="metadata-entry">
                            <input class="input jsonKey" type="text" spellcheck="false" value="name" readonly style="margin-right: 5px;">
                            <input class="input jsonValue" type="text" spellcheck="false" placeholder="Resource Title" required>
                        </div>
                        <div class="metadata-entry">
                            <input class="input jsonKey" type="text" spellcheck="false" value="external" readonly style="margin-right: 5px;">
                            <select class="input jsonValue" id="externalMetadataToggle" style="width:422px;" required>
                                <option value="default" selected disabled>Resource Is External</option>
                                <option value="true">True</option>
                                <option value="false">False</option>
                            </select>
                        </div>
                        <div class="metadata-entry">
                            <input class="input jsonKey" type="text" spellcheck="false" value="mimetype" readonly style="margin-right: 5px;">
                            <input class="input jsonValue" type="text" spellcheck="false" placeholder="Resource MimeType" required>
                        </div>
                        <div class="metadata-entry">
                            <input class="input jsonKey" type="text" spellcheck="false" value="location" readonly style="margin-right: 5px;">
                            <input class="input jsonValue" type="text" spellcheck="false" placeholder="Resource Location (domain & owner email)" required>
                        </div>
                        <div class="metadata-entry">
                            <input class="input jsonKey" type="text" spellcheck="false" value="description" readonly style="margin-right: 5px;">
                            <input class="input jsonValue" type="text" spellcheck="false" placeholder="Resource Description">
                        </div>
                        <div id="cloud-storage-id-div" class="metadata-entry" style="display:none;">
                            <input class="input jsonKey" type="text" spellcheck="false" value="cloud_storage_id" readonly style="margin-right:-5px;">
                            <input class="input jsonValue" type="text" spellcheck="false" placeholder="Resource Cloud Storage ID">
                        </div>
                        <div id="md5checksum-div" class="metadata-entry" style="display:none;">
                            <input class="input jsonKey" type="text" spellcheck="false" value="md5Checksum" readonly style="margin-right:-5px;">
                            <input class="input jsonValue" type="text" spellcheck="false" placeholder="Resource MD5 Checksum">
                        </div>
                    </div>
                </div>
            `;
            document.getElementById("clear-insert").style.display = "inline";
            document.getElementById("generate-uuid-button").addEventListener("click", generateUUID);
            document.getElementById("insert-entry-button").addEventListener("click", insertEntry);
            document.getElementById("externalMetadataToggle").addEventListener("change", updateExternalInternalMetadataSelector);
            break;

        case "fromFile":
            dynamicInsertEntryDiv.innerHTML = `
                <div class="section" id="file-insert-entry">
                    <input class="input" id="insertFile" type="file" accept=".csv" style="border:0px;background-color:transparent;"/>
                    <button class="btn" id="insert-file-button">Insert</button>
                </div>
            `;
            document.getElementById("clear-insert").style.display = "inline";
            document.getElementById("insertFile").addEventListener("change", (event) => readFileContent(event.target));
            document.getElementById("insert-file-button").addEventListener("click", insertEntryFromFile);
            break;

        case "getBulkInsertTemplate":
            GenerateCSV();
        default:
            document.getElementById("insertEntrySelector").value = "default";
            dynamicInsertEntryDiv.innerHTML = `<div class="result section"></div>`;
            document.getElementById("clear-insert").style.display = "none";
            break;
    }
};

window.searchEntry = async function (event) {
    if (event?.type === "keydown" && event.key !== "Enter") {
        return;
    }

    const entryValueElement = document.getElementById("entryValue");
    const searchTypeElement = document.getElementById("searchType");

    const value = entryValueElement.value.trim();
    const searchType = searchTypeElement.value;

    if (searchType !== "GetAllEntries" && value === "") {
        ShowAlert("Please enter a search value.");
        return;
    }

    try {
        window.cachedEntries = [];

        switch (searchType) {
            case "GetEntryByName": {
                const entryByName = await GetEntryByName(getUUIDFromString(value));
                if (entryByName?.Name) {
                    window.cachedEntries.push(entryByName);
                }
                break;
            }

            case "GetEntryByValue": {
                const entryByValue = await GetEntryByValue(value);
                if (entryByValue?.Name) {
                    window.cachedEntries.push(entryByValue);
                }
                break;
            }

            case "GetEntriesByValue":
                window.cachedEntries = await GetEntriesByValue(value) ?? [];
                break;

            case "GetAllEntries":
                window.cachedEntries = await GetAllEntries() ?? [];
                break;

            default:
                updateResults("Invalid search type.");
                return;
        }

        if (window.cachedEntries.length > 0) {
            initializeFuse(window.cachedEntries);
            displayEntries(window.cachedEntries);
        } else {
            updateResults("No entries found for the provided value.");
        }
    } catch (err) {
        updateResults(`An error occurred while fetching the entries. ${err}`);
    }
};

window.removeMetaDataEntryField = function () {
    const entryMetadata = document.getElementById("entryMetadata");
    entryMetadata.removeChild(entryMetadata.lastChild);
};

window.deleteEntry = async function (event) {
    if (event?.type === "keydown" && event.key !== "Enter") {
        return;
    }

    try {
        const uuid = getUUIDFromString(document.getElementById("deleteEntryName").value);
        if (uuid === '') {
            ShowAlert("Must enter a valid UUID\nUse Search All to see a list of all current UUIDs");
            clearDeleteField();
            return;
        }

        await DeleteKeyValue(uuid);
        await DeleteName(uuid);
        clearDeleteField();
    } catch (err) {
        ShowAlert(`Error deleting record. ${err}`);
    }
};

window.clearResults = function () {
    const entryValueElement = document.getElementById("entryValue");
    updateResults();
    if (entryValueElement) {
        entryValueElement.value = '';
    }
};

window.clearDeleteField = function () {
    const deleteField = document.getElementById("deleteEntryName");
    if (deleteField) {
        deleteField.value = '';
    }
};

window.generateUUID = function () {
    const entryNameInput = document.getElementById('insertEntryName');
    entryNameInput.value = crypto.randomUUID();
};

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
            if (key === 'external') {
                value = (value === 'true');
            }
        } else {
            value = valueInput.value.trim();
        }

        if (key && value !== '' && value !== 'default') {
            metadata[key] = value;
        }
    });

    const value = document.getElementById("insertEntryValue").value.trim();
    const name = document.getElementById("insertEntryName").value.trim();

    if (!name || !value) {
        ShowAlert("Please provide both Name and Value.");
        return;
    }

    if (metadata['external'] === undefined || metadata['external'] === 'default') {
        ShowAlert("Please select whether the resource is external.");
        return;
    }

    try {
        const metadataString = JSON.stringify(metadata);
        const response = await InsertKVEntry(name, value, metadataString);

        if (response && response.success) {
            await InsertKVEntryIntoDatabase(name, value, metadataString);
            updateInsertEntry("");
            ShowAlert(`Successfully inserted ${metadata["name"]}`);
        } else {
            ShowAlert('Failed to insert entry: ' + response.errors.join(', '));
        }
    } catch (error) {
        ShowAlert(`An error occurred while inserting the entry. ${error}`);
    }
};

window.insertFromFileContent = null;
window.insertFromFileContentResolver = null;

window.clearInsertFromFile = function () {
    window.insertFromFileContent = null;
    const insertFile = document.getElementById('insertFile');
    if (insertFile) {
        insertFile.value = '';
    }
};

window.readFileContent = function (input) {
    const file = input.files[0];
    if (file) {
        const reader = new FileReader();
        const fileContentPromise = new Promise((resolve) => {
            window.insertFromFileContentResolver = resolve;
        });

        reader.onload = function (e) {
            window.insertFromFileContent = e.target.result;
            if (window.insertFromFileContentResolver) {
                window.insertFromFileContentResolver(window.insertFromFileContent);
            }
        };

        reader.readAsText(file);
        return fileContentPromise;
    }

    return Promise.reject("No file selected");
};

window.insertEntryFromFile = async function () {
    const content =
        window.insertFromFileContent ||
        await window.readFileContent(document.getElementById("insertFile"));

    if (!content || !content.trim()) return;

    const parsed = Papa.parse(content, {
        header: true,
        skipEmptyLines: "greedy",
        transformHeader: (header) => header.trim().replace(/\r/g, ""),
        transform: (value) => value.trim().replace(/\r/g, "")
    });

    if (parsed.errors?.length) {
        const firstError = parsed.errors[0];
        ShowAlert(`CSV parse error on row ${firstError.row ?? "unknown"}: ${firstError.message}`);
        return;
    }

    const rows = parsed.data || [];
    const errors = [];
    let insertedCount = 0;

    for (let index = 0; index < rows.length; index++) {
        const rowNumber = index + 2;
        const rowData = rows[index];

        const name = rowData.name || "";
        const value = rowData.value || "";

        if (!name || !value) {
            errors.push(`Row ${rowNumber}: missing name or value.`);
            continue;
        }

        const metadata = {};

        for (const [key, rawValue] of Object.entries(rowData)) {
            if (!key.startsWith("metadata_")) continue;
            if (rawValue === "") continue;

            const metaKey = key.slice("metadata_".length);
            let parsedValue = rawValue;

            if (metaKey === "external") {
                const normalized = String(rawValue).toLowerCase();
                if (normalized === "true") parsedValue = true;
                else if (normalized === "false") parsedValue = false;
                else {
                    errors.push(`Row ${rowNumber}: metadata_external must be true or false.`);
                    parsedValue = undefined;
                }
            }

            if (parsedValue !== undefined) {
                metadata[metaKey] = parsedValue;
            }
        }

        if (metadata.external === undefined) {
            errors.push(`Row ${rowNumber}: metadata_external is required.`);
            continue;
        }

        try {
            const metadataString = JSON.stringify(metadata);
            const response = await InsertKVEntry(name, value, metadataString);

            if (!response || response.success === undefined || response.success) {
                await InsertKVEntryIntoDatabase(name, value, metadataString);
                insertedCount++;
            } else {
                const errorText = Array.isArray(response.errors)
                    ? response.errors.join(", ")
                    : "Unknown error";
                errors.push(`Row ${rowNumber}: failed to insert "${name}" - ${errorText}`);
            }
        } catch (error) {
            errors.push(`Row ${rowNumber}: exception while inserting "${name}" - ${error}`);
        }
    }

    clearInsertFromFile();

    if (errors.length > 0) {
        ShowAlert(
            `Inserted ${insertedCount} entr${insertedCount === 1 ? "y" : "ies"}.\n\nErrors:\n${errors.join("\n")}`
        );
        return;
    }

    ShowAlert(`Successfully inserted ${insertedCount} entr${insertedCount === 1 ? "y" : "ies"}.`);
};

function updateResults(content = '') {
    const resultElement = document.getElementById("entryResult");
    const clearResultsButton = document.getElementById("clear");

    if (!resultElement || !clearResultsButton) return;

    resultElement.innerHTML = content;
    clearResultsButton.style.display = content ? 'inline' : 'none';
}

function displayEntries(entries) {
    let tableHTML = `
        <div class="section" id="table-search">
            <label for="approximateSearchValue" style="font-style:italic;">Search Table:</label>
            <input class="input" id="approximateSearchValue" type="text" autocomplete="off" spellcheck="false" placeholder="Search..." style="width:400px;"/>
            <span id="numberOfRecords" style="font-style:italic;">${entries.length} Records</span>
        </div>
        <table id="resultTable" style="margin-bottom:10px;table-layout:fixed; width:100%;">
            <colgroup>
                <col style="width:400px;">
                <col style="width:400px;">
                <col style="width:400px;">
                <col style="width:250px;">
                <col style="width:250px;">
                <col style="width:350px;">
                <col style="width:320px;">
                <col style="width:400px;">
            </colgroup>
            <thead>
                <tr>
                    <th data-column="UUID" class="sortable table-header">
                        ID
                        <span class="glyph sort-trigger">&#8645;</span>
                    </th>
                    <th data-column="Value" class="sortable table-header">
                        Value
                        <span class="glyph sort-trigger table-header">&#8645;</span>
                    </th>
                    <th data-column="Name" class="sortable table-header">
                        Name
                        <span class="glyph sort-trigger">&#8645;</span>
                    </th>
                    <th data-column="MimeType" class="sortable table-header">
                        Mime Type
                        <span class="glyph sort-trigger">&#8645;</span>
                    </th>
                    <th data-column="Location" class="sortable table-header">
                        Location
                        <span class="glyph sort-trigger">&#8645;</span>
                    </th>
                    <th data-column="CloudStorageId" class="sortable table-header">
                        Cloud Storage ID
                        <span class="glyph sort-trigger">&#8645;</span>
                    </th>
                    <th data-column="MD5Checksum" class="sortable table-header">
                        MD5 Checksum
                        <span class="glyph sort-trigger">&#8645;</span>
                    </th>
                    <th data-column="Description" class="sortable table-header">
                        Description
                        <span class="glyph sort-trigger">&#8645;</span>
                    </th>
                </tr>
            </thead>
            <tbody id="resultTableBody">
    `;

    entries.forEach(entry => {
        tableHTML += `
            <tr>
                <td>
                    <span class="copyonclick">${entry.Name}</span>
                    <span class="copyonclick glyphicon glyphicon-link" data-copy="${buildEntryLink(entry.Name)}"></span>
                </td>
                <td class="copyonclick">${entry.Value}</td>
                <td class="copyonclick">${entry.Metadata?.name ?? ''}</td>
                <td class="copyonclick">${entry.Metadata?.mimetype ?? ''}</td>
                <td class="copyonclick">${entry.Metadata?.location ?? ''}</td>
                <td class="copyonclick">${entry.Metadata?.cloud_storage_id ?? ''}</td>
                <td class="copyonclick">${entry.Metadata?.md5Checksum ?? ''}</td>
                <td class="copyonclick">${entry.Metadata?.description ?? ''}</td>
            </tr>
        `;
    });

    tableHTML += `
            </tbody>
        </table>
    `;

    updateResults(tableHTML);

    const approximateSearchValue = document.getElementById("approximateSearchValue");
    if (approximateSearchValue) {
        approximateSearchValue.addEventListener("input", approximateSearch);
    }

    enableSorting();
    enableCopyOnClick();
}

function enableSorting() {
    const table = document.getElementById("resultTable");
    if (!table) return;

    const headers = table.querySelectorAll("th.sortable");
    let sortDirection = 1;

    const sortTriggers = table.querySelectorAll(".sort-trigger");
    sortTriggers.forEach(trigger => {
        trigger.addEventListener("click", (event) => {
            const header = event.target.closest("th");
            const columnIndex = Array.from(headers).indexOf(header) + 1;
            const rows = Array.from(table.querySelector("tbody").rows);

            rows.sort((a, b) => {
                const aText = a.querySelector(`td:nth-child(${columnIndex})`).textContent.trim();
                const bText = b.querySelector(`td:nth-child(${columnIndex})`).textContent.trim();

                if (aText === '' && bText === '') return 0;
                if (aText === '') return 1;
                if (bText === '') return -1;

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

function displayClipboardMessage(message) {
    const textContent = message.trim();
    if (textContent === '') return;

    const messageElement = document.createElement('div');
    messageElement.innerText = `Copied: ${textContent}`;
    messageElement.className = "clipboard-message";

    navigator.clipboard.writeText(textContent).then(() => {
        document.body.appendChild(messageElement);
        setTimeout(() => {
            document.body.removeChild(messageElement);
        }, 2000);
    }).catch(err => {
        ShowAlert(`Error copying to clipboard: ${err}`);
    });
}

function enableCopyOnClick() {
    document.querySelectorAll('.copyonclick').forEach(element => {
        element.addEventListener('click', (e) => {
            let textValue = e.target.dataset.copy;
            if (textValue === undefined) {
                textValue = e.target.innerText;
            }
            displayClipboardMessage(textValue || '');
        });
    });
}

function getUUIDFromString(stringContainingUUID) {
    const uuidPattern = /\b[0-9a-fA-F]{8}(?:-[0-9a-fA-F]{4}){3}-[0-9a-fA-F]{12}\b/;
    const match = stringContainingUUID.match(uuidPattern);
    return match ? match[0] : '';
}

window.approximateSearch = function () {
    const query = document.getElementById("approximateSearchValue").value.trim();

    if (query === "") {
        displayEntries(window.cachedEntries);
        return;
    }

    const results = fuse.search(query);
    const filteredData = results.map((result) => result.item);

    displayApproximateSearchSort(filteredData);
};

function displayApproximateSearchSort(data) {
    const tableBody = document.getElementById("resultTableBody");
    if (!tableBody) return;

    tableBody.innerHTML = "";

    data.forEach(entry => {
        tableBody.innerHTML += `
            <tr>
                <td>
                    <span class="copyonclick">${entry.Name}</span>
                    <span class="copyonclick glyphicon glyphicon-link" data-copy="${buildEntryLink(entry.Name)}"></span>
                </td>
                <td class="copyonclick">${entry.Value}</td>
                <td class="copyonclick">${entry.Metadata?.name ?? ''}</td>
                <td class="copyonclick">${entry.Metadata?.mimetype ?? ''}</td>
                <td class="copyonclick">${entry.Metadata?.location ?? ''}</td>
                <td class="copyonclick">${entry.Metadata?.cloud_storage_id ?? ''}</td>
                <td class="copyonclick">${entry.Metadata?.md5Checksum ?? ''}</td>
                <td class="copyonclick">${entry.Metadata?.description ?? ''}</td>
            </tr>
        `;
    });

    document.getElementById("numberOfRecords").innerHTML = `${data.length} Records`;
    enableCopyOnClick();
}