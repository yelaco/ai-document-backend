import { Injectable } from '@nestjs/common';
import { AiService } from './ai.service';
import { Observable } from 'rxjs';
import { AiResponseDto } from './dto/ai-response.dto';

@Injectable()
export class AnthropicAiService implements AiService {
  answer(question: string): Promise<Observable<AiResponseDto>> {
    throw new Error('Method not implemented.');
  }
}
