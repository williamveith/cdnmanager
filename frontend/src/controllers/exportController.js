import { GenerateDatabaseCSV } from '../services/appService';

export function bindExportEvents() {
    const exportButton = document.getElementById('export-button');
    exportButton.addEventListener('click', exportDatabase);
}

function exportDatabase(event) {
    GenerateDatabaseCSV();
}