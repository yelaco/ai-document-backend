import { Observable } from 'rxjs';

export interface AiService {
  generateEmbedding(texts: string[]): Promise<number[][]>;
  answer(question: string): Promise<Observable<string>>;
}
