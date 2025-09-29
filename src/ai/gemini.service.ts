import { Inject, Injectable } from '@nestjs/common';
import { AiService } from './ai.service';
import { GoogleGenAI } from '@google/genai';
import { ConfigService } from '@nestjs/config';

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

  summarizeDocument(text: string): Promise<string> {
    throw new Error('Method not implemented.');
  }
}
