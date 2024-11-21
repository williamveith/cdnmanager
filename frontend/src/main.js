import './style.css';
import './app.css';

import { GetEntryByName, GetEntryByValue, GetEntriesByValue, GetAllEntries } from '../wailsjs/go/database/Database';

document.querySelector('#app').innerHTML = `
    <div class="input-box" id="entry-input">
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

let searchTypeElement = document.getElementById("searchType");
let entryValueElement = document.getElementById("entryValue");
let resultElement = document.getElementById("entryResult");
let clearResultsButton = document.getElementById("clear")

searchTypeElement.addEventListener('change', () => {
    if (searchTypeElement.value === "GetAllEntries") {
        entryValueElement.style.display = 'none';
        entryValueElement.value = '';
    } else {
        entryValueElement.style.display = 'inline';
    }
});

window.searchEntry = async function () {
    let value = entryValueElement.value.trim();
    let searchType = searchTypeElement.value;

    if (searchType !== "GetAllEntries" && value === "") {
        resultElement.innerText = "Please enter a search value.";
        clearResultsButton.style.display = 'inline';
        return;
    }

    try {
        let entries = [];
        switch (searchType) {
            case "GetEntryByName":
                const entryByName = await GetEntryByName(value);
                if (entryByName?.Name) {
                    entries.push(entryByName);
                }
                break;
            case "GetEntryByValue":
                const entryByValue = await GetEntryByValue(value);
                if (entryByValue?.Name) {
                    entries.push(entryByValue);
                }
                break;
            case "GetEntriesByValue":
                entries = await GetEntriesByValue(value);
                break;
            case "GetAllEntries":
                entries = await GetAllEntries();
                break;
            default:
                resultElement.innerText = "Invalid search type.";
                clearResultsButton.style.display = 'inline';
                return;
        }

        if (entries.length > 0) {
            displayEntries(entries);
        } else {
            resultElement.innerText = "No entries found for the provided value.";
            clearResultsButton.style.display = 'inline';
        }
    } catch (err) {
        console.error(err);
        resultElement.innerText = "An error occurred while fetching the entries.";
        clearResultsButton.style.display = 'inline';
    }
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
        <table>
            <thead>
                <tr>
                    <th>Name</th>
                    <th>Value</th>
                    <th>Metadata</th>
                </tr>
            </thead>
            <tbody>
    `;

    entries.forEach(entry => {
        const metadataObject = JSON.parse(entry.Metadata);
        let metadataFormatted = '';
        for (const key in metadataObject) {
            if (metadataObject.hasOwnProperty(key)) {
                metadataFormatted += `${key}: ${metadataObject[key]}\n`;
            }
        }

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

    resultElement.innerHTML = tableHTML;
    const tdElements = document.querySelectorAll('.clickable');
    tdElements.forEach(td => {
        td.addEventListener('click', () => {
            navigator.clipboard.writeText(td.textContent.trim()).then(() => {
                displayClipboardMessage(`Copied: ${td.textContent.trim()}`);
            }).catch(err => {
                console.error('Error copying to clipboard:', err);
            });
        });
    });
    clearResultsButton.style.display = 'inline';
}

window.clearResults = function () {
    resultElement.innerHTML = '';
    clearResultsButton.style.display = 'none';
    entryValueElement.value = '';
}
