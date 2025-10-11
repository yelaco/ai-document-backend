export class PaginationResponseDto<T> {
  items: T[];
  meta: PaginationMetaDto;
}

export class PaginationMetaDto {
  total: number;
  page: number;
  pageSize: number;
}
