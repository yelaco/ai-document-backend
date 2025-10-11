import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { SecretManagerService } from './secret-manager.service';

@Module({
  imports: [ConfigModule],
  providers: [SecretManagerService],
  exports: [SecretManagerService],
})
export class SecretManagerModule {}
