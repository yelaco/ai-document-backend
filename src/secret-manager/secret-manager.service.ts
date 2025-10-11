import { Injectable, OnModuleInit } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { SecretManagerServiceClient } from '@google-cloud/secret-manager';
import { DEV_SECRETS } from './secrets.dev';

export interface AppSecrets {
  jwtSecret: string;
}

@Injectable()
export class SecretManagerService implements OnModuleInit {
  private client: SecretManagerServiceClient;
  private projectId?: string;

  // Cache for secrets to avoid multiple fetches
  private secrets: AppSecrets;

  constructor(configService: ConfigService) {
    if (configService.get<string>('appEnv') === 'development') {
      this.useDevSecrets();
    } else {
      this.client = new SecretManagerServiceClient();
      this.projectId = configService.get<string>('gcp.projectId');
    }
  }

  async onModuleInit() {
    if (!this.secrets) {
      const secretString = await this.getSecret('backend');
      this.secrets = JSON.parse(secretString) as AppSecrets;
    }
  }

  useDevSecrets() {
    this.secrets = DEV_SECRETS;
  }

  getAppSecrets(): AppSecrets {
    return this.secrets;
  }

  async getSecret(secretName: string): Promise<string> {
    try {
      const [version] = await this.client.accessSecretVersion({
        name: `projects/${this.projectId}/secrets/${secretName}/versions/latest`,
      });

      const payload = version.payload?.data?.toString();

      if (!payload) {
        throw new Error(`Secret ${secretName} has no payload.`);
      }

      return payload;
    } catch (error) {
      throw new Error(
        `Failed to load secret ${secretName}: ${(error as Error).message}`,
      );
    }
  }
}
