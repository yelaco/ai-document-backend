import { Inject, Injectable } from '@nestjs/common';
import { AiService } from './ai.service';
import { GoogleGenAI } from '@google/genai';
import { from, map, Observable } from 'rxjs';
import { GEMINI_AI } from './ai.constants';

@Injectable()
export class GeminiAiService implements AiService {
  constructor(
    @Inject(GEMINI_AI)
    private readonly ai: GoogleGenAI,
  ) {}

  async answer(question: string): Promise<Observable<string>> {
    const response = await this.ai.models.generateContentStream({
      model: 'gemini-2.5-flash',
      contents: question,
      config: {
        candidateCount: 1,
        thinkingConfig: {
          thinkingBudget: 0,
        },
      },
    });
    return from(response).pipe(map((resp) => resp.text || ''));
  }
}
