import { Expose, instanceToPlain, plainToInstance } from 'class-transformer';
import { DocumentDto } from '../../documents/dto/document.dto';
import { User } from '../entities/user.entity';

export class UserDto {
  @Expose()
  id: number;

  @Expose()
  email: string;

  @Expose()
  fullName: string;

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
