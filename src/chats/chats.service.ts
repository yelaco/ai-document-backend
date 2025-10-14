import { Injectable, UnauthorizedException } from '@nestjs/common';
import { CreateChatDto } from './dto/create-chat.dto';
import { UpdateChatDto } from './dto/update-chat.dto';
import { Repository } from 'typeorm';
import { Chat } from './entities/chat.entity';
import { InjectRepository } from '@nestjs/typeorm';
import { RequestContext } from '../shared/interceptors/request-context.interceptor';

@Injectable()
export class ChatsService {
  constructor(
    @InjectRepository(Chat)
    private chatRepo: Repository<Chat>,
  ) {}

  async create(createChatDto: CreateChatDto) {
    const chat = this.chatRepo.create({
      title: createChatDto.title,
      documentId: createChatDto.documentId,
    });

    return this.chatRepo.save(chat);
  }

  async findPaginated(
    ctx: RequestContext,
    page: number,
    pageSize: number,
  ): Promise<{ data: Chat[]; total: number }> {
    const [data, total] = await this.chatRepo.findAndCount({
      skip: (page - 1) * pageSize,
      take: pageSize,
      order: { createdAt: 'DESC' },
    });
    return { data, total };
  }

  update(id: number, updateChatDto: UpdateChatDto) {
    return `This action updates a #${id} chat`;
  }

  async findOne(ctx: RequestContext, id: string): Promise<Chat | null> {
    if (!ctx.userId) {
      throw new UnauthorizedException();
    }

    return this.chatRepo.findOne({
      where: { id: id, userId: ctx.userId },
    });
  }

  async remove(ctx: RequestContext, id: string): Promise<void> {
    await this.chatRepo.delete(id);
  }
}
