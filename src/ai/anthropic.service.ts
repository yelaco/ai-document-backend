import { Injectable } from '@nestjs/common';
import { AiService } from './ai.service';

@Injectable()
export class AnthropicAiService implements AiService {
  generateEmbedding(texts: string[]): Promise<number[][]> {
    throw new Error('Method not implemented.');
  }
  summarizeDocument(text: string): Promise<string> {
    throw new Error('Method not implemented.');
  }
}
