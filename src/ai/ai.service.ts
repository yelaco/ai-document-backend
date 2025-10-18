import { Observable } from 'rxjs';
import { AiResponseDto } from './dto/ai-response.dto';

export interface AiService {
  answer(question: string): Promise<Observable<AiResponseDto>>;
}
