import { Module } from '@nestjs/common';
import { AuthController } from './auth.controller';
import { AuthService } from './auth.service';
import { JwtModule } from '@nestjs/jwt';
import { ACCESS_TOKEN_EXPIRES_IN, ARGON2_OPTIONS } from './auth.constants';
import { argon2id } from 'argon2';
import { PassportModule } from '@nestjs/passport';
import { LocalStrategy } from './strategy/local.strategy';
import { UsersModule } from '../users/users.module';
import { SecretManagerModule } from '../secret-manager/secret-manager.module';
import { SecretManagerService } from '../secret-manager/secret-manager.service';
import { JwtStrategy } from './strategy/jwt.strategy';

@Module({
  imports: [
    UsersModule,
    SecretManagerModule,
    JwtModule.registerAsync({
      imports: [SecretManagerModule],
      inject: [SecretManagerService],
      useFactory: (secretManagerService: SecretManagerService) => {
        return {
          global: true,
          secret: secretManagerService.getAppSecrets().jwtSecret,
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
    JwtStrategy,
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
