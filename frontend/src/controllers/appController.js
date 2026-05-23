import { IsConfigured, SyncFromCloudflare, GetDomain, ShowAlert } from '../services/appService';
import { normalizeDomain } from '../utils/domain';
import { appState } from '../state/appState';
import { renderConfigView } from '../views/configView';
import { renderMainShell } from '../views/shellView';
import { bindConfigEvents } from './configController';
import { bindSearchEvents } from './searchController';
import { bindInsertEvents } from './insertController';
import { bindDeleteEvents } from './deleteController';
import { bindExportEvents } from './exportController';

const appRoot = document.querySelector('#app');

export async function initializeApp() {
  try {
    const configured = await IsConfigured();

    if (!configured) {
      renderConfigView(appRoot);
      bindConfigEvents();
      return;
    }

    appState.appDomain = normalizeDomain(await GetDomain());
    await SyncFromCloudflare();

    renderMainShell(appRoot);
    bindSearchEvents();
    bindInsertEvents();
    bindDeleteEvents();
    bindExportEvents();
  } catch (err) {
    renderConfigView(appRoot);
    bindConfigEvents();
    ShowAlert(`Startup failed: ${err}`);
  }
}