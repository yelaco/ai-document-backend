import { Inject, Injectable } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Document } from './entities/document.entity';
import { Repository } from 'typeorm';
import pdf from 'pdf-parse';
import { RecursiveCharacterTextSplitter } from '@langchain/textsplitters';
import { AI_SERVICE } from 'src/ai/ai.constants';
import { type AiService } from 'src/ai/ai.service';

@Injectable()
export class DocumentsService {
  constructor(
    @InjectRepository(Document)
    private documentsRepository: Repository<Document>,

    @Inject(AI_SERVICE)
    private readonly aiService: AiService,
  ) {}

  async process(file: Express.Multer.File) {
    // Extract text from PDF
    const data = await pdf(file.buffer);

    // Chunk the text
    const splitter = new RecursiveCharacterTextSplitter({
      chunkSize: 1000,
      chunkOverlap: 200,
    });

    const chunks: string[] = await splitter.splitText(data.text);

    // Embed the chunks
    const vectors = await this.aiService.generateEmbedding(chunks);
    console.log(vectors);
  }
}
