import { Expose, instanceToPlain, plainToInstance } from 'class-transformer';
import { Chat } from '../entities/chat.entity';

export class ChatDto {
  @Expose()
  id: string;

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
