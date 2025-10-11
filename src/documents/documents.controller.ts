import {
  Controller,
  Post,
  UseInterceptors,
  UploadedFile,
  Sse,
  Query,
  UseGuards,
  Request,
} from '@nestjs/common';
import { DocumentsService } from './documents.service';
import { FileInterceptor } from '@nestjs/platform-express';
import { DocumentDto, toDocumentDto } from './dto/document.dto';
import { JwtAuthGuard } from '../auth/guard/jwt-auth.guard';
import {
  Context,
  type RequestContext,
} from '../shared/interceptors/request-context.interceptor';

@Controller('documents')
export class DocumentsController {
  constructor(private readonly documentService: DocumentsService) {}

  @UseGuards(JwtAuthGuard)
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
  ask(
    @Query('documentId') documentId: string,
    @Query('question') question: string,
  ) {
    return this.documentService.ask({ documentId, question });
  }
}
