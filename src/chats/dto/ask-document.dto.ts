import { IsNotEmpty, IsString } from 'class-validator';

import { ApiProperty } from '@nestjs/swagger';

export class AskDocumentDto {
  @ApiProperty({ example: 'chat1' })
  @IsString()
  @IsNotEmpty()
  chatId: string;

  @ApiProperty({ example: 'What is the summary?' })
  @IsString()
  @IsNotEmpty()
  question: string;
}

export class AnswerDocumentDto {
  status: 'in-progress' | 'completed';
  text?: string;
}
