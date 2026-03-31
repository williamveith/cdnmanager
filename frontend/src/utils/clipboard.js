export async function copyWithToast(text, showAlertFallback) {
  const textContent = (text ?? '').trim();
  if (!textContent) return;

  const messageElement = document.createElement('div');
  messageElement.innerText = `Copied: ${textContent}`;
  messageElement.className = 'clipboard-message';

  try {
    await navigator.clipboard.writeText(textContent);
    document.body.appendChild(messageElement);
    setTimeout(() => {
      if (messageElement.parentNode) {
        messageElement.parentNode.removeChild(messageElement);
      }
    }, 2000);
  } catch (err) {
    showAlertFallback(`Error copying to clipboard: ${err}`);
  }
}