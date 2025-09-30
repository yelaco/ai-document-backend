import { Injectable } from '@nestjs/common';
import { AiService } from './ai.service';
import { Observable } from 'rxjs';

@Injectable()
export class AnthropicAiService implements AiService {
  generateEmbedding(texts: string[]): Promise<number[][]> {
    throw new Error('Method not implemented.');
  }
  answer(question: string): Promise<Observable<string>> {
    throw new Error('Method not implemented.');
  }
}
