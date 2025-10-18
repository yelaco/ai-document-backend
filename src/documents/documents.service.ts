import {
  Injectable,
  InternalServerErrorException,
  UnauthorizedException,
} from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Document } from './entities/document.entity';
import { Repository } from 'typeorm';
import pdf from 'pdf-parse';
import { CreateDocumentDto } from './dto/create-document.dto';
import { EmbeddingService } from '../embedding/embedding.service';
import { RequestContext } from '../shared/interceptors/request-context.interceptor';

@Injectable()
export class DocumentsService {
  constructor(
    @InjectRepository(Document)
    private documentRepo: Repository<Document>,
    private readonly embeddingService: EmbeddingService,
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
    const userId = ctx.userId;
    if (!userId) {
      throw new UnauthorizedException();
    }

    const [data, total] = await this.documentRepo.findAndCount({
      where: { userId },
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
