import { IsOptional, IsUUID } from 'class-validator';
import { PaginationParamsDto } from 'src/shared/dto/pagination-params.dto';

export class ListChatDto extends PaginationParamsDto {
  @IsUUID()
  @IsOptional()
  documentId?: string;
}
