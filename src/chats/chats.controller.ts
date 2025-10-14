import {
  Controller,
  Get,
  Post,
  Body,
  Patch,
  Param,
  Delete,
  NotFoundException,
  Query,
  Sse,
  UseGuards,
  MessageEvent,
} from '@nestjs/common';
import { ChatsService } from './chats.service';
import { CreateChatDto } from './dto/create-chat.dto';
import { UpdateChatDto } from './dto/update-chat.dto';
import {
  Context,
  type RequestContext,
} from '../shared/interceptors/request-context.interceptor';
import { ListChatDto } from './dto/list-chat.dto';
import { toChatDto } from './dto/chat.dto';
import { map, Observable } from 'rxjs';
import { AskDocumentDto } from '../chats/dto/ask-document.dto';
import { JwtAuthGuard } from '../auth/guard/jwt-auth.guard';
import { ApiOperation, ApiResponse, ApiBody } from '@nestjs/swagger';

@UseGuards(JwtAuthGuard)
@Controller('chats')
export class ChatsController {
  constructor(private readonly chatsService: ChatsService) {}

  @ApiOperation({
    summary: 'Ask document',
    description: 'Ask a question about a document.',
  })
  @ApiResponse({
    status: 200,
    description: 'Answer streamed as SSE',
    schema: {
      example: {
        data: { status: 'completed', text: 'Document answer here...' },
      },
    },
  })
  @Sse('ask')
  async ask(
    @Context() ctx: RequestContext,
    @Query() askDocumentDto: AskDocumentDto,
  ): Promise<Observable<MessageEvent>> {
    const answerObs = await this.chatsService.ask(
      ctx,
      askDocumentDto.chatId,
      askDocumentDto.question,
    );

    return answerObs.pipe(
      map(
        (dto) =>
          ({
            data: dto,
          }) as MessageEvent,
      ),
    );
  }

  @ApiOperation({ summary: 'Create chat', description: 'Create a new chat.' })
  @ApiBody({
    type: CreateChatDto,
    examples: {
      default: { value: { title: 'Chat about doc', documentId: 'doc123' } },
    },
  })
  @ApiResponse({
    status: 201,
    description: 'Chat created.',
    schema: { example: { id: 'chat1', title: 'Chat about doc' } },
  })
  @Post()
  create(@Context() ctx: RequestContext, @Body() createChatDto: CreateChatDto) {
    return this.chatsService.create(ctx, createChatDto);
  }

  @ApiOperation({
    summary: 'List chats',
    description: 'Get a paginated list of chats.',
  })
  @ApiResponse({
    status: 200,
    description: 'Paginated chat list.',
    schema: {
      example: {
        items: [{ id: 'chat1', title: 'First chat' }],
        meta: { total: 1, page: 1, pageSize: 10 },
      },
    },
  })
  @Get()
  async findPaginated(
    @Context() ctx: RequestContext,
    listChatDto: ListChatDto,
  ) {
    const { data, total } = await this.chatsService.findPaginated(
      ctx,
      listChatDto.page || 1,
      listChatDto.pageSize || 10,
    );
    return {
      items: data.map(toChatDto),
      meta: {
        total,
        page: listChatDto.page || 1,
        pageSize: listChatDto.pageSize || 10,
      },
    };
  }

  @ApiOperation({
    summary: 'Get chat by id',
    description: 'Retrieve a single chat by its ID.',
  })
  @ApiResponse({
    status: 200,
    description: 'Single chat details.',
    schema: { example: { id: 'chat1', title: 'Chat title' } },
  })
  @Get(':id')
  async findOne(@Context() ctx: RequestContext, @Param('id') id: string) {
    const chat = await this.chatsService.findOne(ctx, id);
    if (!chat) {
      throw new NotFoundException('Chat not found');
    }
    return toChatDto(chat);
  }

  @ApiOperation({
    summary: 'Update chat',
    description: 'Update the chat title.',
  })
  @ApiBody({
    type: UpdateChatDto,
    examples: { default: { value: { title: 'Updated chat title' } } },
  })
  @ApiResponse({
    status: 200,
    description: 'Updated chat.',
    schema: { example: { id: 'chat1', title: 'Updated chat title' } },
  })
  @Patch(':id')
  async update(
    @Context() ctx: RequestContext,
    @Param('id') id: string,
    @Body() updateChatDto: UpdateChatDto,
  ) {
    return this.chatsService.update(ctx, id, updateChatDto);
  }

  @ApiOperation({
    summary: 'Delete chat',
    description: 'Delete a chat by its ID.',
  })
  @ApiResponse({
    status: 200,
    description: 'Chat deleted.',
    schema: { example: { success: true } },
  })
  @Delete(':id')
  remove(@Context() ctx: RequestContext, @Param('id') id: string) {
    return this.chatsService.remove(ctx, id);
  }
}
