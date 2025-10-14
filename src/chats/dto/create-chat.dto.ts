import { IsNotEmpty, IsOptional, IsString } from 'class-validator';

export class CreateChatDto {
  @IsOptional()
  @IsString()
  title?: string;

  @IsString()
  @IsNotEmpty()
  documentId: string;
}
