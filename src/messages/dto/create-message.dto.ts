import { IsDateString, IsEnum, IsString, IsUUID } from 'class-validator';
import { ChatRole } from '../../chats/chats.enum';

export class CreateMessageDto {
  @IsUUID()
  chatId: string;

  @IsEnum(ChatRole)
  role: string;

  @IsString()
  content: string;

  @IsDateString()
  timestamp: Date;
}
