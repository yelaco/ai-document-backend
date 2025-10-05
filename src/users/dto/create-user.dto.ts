import { IsEmail, IsNotEmpty, IsString, Length } from 'class-validator';

export class CreateUserDto {
  @IsString()
  @IsEmail({})
  @IsNotEmpty()
  email: string;

  @IsString()
  @Length(8, 32)
  password: string;

  @IsString()
  @IsNotEmpty()
  fullName: string;
}
