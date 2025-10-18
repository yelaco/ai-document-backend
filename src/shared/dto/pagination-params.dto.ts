import { IsNumber, IsOptional, Max, Min } from 'class-validator';

export class PaginationParamsDto {
  @IsNumber()
  @IsOptional()
  @Min(1)
  page?: number;

  @IsNumber()
  @IsOptional()
  @Min(1)
  @Max(50)
  pageSize?: number;
}
