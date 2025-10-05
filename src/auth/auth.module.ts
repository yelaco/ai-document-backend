import { Module } from '@nestjs/common';
import { AuthController } from './auth.controller';
import { AuthService } from './auth.service';
import { JwtModule } from '@nestjs/jwt';
import { ACCESS_TOKEN_EXPIRES_IN, ARGON2_OPTIONS } from './auth.constants';
import { argon2id } from 'argon2';
import { PassportModule } from '@nestjs/passport';
import { LocalStrategy } from './local.strategy';
import { UsersModule } from '../users/users.module';

@Module({
  imports: [
    UsersModule,
    JwtModule.register({
      global: true,
      secret: 'ai-document', // TODO: use secret vault
      signOptions: { expiresIn: ACCESS_TOKEN_EXPIRES_IN },
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
