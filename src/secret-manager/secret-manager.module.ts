import { Module } from '@nestjs/common';
import { GcpSecretManagerService } from './gcp.service';
import { SECRET_MANAGER_SERVICE } from './secret-manager.constants';
import { ConfigModule, ConfigService } from '@nestjs/config';
import { AwsSecretManagerService } from './aws.service';

@Module({
  imports: [ConfigModule],
  providers: [
    GcpSecretManagerService,
    AwsSecretManagerService,
    {
      provide: SECRET_MANAGER_SERVICE,
      inject: [ConfigService, GcpSecretManagerService, AwsSecretManagerService],
      useFactory: (
        configService: ConfigService,
        gcpSecretManagerService: GcpSecretManagerService,
        awsSecretManagerService: AwsSecretManagerService,
      ) => {
        const secretManagerServiceType = configService.get<string>(
          'gcp.secretManagerServiceType',
        );

        if (secretManagerServiceType === 'aws') {
          return awsSecretManagerService;
        }
        return gcpSecretManagerService;
      },
    },
  ],
  exports: [SECRET_MANAGER_SERVICE],
})
export class SecretManagerModule {}
