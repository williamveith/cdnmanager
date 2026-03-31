export function getUUIDFromString(value) {
  const uuidPattern = /\b[0-9a-fA-F]{8}(?:-[0-9a-fA-F]{4}){3}-[0-9a-fA-F]{12}\b/;
  const match = value.match(uuidPattern);
  return match ? match[0] : '';
}

export function generateUUID() {
  return crypto.randomUUID();
}