import { Injectable, InternalServerErrorException } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { SecretManagerServiceClient } from '@google-cloud/secret-manager';
import { AppSecrets, SecretManagerService } from './secret-manager.service';

@Injectable()
export class GcpSecretManagerService implements SecretManagerService {
  private client: SecretManagerServiceClient;
  private projectId: string;

  // Cache for secrets to avoid multiple fetches
  private secrets: AppSecrets;

  constructor(configService: ConfigService) {
    this.client = new SecretManagerServiceClient();
    // this.projectId = configService.getOrThrow<string>('gcp.projectId');
  }

  async getAppSecrets(): Promise<AppSecrets> {
    if (!this.secrets) {
      const secretString = await this.getSecret('backend');
      this.secrets = JSON.parse(secretString) as AppSecrets;
    }
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
      throw new InternalServerErrorException(
        `Failed to load secret ${secretName}: ${(error as Error).message}`,
      );
    }
  }
}
