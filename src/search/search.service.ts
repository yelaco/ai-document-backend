import { Inject, Injectable } from '@nestjs/common';
import { AI_SERVICE } from 'src/ai/ai.constants';
import { type AiService } from 'src/ai/ai.service';

@Injectable()
export class SearchService {
  constructor(@Inject(AI_SERVICE) private readonly aiService: AiService) {}

  async performSemanticSearch(query: string) {
    const queryVector = await this.aiService.generateEmbedding([query]);
    if (queryVector) {
      return { results: `Vector search complete: ${queryVector.toString()}` };
    }
  }
}
