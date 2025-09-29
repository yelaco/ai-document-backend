export interface AiService {
  generateEmbedding(texts: string[]): Promise<number[][]>;
  summarizeDocument(text: string): Promise<string>;
}
