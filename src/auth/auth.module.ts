import { Module } from '@nestjs/common';
import { AuthController } from './auth.controller';
import { AuthService } from './auth.service';
import { JwtModule } from '@nestjs/jwt';
import { ACCESS_TOKEN_EXPIRES_IN, ARGON2_OPTIONS } from './auth.constants';
import { argon2id } from 'argon2';
import { PassportModule } from '@nestjs/passport';
import { LocalStrategy } from './local.strategy';
import { UsersModule } from '../users/users.module';
import { SecretManagerModule } from '../secret-manager/secret-manager.module';
import { ConfigModule, ConfigService } from '@nestjs/config';
import { SECRET_MANAGER_SERVICE } from '../secret-manager/secret-manager.constants';
import { SecretManagerService } from '../secret-manager/secret-manager.service';

@Module({
  imports: [
    UsersModule,
    ConfigModule,
    SecretManagerModule,
    JwtModule.registerAsync({
      imports: [ConfigModule, SecretManagerModule],
      inject: [ConfigService, SECRET_MANAGER_SERVICE],
      useFactory: async (
        configService: ConfigService,
        secretManagerService: SecretManagerService,
      ) => {
        let jwtSecret = 'ai-document-secret';
        if (configService.get<string>('appEnv') === 'production') {
          const appSecrets = await secretManagerService.getAppSecrets();
          jwtSecret = appSecrets.JWT_SECRET;
        }
        return {
          global: true,
          secret: jwtSecret,
          signOptions: {
            expiresIn: ACCESS_TOKEN_EXPIRES_IN,
          },
        };
      },
    }),
    PassportModule,
  ],
  controllers: [AuthController],
  providers: [
    AuthService,
    LocalStrategy,
    {
      provide: ARGON2_OPTIONS,
      useValue: {
        type: argon2id,
        memoryCost: 65536, // 64 MiB
        timeCost: 4, // 4 iterations
        parallelism: 1,
      },
    },
  ],
})
export class AuthModule {}
