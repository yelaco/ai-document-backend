import { IsNotEmpty, IsString } from 'class-validator';

export class AskDocumentDto {
  @IsString()
  documentId: string;

  @IsString()
  @IsNotEmpty()
  question: string;
}
