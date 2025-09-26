export interface AiService {
  generateEmbedding(text: string): Promise<number[]>;
  summarizeDocument(text: string): Promise<string>;
}
