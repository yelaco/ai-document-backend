export class SearchMessageDto {
  chatId: string;
  query: string;
  topK?: number;
  timestamp?: Date;
}
