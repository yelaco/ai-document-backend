import { AppSecrets, SecretManagerService } from './secret-manager.service';

export class AwsSecretManagerService implements SecretManagerService {
  getAppSecrets(): Promise<AppSecrets> {
    throw new Error('Method not implemented.');
  }
}
