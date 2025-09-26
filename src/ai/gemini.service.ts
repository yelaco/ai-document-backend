import { Injectable } from '@nestjs/common';
import { AiService } from './ai.service';

@Injectable()
export class GeminiAiService implements AiService {
  generateEmbedding(text: string): Promise<number[]> {
    throw new Error('Method not implemented.');
  }
  summarizeDocument(text: string): Promise<string> {
    throw new Error('Method not implemented.');
  }
}
