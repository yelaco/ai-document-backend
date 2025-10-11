import {
  Controller,
  Post,
  UseInterceptors,
  UploadedFile,
  Sse,
  Query,
  UseGuards,
  NotFoundException,
  Get,
  Param,
  Delete,
} from '@nestjs/common';
import { DocumentsService } from './documents.service';
import { FileInterceptor } from '@nestjs/platform-express';
import { DocumentDto, toDocumentDto } from './dto/document.dto';
import { JwtAuthGuard } from '../auth/guard/jwt-auth.guard';
import {
  Context,
  type RequestContext,
} from '../shared/interceptors/request-context.interceptor';
import { ListDocumentDto } from './dto/list-document.dto';
import { AskDocumentDto } from './dto/ask-document.dto';
import { PaginationResponseDto } from '../shared/dto/pagination-response.dto';

@UseGuards(JwtAuthGuard)
@Controller('documents')
export class DocumentsController {
  constructor(private readonly documentService: DocumentsService) {}

  @Post()
  @UseInterceptors(FileInterceptor('file'))
  async upload(
    @Context() ctx: RequestContext,
    @UploadedFile() file: Express.Multer.File,
  ): Promise<DocumentDto> {
    const document = await this.documentService.create(ctx, {
      title: file.originalname,
    });
    await this.documentService.process(ctx, document, file);
    return toDocumentDto(document);
  }

  @Sse('ask')
  async ask(
    @Context() ctx: RequestContext,
    @Query() askDocumentDto: AskDocumentDto,
  ) {
    const document = await this.documentService.findOne(
      ctx,
      askDocumentDto.documentId,
    );
    if (!document) {
      throw new NotFoundException('Document not found');
    }
    return this.documentService.ask(
      ctx,
      askDocumentDto.documentId,
      askDocumentDto.question,
    );
  }

  @Get()
  async findPaginated(
    @Context() ctx: RequestContext,
    @Query() listDocumentDto: ListDocumentDto,
  ): Promise<PaginationResponseDto<DocumentDto>> {
    const { data, total } = await this.documentService.findPaginated(
      ctx,
      listDocumentDto.page || 1,
      listDocumentDto.pageSize || 10,
    );
    return {
      items: data.map(toDocumentDto),
      meta: {
        total,
        page: listDocumentDto.page || 1,
        pageSize: listDocumentDto.pageSize || 10,
      },
    };
  }

  @Get(':id')
  async findOne(
    @Context() ctx: RequestContext,
    @Param('id') id: string,
  ): Promise<DocumentDto> {
    const document = await this.documentService.findOne(ctx, id);
    if (!document) {
      throw new NotFoundException('Document not found');
    }
    return toDocumentDto(document);
  }

  @Delete(':id')
  async remove(
    @Context() ctx: RequestContext,
    @Param('id') id: string,
  ): Promise<void> {
    await this.documentService.remove(ctx, id);
  }
}
