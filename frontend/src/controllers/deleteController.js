import { Delete, ShowAlert } from '../services/appService';
import { getUUIDFromString } from '../utils/uuid';

export function bindDeleteEvents() {
  const deleteButton = document.getElementById('delete-button');
  const deleteEntryName = document.getElementById('deleteEntryName');

  if (!deleteButton || !deleteEntryName) return;

  deleteButton.addEventListener('click', deleteEntry);
  deleteEntryName.addEventListener('keydown', deleteEntry);
}

async function deleteEntry(event) {
  if (event?.type === 'keydown' && event.key !== 'Enter') return;

  try {
    const deleteField = document.getElementById('deleteEntryName');
    const uuid = getUUIDFromString(deleteField.value);

    if (uuid === '') {
      ShowAlert('Must enter a valid UUID\nUse Search All to see a list of all current UUIDs');
      deleteField.value = '';
      return;
    }

    await Delete(uuid);
    deleteField.value = '';
  } catch (err) {
    ShowAlert(`Error deleting record. ${err}`);
  }
}