import { IsNumber, IsOptional } from 'class-validator';

export class ListChatDto {
  @IsNumber()
  @IsOptional()
  page?: number;

  @IsNumber()
  @IsOptional()
  pageSize?: number;
}
