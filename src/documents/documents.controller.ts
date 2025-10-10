import {
  Controller,
  Post,
  UseInterceptors,
  UploadedFile,
  Sse,
  Query,
} from '@nestjs/common';
import { DocumentsService } from './documents.service';
import { FileInterceptor } from '@nestjs/platform-express';
import { DocumentDto, toDocumentDto } from './dto/document.dto';

@Controller('documents')
export class DocumentsController {
  constructor(private readonly documentService: DocumentsService) {}

  @Post()
  @UseInterceptors(FileInterceptor('file'))
  async upload(
    @UploadedFile() file: Express.Multer.File,
  ): Promise<DocumentDto> {
    const document = await this.documentService.create({
      title: file.originalname,
    });
    await this.documentService.process(file, document);
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
