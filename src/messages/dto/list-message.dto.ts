import { IsUUID } from 'class-validator';
import { PaginationParamsDto } from '../../shared/dto/pagination-params.dto';

export class ListMessageDto extends PaginationParamsDto {
  @IsUUID()
  chatId: string;
}
