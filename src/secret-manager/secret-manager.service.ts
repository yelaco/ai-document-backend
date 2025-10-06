export interface AppSecrets {
  JWT_SECRET: string;
}

export interface SecretManagerService {
  getAppSecrets(): Promise<AppSecrets>;
}
