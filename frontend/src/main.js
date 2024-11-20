import './style.css';
import './app.css';

import { GetEntryByName } from '../wailsjs/go/database/Database';

document.querySelector('#app').innerHTML = `
    <div class="header" id="header">Enter Entry ID</div>
    <div class="input-box" id="entry-input">
        <input class="input" id="entryValue" type="text" autocomplete="off" />
        <button class="btn" onclick="searchEntry()">Search</button>
    </div>
    <div class="result" id="entryResult"></div>
`;

let entryValueElement = document.getElementById("entryValue");
entryValueElement.focus();
let resultElement = document.getElementById("entryResult");

window.searchEntry = async function () {
    // Get value
    let value = entryValueElement.value;

    // Check if the input is empty
    if (value === "") return;

    // Call Database.GetEntryByValue(value)
    try {
        const entry = await GetEntryByName(value);
        if (entry.Name) {
            resultElement.innerText = `
                Name: ${entry.Name}\n
                Value: ${entry.Value}\n
                Metadata: ${JSON.stringify(entry.Metadata, null, 2)}
            `;
        } else {
            resultElement.innerText = "No entry found for the provided value.";
        }
    } catch (err) {
        console.error(err);
        resultElement.innerText = "An error occurred while fetching the entry.";
    };
};
