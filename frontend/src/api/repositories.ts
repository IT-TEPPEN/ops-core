// API functions for managing repositories

// Updated to handle `window` object safely for testing environments
const apiHost = import.meta.env.VITE_API_HOST || (typeof window !== 'undefined' ? window.location.host : 'localhost');
const apiUrl = `${typeof window !== 'undefined' ? window.location.protocol : 'http:'}//${apiHost}/api/v1`;

export async function fetchRepositories(signal?: AbortSignal) {
  const response = await fetch(`${apiUrl}/repositories`, {
    signal,
    headers: {
      "Cache-Control": "max-age=60",
    },
  });
  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }
  return response.json();
}

export async function registerRepository(url: string) {
  const response = await fetch(`${apiUrl}/repositories`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ url }),
  });
  if (!response.ok) {
    const data = await response.json();
    throw new Error(data.message || "Failed to register repository");
  }
  return response.json();
}

export async function fetchRepositoryDetails(repoId: string) {
  const response = await fetch(`${apiUrl}/repositories/${repoId}`);
  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }
  return response.json();
}

export async function fetchRepositoryFiles(repoId: string) {
  const response = await fetch(`${apiUrl}/repositories/${repoId}/files`);
  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }
  return response.json();
}

export async function updateAccessToken(repoId: string, accessToken: string) {
  const response = await fetch(`${apiUrl}/repositories/${repoId}/token`, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ accessToken }),
  });
  if (!response.ok) {
    const data = await response.json();
    throw new Error(data.message || "Failed to update access token");
  }
  return response.json();
}

export async function selectFiles(repoId: string, filePaths: string[]) {
  const response = await fetch(`${apiUrl}/repositories/${repoId}/files/select`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ filePaths }),
  });
  if (!response.ok) {
    const data = await response.json();
    throw new Error(data.message || "Failed to select files");
  }
  return response.json();
}