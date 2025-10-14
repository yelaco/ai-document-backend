import { IsNotEmpty, IsString } from 'class-validator';

export class AskDocumentDto {
  @IsString()
  @IsNotEmpty()
  chatId: string;

  @IsString()
  @IsNotEmpty()
  question: string;
}

export class AnswerDocumentDto {
  status: 'in-progress' | 'completed';
  text?: string;
}
