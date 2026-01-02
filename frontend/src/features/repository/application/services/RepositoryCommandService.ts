export interface CreateRepositoryRequest {
  name: string;
  url: string;
  provider: string;
  accessToken?: string;
}

export interface UpdateRepositoryRequest {
  name?: string;
  url?: string;
}
