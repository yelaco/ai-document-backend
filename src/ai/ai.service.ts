import { Observable } from 'rxjs';

export interface AiService {
  answer(question: string): Promise<Observable<string>>;
}
