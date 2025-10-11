import {
  Inject,
  Injectable,
  InternalServerErrorException,
  MessageEvent,
} from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Document } from './entities/document.entity';
import { Repository } from 'typeorm';
import pdf from 'pdf-parse';
import { AI_SERVICE } from '../ai/ai.constants';
import { type AiService } from '../ai/ai.service';
import { CreateDocumentDto } from './dto/create-document.dto';
import { AskDocumentDto } from './dto/ask-document.dto';
import { endWith, map, Observable } from 'rxjs';
import { buildPrompt } from '../ai/ai.prompts';
import { EmbeddingService } from '../embedding/embedding.service';
import { RequestContext } from 'src/shared/interceptors/request-context.interceptor';

@Injectable()
export class DocumentsService {
  constructor(
    @InjectRepository(Document)
    private documentsRepository: Repository<Document>,
    private readonly embeddingService: EmbeddingService,
    @Inject(AI_SERVICE)
    private readonly aiService: AiService,
  ) {}

  async process(
    ctx: RequestContext,
    document: Document,
    file: Express.Multer.File,
  ): Promise<void> {
    try {
      // Extract text from PDF
      const data = await pdf(file.buffer);

      await this.embeddingService.embedDocument(data.text);
    } catch (error) {
      throw new InternalServerErrorException(
        'Failed to process document: ' + (error as Error).message,
      );
    }
  }

  async ask(askDocumentDto: AskDocumentDto): Promise<Observable<MessageEvent>> {
    const questionContext = await this.embeddingService.vectorSearchDocument(
      askDocumentDto.question,
    );

    const observableAnswer = await this.aiService.answer(
      buildPrompt(askDocumentDto.question, questionContext),
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

  async create(
    ctx: RequestContext,
    createDocumentDto: CreateDocumentDto,
  ): Promise<Document> {
    const document = this.documentsRepository.create({
      title: createDocumentDto.title,
      userId: ctx.userId,
    });

    return await this.documentsRepository.save(document);
  }
}
