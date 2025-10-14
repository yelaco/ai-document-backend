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

@UseGuards(JwtAuthGuard)
@Controller('chats')
export class ChatsController {
  constructor(private readonly chatsService: ChatsService) {}

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

  @Post()
  create(@Context() ctx: RequestContext, @Body() createChatDto: CreateChatDto) {
    return this.chatsService.create(ctx, createChatDto);
  }

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

  @Get(':id')
  async findOne(@Context() ctx: RequestContext, @Param('id') id: string) {
    const chat = await this.chatsService.findOne(ctx, id);
    if (!chat) {
      throw new NotFoundException('Chat not found');
    }
    return toChatDto(chat);
  }

  @Patch(':id')
  async update(
    @Context() ctx: RequestContext,
    @Param('id') id: string,
    @Body() updateChatDto: UpdateChatDto,
  ) {
    return this.chatsService.update(ctx, id, updateChatDto);
  }

  @Delete(':id')
  remove(@Context() ctx: RequestContext, @Param('id') id: string) {
    return this.chatsService.remove(ctx, id);
  }
}
