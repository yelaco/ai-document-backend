import {
  Inject,
  Injectable,
  InternalServerErrorException,
  MessageEvent,
  UnauthorizedException,
} from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Document } from './entities/document.entity';
import { Repository } from 'typeorm';
import pdf from 'pdf-parse';
import { AI_SERVICE } from '../ai/ai.constants';
import { type AiService } from '../ai/ai.service';
import { CreateDocumentDto } from './dto/create-document.dto';
import { endWith, map, Observable } from 'rxjs';
import { buildPrompt } from '../ai/ai.prompts';
import { EmbeddingService } from '../embedding/embedding.service';
import { RequestContext } from '../shared/interceptors/request-context.interceptor';
import { AnswerDocumentDto } from './dto/ask-document.dto';

@Injectable()
export class DocumentsService {
  constructor(
    @InjectRepository(Document)
    private documentRepo: Repository<Document>,
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

      await this.embeddingService.embedDocument(document.id, data.text);
    } catch (error) {
      throw new InternalServerErrorException(
        'Failed to process document: ' + (error as Error).message,
      );
    }
  }

  async ask(
    ctx: RequestContext,
    documentId: string,
    question: string,
  ): Promise<Observable<AnswerDocumentDto>> {
    const questionContext = await this.embeddingService.vectorSearchDocument(
      documentId,
      question,
    );

    const observableAnswer = await this.aiService.answer(
      buildPrompt(question, questionContext),
    );

    return observableAnswer.pipe(
      map(
        (text) =>
          ({
            text: text,
            status: 'in-progress',
          }) as AnswerDocumentDto,
      ),
      endWith({
        status: 'completed',
      } as AnswerDocumentDto),
    );
  }

  async create(
    ctx: RequestContext,
    createDocumentDto: CreateDocumentDto,
  ): Promise<Document> {
    const document = this.documentRepo.create({
      title: createDocumentDto.title,
      userId: ctx.userId,
    });

    return this.documentRepo.save(document);
  }

  async findPaginated(
    ctx: RequestContext,
    page: number,
    pageSize: number,
  ): Promise<{ data: Document[]; total: number }> {
    const [data, total] = await this.documentRepo.findAndCount({
      skip: (page - 1) * pageSize,
      take: pageSize,
      order: { createdAt: 'DESC' },
    });
    return { data, total };
  }

  async findOne(ctx: RequestContext, id: string): Promise<Document | null> {
    if (!ctx.userId) {
      throw new UnauthorizedException();
    }

    return this.documentRepo.findOne({
      where: { id: id, userId: ctx.userId },
    });
  }

  async remove(ctx: RequestContext, id: string): Promise<void> {
    await this.documentRepo.delete(id);
  }
}
