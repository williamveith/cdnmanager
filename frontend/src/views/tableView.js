import Fuse from 'fuse.js';
import { appState, setFuse } from '../state/appState';
import { buildEntryLink } from '../utils/domain';
import { copyWithToast } from '../utils/clipboard';
import { ShowAlert } from '../services/appService';

const fuseOptions = {
  keys: ['Metadata.name', 'Metadata.mimetype', 'Metadata.location', 'Metadata.description'],
  threshold: 0.3,
  includeScore: true
};

export function initializeFuse(entries) {
  setFuse(new Fuse(entries, fuseOptions));
}

export function updateResults(content = '') {
  const resultElement = document.getElementById('entryResult');
  const clearResultsButton = document.getElementById('clear');

  if (!resultElement || !clearResultsButton) return;

  resultElement.innerHTML = content;
  clearResultsButton.style.display = content ? 'inline' : 'none';
}

export function displayEntries(entries) {
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
        ${entries.map(renderRow).join('')}
      </tbody>
    </table>
  `;

  updateResults(tableHTML);

  const approximateSearchValue = document.getElementById('approximateSearchValue');
  if (approximateSearchValue) {
    approximateSearchValue.addEventListener('input', approximateSearch);
  }

  enableCopyOnClick();
  enableSorting();
}

function renderRow(entry) {
  return `
    <tr>
      <td>
        <span class="copyonclick">${entry.Name}</span>
        <span class="copyonclick glyphicon glyphicon-link" data-copy="${buildEntryLink(appState.appDomain, entry.Name)}"></span>
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

function enableCopyOnClick() {
  document.querySelectorAll('.copyonclick').forEach(element => {
    element.addEventListener('click', async (e) => {
      const textValue = e.target.dataset.copy ?? e.target.innerText;
      await copyWithToast(textValue || '', ShowAlert);
    });
  });
}

export function approximateSearch() {
  const query = document.getElementById('approximateSearchValue').value.trim();

  if (query === '') {
    displayEntries(appState.cachedEntries);
    return;
  }

  const results = appState.fuse.search(query);
  const filteredData = results.map(result => result.item);

  const tableBody = document.getElementById('resultTableBody');
  if (!tableBody) return;

  tableBody.innerHTML = filteredData.map(renderRow).join('');
  document.getElementById('numberOfRecords').innerHTML = `${filteredData.length} Records`;
  enableCopyOnClick();
}