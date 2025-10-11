import { Injectable } from '@nestjs/common';
import { PassportStrategy } from '@nestjs/passport';
import { ExtractJwt, Strategy } from 'passport-jwt';
import { SecretManagerService } from '../secret-manager/secret-manager.service';

@Injectable()
export class JwtStrategy extends PassportStrategy(Strategy) {
  constructor(secretManagerService: SecretManagerService) {
    super({
      jwtFromRequest: ExtractJwt.fromAuthHeaderAsBearerToken(),
      ignoreExpiration: false,
      secretOrKey: secretManagerService.getAppSecrets().jwtSecret,
    });
  }

  async validate(payload: any) {
    return Promise.resolve({ userId: payload.sub, email: payload.username });
  }
}
