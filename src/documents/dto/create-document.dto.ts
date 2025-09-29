import { IsNotEmpty, IsString, Length } from 'class-validator';

export class CreateDocumentDto {
  @IsString()
  @Length(1, 100)
  title: string;

  @IsNotEmpty()
  @IsString()
  content: string;

  @IsNotEmpty()
  vector: number[];
}
