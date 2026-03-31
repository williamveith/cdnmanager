export function normalizeDomain(domain) {
  const trimmed = (domain ?? '').trim();
  if (!trimmed) return '';

  if (/^https?:\/\//i.test(trimmed)) {
    return trimmed.replace(/\/+$/, '');
  }

  return `https://${trimmed.replace(/\/+$/, '')}`;
}

export function buildEntryLink(appDomain, id) {
  if (!appDomain) {
    return `?id=${id}`;
  }

  return `${appDomain}/?id=${id}`;
}