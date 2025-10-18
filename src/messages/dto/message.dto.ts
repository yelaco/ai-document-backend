import { IsDateString, IsString, IsUUID } from 'class-validator';
import { Message } from '../entities/message.entity';
import { instanceToPlain, plainToInstance } from 'class-transformer';

export class MessageDto {
  @IsUUID()
  id: string;

  @IsString()
  role: string;

  @IsString()
  content: string;

  @IsString()
  chatId: string;

  @IsDateString()
  timestamp: Date;

  @IsDateString()
  createdAt: Date;

  @IsDateString()
  updatedAt: Date;
}

export function toMessageDto(message: Message): MessageDto {
  const plainMessage = instanceToPlain(message);
  return plainToInstance(MessageDto, plainMessage);
}

export function toMessageEntity(dto: any): Message {
  const plainMessage = instanceToPlain(dto);
  return plainToInstance(Message, plainMessage);
}
