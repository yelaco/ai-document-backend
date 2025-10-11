import { IsNumber, IsOptional } from 'class-validator';

export class ListDocumentDto {
  @IsNumber()
  @IsOptional()
  page?: number;

  @IsNumber()
  @IsOptional()
  pageSize?: number;
}
