/**
 * テキストをクリップボードにコピーする
 * HTTPS環境ではnavigator.clipboard、HTTP環境ではdocument.execCommandを使用
 */
export async function copyToClipboard(text: string): Promise<void> {
  try {
    if (navigator.clipboard && navigator.clipboard.writeText) {
      await navigator.clipboard.writeText(text);
    } else {
      // HTTP環境：古い手法を使用
      const textarea = document.createElement("textarea");
      textarea.value = text;
      textarea.style.position = "fixed";
      textarea.style.opacity = "0";
      document.body.appendChild(textarea);
      textarea.select();
      document.execCommand("copy");
      document.body.removeChild(textarea);
    }
  } catch (error) {
    console.error("Failed to copy:", error);
    throw error;
  }
}
