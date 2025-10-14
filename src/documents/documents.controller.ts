import {
  Controller,
  Post,
  UseInterceptors,
  UploadedFile,
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
import { PaginationResponseDto } from '../shared/dto/pagination-response.dto';
import { ApiTags, ApiOperation, ApiResponse } from '@nestjs/swagger';

@UseGuards(JwtAuthGuard)
@UseGuards(JwtAuthGuard)
@ApiTags('documents')
@Controller('documents')
export class DocumentsController {
  constructor(private readonly documentService: DocumentsService) {}

  @ApiOperation({
    summary: 'Upload document',
    description: 'Upload a new document via file.',
  })
  @ApiResponse({
    status: 201,
    description: 'Document uploaded.',
    schema: { example: { id: 'doc123', title: 'Uploaded file.pdf' } },
  })
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

  @ApiOperation({
    summary: 'List documents',
    description: 'Get a paginated list of documents.',
  })
  @ApiResponse({
    status: 200,
    description: 'Paginated document list.',
    schema: {
      example: {
        items: [{ id: 'doc123', title: 'First document' }],
        meta: { total: 1, page: 1, pageSize: 10 },
      },
    },
  })
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

  @ApiOperation({
    summary: 'Get document by id',
    description: 'Retrieve a single document by its ID.',
  })
  @ApiResponse({
    status: 200,
    description: 'Single document details.',
    schema: { example: { id: 'doc123', title: 'Document title' } },
  })
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
