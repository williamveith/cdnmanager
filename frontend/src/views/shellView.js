export function renderMainShell(appRoot) {
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
}