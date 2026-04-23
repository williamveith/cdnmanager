export function renderConfigView(appRoot) {
  appRoot.innerHTML = `
    <form id="config-form">
      <div class="section" id="config-form-section">
        <div style="font-size:24px;font-weight:bold;margin-bottom:10px;">CDN Manager Setup</div>
        <div style="margin-bottom:16px;">Enter your Cloudflare configuration to initialize the application.</div>

        <div class="section">
          <input class="input"
              id="config-cloudflare-api-token"
              type="password"
              required
              pattern="[A-Za-z0-9_-]{30,}"
              title="Cloudflare API Token must contain only letters, numbers, underscores, or hyphens"
              spellcheck="false"
              placeholder="Cloudflare API Token"
              style="width:500px;" />
        </div>

        <div class="section">
          <input class="input"
              id="config-account-id"
              type="text"
              required
              pattern="[a-f0-9]{32}"
              title="Account ID must be 32 lowercase hexadecimal characters"
              spellcheck="false"
              placeholder="Account ID"
              style="width:500px;" />
        </div>

        <div class="section">
          <input class="input"
              id="config-namespace-id"
              type="text"
              required
              pattern="[a-f0-9]{32}"
              title="Namespace ID must be 32 lowercase hexadecimal characters"
              spellcheck="false"
              placeholder="Namespace ID"
              style="width:500px;" />
        </div>

        <div class="section">
          <input class="input"
              id="config-domain"
              type="text"
              required
              pattern="(https?:\\/\\/)?([a-zA-Z0-9-]+\\.)+[a-zA-Z]{2,}"
              title="Enter a valid domain such as cdn.example.com"
              spellcheck="false"
              placeholder="Domain"
              style="width:500px;" />
        </div>

        <div class="section" style="width:auto">
          <button class="btn" id="save-config-button" type="submit">Save & Sync</button>
        </div>

        <div class="result section" id="config-status"></div>
      </div>
    </form>
  `;
}