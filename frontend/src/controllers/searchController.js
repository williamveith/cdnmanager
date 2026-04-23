import { GetAllEntries, GetEntryByName, GetEntryByValue, GetEntriesByValue } from '../services/dbService';
import { ShowAlert } from '../services/appService';
import { getUUIDFromString } from '../utils/uuid';
import { appState } from '../state/appState';
import { displayEntries, initializeFuse, updateResults } from '../views/tableView';

export function bindSearchEvents() {
  const searchTypeElement = document.getElementById('searchType');
  const entryValueElement = document.getElementById('entryValue');
  const clearResultsButton = document.getElementById('clear');
  const searchButton = document.getElementById('search-button');

  if (!searchTypeElement || !entryValueElement || !clearResultsButton || !searchButton) return;

  searchTypeElement.addEventListener('change', () => {
    if (searchTypeElement.value === 'GetAllEntries') {
      entryValueElement.style.display = 'none';
      entryValueElement.value = '';
    } else {
      entryValueElement.style.display = 'inline';
    }
  });

  searchButton.addEventListener('click', searchEntry);
  entryValueElement.addEventListener('keydown', searchEntry);
  clearResultsButton.addEventListener('click', clearResults);
}

async function searchEntry(event) {
  if (event?.type === 'keydown' && event.key !== 'Enter') return;

  const entryValueElement = document.getElementById('entryValue');
  const searchTypeElement = document.getElementById('searchType');

  const value = entryValueElement.value.trim();
  const searchType = searchTypeElement.value;

  if (searchType !== 'GetAllEntries' && value === '') {
    ShowAlert('Please enter a search value.');
    return;
  }

  try {
    appState.cachedEntries = [];

    switch (searchType) {
      case 'GetEntryByName': {
        const entry = await GetEntryByName(getUUIDFromString(value));
        if (entry?.Name) appState.cachedEntries.push(entry);
        break;
      }
      case 'GetEntryByValue': {
        const entry = await GetEntryByValue(value);
        if (entry?.Name) appState.cachedEntries.push(entry);
        break;
      }
      case 'GetEntriesByValue':
        appState.cachedEntries = await GetEntriesByValue(value) ?? [];
        break;
      case 'GetAllEntries':
        appState.cachedEntries = await GetAllEntries() ?? [];
        break;
      default:
        updateResults('Invalid search type.');
        return;
    }

    if (appState.cachedEntries.length > 0) {
      initializeFuse(appState.cachedEntries);
      displayEntries(appState.cachedEntries);
    } else {
      updateResults('No entries found for the provided value.');
    }
  } catch (err) {
    updateResults(`An error occurred while fetching the entries. ${err}`);
  }
}

function clearResults() {
  const entryValueElement = document.getElementById('entryValue');
  updateResults();
  if (entryValueElement) {
    entryValueElement.value = '';
  }
}