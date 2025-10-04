import { Inject, Injectable, MessageEvent, OnModuleInit } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Document } from './entities/document.entity';
import { Repository } from 'typeorm';
import pdf from 'pdf-parse';
import { RecursiveCharacterTextSplitter } from '@langchain/textsplitters';
import { AI_SERVICE } from '../ai/ai.constants';
import { type AiService } from '../ai/ai.service';
import { ChromaClient } from 'chromadb';
import { CreateDocumentDto } from './dto/create-document.dto';
import { AskDocumentDto } from './dto/ask-document.dto';
import { endWith, map, Observable } from 'rxjs';
import { buildPrompt } from '../ai/ai.prompts';
import { CHROMA_CLIENT } from './documents.constants';

@Injectable()
export class DocumentsService implements OnModuleInit {
  constructor(
    @InjectRepository(Document)
    private documentsRepository: Repository<Document>,

    @Inject(AI_SERVICE)
    private readonly aiService: AiService,

    @Inject(CHROMA_CLIENT)
    private readonly chroma: ChromaClient,
  ) {}

  async onModuleInit() {
    await this.chroma.heartbeat();
  }

  async process(file: Express.Multer.File, documentId: string): Promise<void> {
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

    try {
      const collection = await this.chroma.getOrCreateCollection({
        name: 'documents',
        metadata: {
          created: new Date().toISOString(),
        },
      });

      await collection.add({
        ids: chunks.map((_, idx) => `${documentId}_chunk_${idx}`),
        embeddings: vectors,
        documents: chunks,
      });
    } catch (error) {
      if (error instanceof Error) {
        throw Error('Failed to store embeddings in ChromaDB: ' + error.message);
      }
    }
  }

  async ask(askDocumentDto: AskDocumentDto): Promise<Observable<MessageEvent>> {
    const embeddedQuery = await this.aiService.generateEmbedding([
      askDocumentDto.question,
    ]);
    if (embeddedQuery.length === 0 || embeddedQuery[0].length === 0) {
      throw new Error('Failed to embed the question');
    }

    const collection = await this.chroma.getCollection({ name: 'documents' });
    const result = await collection.query({
      queryEmbeddings: embeddedQuery,
    });
    const retrievedChunks = result.documents[0].filter((doc) => doc !== null);

    const observableAnswer = await this.aiService.answer(
      buildPrompt(askDocumentDto.question, retrievedChunks),
    );

    return observableAnswer.pipe(
      map(
        (text) =>
          ({
            data: {
              text: text,
              status: 'in-progress',
            },
          }) as MessageEvent,
      ),
      endWith({
        data: {
          status: 'completed',
        },
      } as MessageEvent),
    );
  }

  async create(createDocumentDto: CreateDocumentDto): Promise<Document> {
    const document = new Document();
    document.title = createDocumentDto.title;

    return await this.documentsRepository.save(document);
  }
}
