import { IsNotEmpty, IsString } from 'class-validator';

export class Auth {
  @IsString()
  @IsNotEmpty()
  accessToken: string;

  @IsString()
  refreshToken?: string;
}
