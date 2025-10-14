import { Expose, instanceToPlain, plainToInstance } from 'class-transformer';
import { Chat } from '../entities/chat.entity';

import { ApiProperty } from '@nestjs/swagger';

export class ChatDto {
  @ApiProperty({ example: 'chat1' })
  @Expose()
  id: string;

  @ApiProperty({ example: 'Project Chat' })
  @Expose()
  title: string;
}

export function toChatDto(chat: Chat): ChatDto {
  const plainChat = instanceToPlain(chat);
  return plainToInstance(ChatDto, plainChat);
}

export function toChatEntity(dto: any): Chat {
  const plainChat = instanceToPlain(dto);
  return plainToInstance(Chat, plainChat);
}
