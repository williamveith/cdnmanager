import { SetupAndSync, GetDomain, ShowAlert } from '../services/appService';
import { normalizeDomain } from '../utils/domain';
import { appState } from '../state/appState';
import { renderMainShell } from '../views/shellView';
import { bindSearchEvents } from './searchController';
import { bindInsertEvents } from './insertController';
import { bindDeleteEvents } from './deleteController';

const appRoot = document.querySelector('#app');

export function bindConfigEvents() {
  const form = document.getElementById('config-form');
  if (!form) return;

  form.addEventListener('submit', async (e) => {
    e.preventDefault();

    const cfg = {
      cloudflare_api_token: document.getElementById('config-cloudflare-api-token').value.trim(),
      account_id: document.getElementById('config-account-id').value.trim(),
      namespace_id: document.getElementById('config-namespace-id').value.trim(),
      domain: document.getElementById('config-domain').value.trim()
    };

    if (!cfg.cloudflare_api_token || !cfg.account_id || !cfg.namespace_id || !cfg.domain) {
      ShowAlert('All configuration fields are required.');
      return;
    }

    try {
      document.getElementById('config-status').innerHTML = 'Saving configuration and syncing Cloudflare data...';
      await SetupAndSync(cfg);
      appState.appDomain = normalizeDomain(await GetDomain());

      renderMainShell(appRoot);
      bindSearchEvents();
      bindInsertEvents();
      bindDeleteEvents();

      ShowAlert('Configuration saved and database synced.');
    } catch (err) {
      document.getElementById('config-status').innerHTML = '';
      ShowAlert(`Failed to save configuration or sync database. ${err}`);
    }
  });
}