import { Controller, Get, Param, Delete, Query } from '@nestjs/common';
import { MessagesService } from './messages.service';
import { CreateMessageDto } from './dto/create-message.dto';
import { ApiTags, ApiOperation, ApiBody, ApiResponse } from '@nestjs/swagger';
import {
  Context,
  type RequestContext,
} from '../shared/interceptors/request-context.interceptor';
import { ListMessageDto } from './dto/list-message.dto';
import { PaginationResponseDto } from '../shared/dto/pagination-response.dto';
import { MessageDto, toMessageDto } from './dto/message.dto';

@ApiTags('messages')
@Controller('messages')
export class MessagesController {
  constructor(private readonly messagesService: MessagesService) {}

  @ApiOperation({
    summary: 'Create message',
    description: 'Create a new message.',
  })
  @ApiBody({ type: CreateMessageDto, examples: { default: { value: {} } } })
  @ApiResponse({
    status: 201,
    description: 'Message created.',
    schema: { example: { id: 1, content: 'Message text' } },
  })
  @Get()
  async findPaginated(
    @Context() ctx: RequestContext,
    @Query() listMessageDto: ListMessageDto,
  ): Promise<PaginationResponseDto<MessageDto>> {
    const { data, total } = await this.messagesService.findPaginated(
      ctx,
      listMessageDto.chatId,
      listMessageDto.page || 1,
      listMessageDto.pageSize || 10,
    );
    return {
      items: data.map(toMessageDto),
      meta: {
        total,
        page: listMessageDto.page || 1,
        pageSize: listMessageDto.pageSize || 10,
      },
    };
  }

  @Get(':id')
  findOne(@Context() ctx: RequestContext, @Param('id') id: string) {
    return this.messagesService.findOne(ctx, id);
  }

  @Delete(':id')
  remove(@Context() ctx: RequestContext, @Param('id') id: string) {
    return this.messagesService.remove(ctx, id);
  }
}
