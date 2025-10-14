import { Expose, instanceToPlain, plainToInstance } from 'class-transformer';
import { DocumentDto } from '../../documents/dto/document.dto';
import { User } from '../entities/user.entity';

import { ApiProperty } from '@nestjs/swagger';

export class UserDto {
  @ApiProperty({ example: 1 })
  @Expose()
  id: number;

  @ApiProperty({ example: 'user@example.com' })
  @Expose()
  email: string;

  @ApiProperty({ example: 'John Doe' })
  @Expose()
  fullName: string;

  @ApiProperty({ example: [] })
  @Expose()
  documents: DocumentDto[];
}

export function toUserDto(user: User): UserDto {
  const plainUser = instanceToPlain(user);
  return plainToInstance(UserDto, plainUser);
}

export function toUserEntity(dto: any): User {
  const plainUser = instanceToPlain(dto);
  return plainToInstance(User, plainUser);
}
