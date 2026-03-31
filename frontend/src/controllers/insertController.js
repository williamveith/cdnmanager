import Papa from 'papaparse';

import { GenerateCSV, Insert, ShowAlert } from '../services/appService';
import { appState } from '../state/appState';
import { generateUUID } from '../utils/uuid';

export function bindInsertEvents() {
  const insertEntrySelector = document.getElementById('insertEntrySelector');
  const clearInsertButton = document.getElementById('clear-insert');

  if (!insertEntrySelector || !clearInsertButton) return;

  insertEntrySelector.addEventListener('change', () => {
    updateInsertEntry();
  });

  clearInsertButton.addEventListener('click', () => {
    updateInsertEntry('');
  });
}

function updateExternalInternalMetadataSelector() {
  const selectedValue = document.getElementById('externalMetadataToggle')?.value;
  const cloudStorageDiv = document.getElementById('cloud-storage-id-div');
  const md5ChecksumDiv = document.getElementById('md5checksum-div');

  if (!cloudStorageDiv || !md5ChecksumDiv) return;

  switch (selectedValue) {
    case 'true':
      cloudStorageDiv.style.display = 'none';
      md5ChecksumDiv.style.display = 'none';
      break;
    case 'false':
    default:
      cloudStorageDiv.style.display = 'block';
      md5ChecksumDiv.style.display = 'block';
      break;
  }
}

function updateInsertEntry(entryMethod = undefined) {
  const selector = document.getElementById('insertEntrySelector');
  const dynamicInsertEntryDiv = document.getElementById('dynamicInsertEntry');
  const clearInsertButton = document.getElementById('clear-insert');

  if (!selector || !dynamicInsertEntryDiv || !clearInsertButton) return;

  const selectedValue =
    entryMethod === undefined
      ? selector.value
      : entryMethod;

  switch (selectedValue) {
    case 'manual':
      dynamicInsertEntryDiv.innerHTML = `
        <div class="section" id="manual-insert-entry">
          <div style="position: relative; display: inline-block;">
            <input
              class="input"
              id="insertEntryName"
              type="text"
              spellcheck="false"
              placeholder="Enter name"
              size="40"
            />
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
              style="position:absolute;top:50%;right:10px;transform:translateY(-50%);cursor:pointer;"
            >
              <circle cx="12" cy="12" r="10"></circle>
              <line x1="12" y1="8" x2="12" y2="16"></line>
              <line x1="8" y1="12" x2="16" y2="12"></line>
            </svg>
          </div>

          <input
            class="input"
            id="insertEntryValue"
            type="text"
            spellcheck="false"
            placeholder="Enter value"
            style="width:400px;"
          />

          <button class="btn" id="insert-entry-button">Insert</button>

          <div id="entryMetadata" class="section">
            <div class="metadata-entry">
              <input
                class="input jsonKey"
                type="text"
                spellcheck="false"
                value="name"
                readonly
                style="margin-right:5px;"
              />
              <input
                class="input jsonValue"
                type="text"
                spellcheck="false"
                placeholder="Resource Title"
                required
              />
            </div>

            <div class="metadata-entry">
              <input
                class="input jsonKey"
                type="text"
                spellcheck="false"
                value="external"
                readonly
                style="margin-right:5px;"
              />
              <select
                class="input jsonValue"
                id="externalMetadataToggle"
                style="width:422px;"
                required
              >
                <option value="default" selected disabled>Resource Is External</option>
                <option value="true">True</option>
                <option value="false">False</option>
              </select>
            </div>

            <div class="metadata-entry">
              <input
                class="input jsonKey"
                type="text"
                spellcheck="false"
                value="mimetype"
                readonly
                style="margin-right:5px;"
              />
              <input
                class="input jsonValue"
                type="text"
                spellcheck="false"
                placeholder="Resource MimeType"
                required
              />
            </div>

            <div class="metadata-entry">
              <input
                class="input jsonKey"
                type="text"
                spellcheck="false"
                value="location"
                readonly
                style="margin-right:5px;"
              />
              <input
                class="input jsonValue"
                type="text"
                spellcheck="false"
                placeholder="Resource Location (domain & owner email)"
                required
              />
            </div>

            <div class="metadata-entry">
              <input
                class="input jsonKey"
                type="text"
                spellcheck="false"
                value="description"
                readonly
                style="margin-right:5px;"
              />
              <input
                class="input jsonValue"
                type="text"
                spellcheck="false"
                placeholder="Resource Description"
              />
            </div>

            <div id="cloud-storage-id-div" class="metadata-entry" style="display:none;">
              <input
                class="input jsonKey"
                type="text"
                spellcheck="false"
                value="cloud_storage_id"
                readonly
                style="margin-right:-5px;"
              />
              <input
                class="input jsonValue"
                type="text"
                spellcheck="false"
                placeholder="Resource Cloud Storage ID"
              />
            </div>

            <div id="md5checksum-div" class="metadata-entry" style="display:none;">
              <input
                class="input jsonKey"
                type="text"
                spellcheck="false"
                value="md5Checksum"
                readonly
                style="margin-right:-5px;"
              />
              <input
                class="input jsonValue"
                type="text"
                spellcheck="false"
                placeholder="Resource MD5 Checksum"
              />
            </div>
          </div>
        </div>
      `;

      clearInsertButton.style.display = 'inline';

      document.getElementById('generate-uuid-button')?.addEventListener('click', handleGenerateUUID);
      document.getElementById('insert-entry-button')?.addEventListener('click', insertEntry);
      document
        .getElementById('externalMetadataToggle')
        ?.addEventListener('change', updateExternalInternalMetadataSelector);
      break;

    case 'fromFile':
      dynamicInsertEntryDiv.innerHTML = `
        <div class="section" id="file-insert-entry">
          <input
            class="input"
            id="insertFile"
            type="file"
            accept=".csv"
            style="border:0;background-color:transparent;"
          />
          <button class="btn" id="insert-file-button">Insert</button>
        </div>
      `;

      clearInsertButton.style.display = 'inline';

      document.getElementById('insertFile')?.addEventListener('change', async (event) => {
        try {
          await readFileContent(event.target);
        } catch (err) {
          ShowAlert(`Failed to read file. ${err}`);
        }
      });

      document.getElementById('insert-file-button')?.addEventListener('click', insertEntryFromFile);
      break;

    case 'getBulkInsertTemplate':
      GenerateCSV();
    default:
      selector.value = 'default';
      dynamicInsertEntryDiv.innerHTML = `<div class="result section"></div>`;
      clearInsertButton.style.display = 'none';
      clearInsertFromFile();
      break;
  }
}

function handleGenerateUUID() {
  const entryNameInput = document.getElementById('insertEntryName');
  if (!entryNameInput) return;
  entryNameInput.value = generateUUID();
}

async function insertEntry() {
  const metadataEntries = document.querySelectorAll('.metadata-entry');
  const metadata = {};

  metadataEntries.forEach((entry) => {
    const keyInput = entry.querySelector('.jsonKey');
    const valueInput = entry.querySelector('.jsonValue');

    if (!keyInput || !valueInput) return;

    const key = keyInput.value.trim();
    let value;

    if (valueInput.tagName.toLowerCase() === 'select') {
      value = valueInput.options[valueInput.selectedIndex]?.value;
      if (key === 'external') {
        value = value === 'true';
      }
    } else {
      value = valueInput.value.trim();
    }

    if (key && value !== '' && value !== 'default') {
      metadata[key] = value;
    }
  });

  const value = document.getElementById('insertEntryValue')?.value.trim() ?? '';
  const name = document.getElementById('insertEntryName')?.value.trim() ?? '';

  if (!name || !value) {
    ShowAlert('Please provide both Name and Value.');
    return;
  }

  if (metadata.external === undefined || metadata.external === 'default') {
    ShowAlert('Please select whether the resource is external.');
    return;
  }

  try {
    const metadataString = JSON.stringify(metadata);
    await Insert(name, value, metadataString);
    updateInsertEntry('');
    ShowAlert(`Successfully inserted ${metadata.name}`);
  } catch (error) {
    ShowAlert(`An error occurred while inserting the entry. ${error}`);
  }
}

function clearInsertFromFile() {
  appState.insertFromFileContent = null;

  const insertFile = document.getElementById('insertFile');
  if (insertFile) {
    insertFile.value = '';
  }
}

function readFileContent(input) {
  const file = input?.files?.[0];

  if (!file) {
    return Promise.reject('No file selected');
  }

  return new Promise((resolve, reject) => {
    const reader = new FileReader();

    reader.onload = (event) => {
      appState.insertFromFileContent = event.target?.result ?? '';
      resolve(appState.insertFromFileContent);
    };

    reader.onerror = () => {
      reject('Unable to read file');
    };

    reader.readAsText(file);
  });
}

async function insertEntryFromFile() {
  try {
    const fileInput = document.getElementById('insertFile');

    const content =
      appState.insertFromFileContent ||
      (await readFileContent(fileInput));

    if (!content || !content.trim()) return;

    const parsed = Papa.parse(content, {
      header: true,
      skipEmptyLines: 'greedy',
      transformHeader: (header) => header.trim().replace(/\r/g, ''),
      transform: (value) => value.trim().replace(/\r/g, '')
    });

    if (parsed.errors?.length) {
      const firstError = parsed.errors[0];
      ShowAlert(`CSV parse error on row ${firstError.row ?? 'unknown'}: ${firstError.message}`);
      return;
    }

    const rows = parsed.data || [];
    const errors = [];
    let insertedCount = 0;

    for (let index = 0; index < rows.length; index += 1) {
      const rowNumber = index + 2;
      const rowData = rows[index];

      const name = rowData.name || '';
      const value = rowData.value || '';

      if (!name || !value) {
        errors.push(`Row ${rowNumber}: missing name or value.`);
        continue;
      }

      const metadata = {};

      for (const [key, rawValue] of Object.entries(rowData)) {
        if (!key.startsWith('metadata_')) continue;
        if (rawValue === '') continue;

        const metaKey = key.slice('metadata_'.length);
        let parsedValue = rawValue;

        if (metaKey === 'external') {
          const normalized = String(rawValue).toLowerCase();
          if (normalized === 'true') parsedValue = true;
          else if (normalized === 'false') parsedValue = false;
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
        await Insert(name, value, metadataString);
        insertedCount += 1;
      } catch (error) {
        errors.push(`Row ${rowNumber}: exception while inserting "${name}" - ${error}`);
      }
    }

    clearInsertFromFile();

    if (errors.length > 0) {
      ShowAlert(
        `Inserted ${insertedCount} entr${insertedCount === 1 ? 'y' : 'ies'}.\n\nErrors:\n${errors.join('\n')}`
      );
      return;
    }

    ShowAlert(`Successfully inserted ${insertedCount} entr${insertedCount === 1 ? 'y' : 'ies'}.`);
  } catch (error) {
    ShowAlert(`An error occurred while processing the CSV file. ${error}`);
  }
}