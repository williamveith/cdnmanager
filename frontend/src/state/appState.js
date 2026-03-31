export const appState = {
  fuse: null,
  appDomain: '',
  cachedEntries: [],
  insertFromFileContent: null,
  insertFromFileContentResolver: null,
};

export function setFuse(fuseInstance) {
  appState.fuse = fuseInstance;
}