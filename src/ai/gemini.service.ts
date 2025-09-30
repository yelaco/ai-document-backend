import { Inject, Injectable } from '@nestjs/common';
import { AiService } from './ai.service';
import { GoogleGenAI } from '@google/genai';
import { ConfigService } from '@nestjs/config';
import { from, map, Observable } from 'rxjs';

@Injectable()
export class GeminiAiService implements AiService {
  private ai: GoogleGenAI;

  constructor(@Inject() configService: ConfigService) {
    this.ai = new GoogleGenAI({
      apiKey: configService.get<string>('GEMINI_API_KEY'),
    });
  }

  async generateEmbedding(texts: string[]): Promise<number[][]> {
    const response = await this.ai.models.embedContent({
      model: 'gemini-embedding-001',
      contents: texts,
      config: {
        taskType: 'RETRIEVAL_DOCUMENT',
      },
    });
    if (!response.embeddings) {
      return Array<number[]>(texts.length).fill([]);
    }

    return response.embeddings
      .map((e) => e.values)
      .filter((e) => e != undefined);
  }

  async answer(question: string): Promise<Observable<string>> {
    const response = await this.ai.models.generateContentStream({
      model: 'gemini-2.0-flash-001',
      contents: question,
      config: {
        maxOutputTokens: 300,
        candidateCount: 1,
      },
    });
    return from(response).pipe(map((resp) => resp.text || ''));
  }
}
