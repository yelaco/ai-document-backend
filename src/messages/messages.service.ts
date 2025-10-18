import {
  Injectable,
  InternalServerErrorException,
  NotFoundException,
} from '@nestjs/common';
import { CreateMessageDto } from './dto/create-message.dto';
import { InjectQueue } from '@nestjs/bullmq';
import { JobsOptions, Queue } from 'bullmq';
import { InjectRepository } from '@nestjs/typeorm';
import { Message } from './entities/message.entity';
import { Repository } from 'typeorm';
import { RequestContext } from '../shared/interceptors/request-context.interceptor';

@Injectable()
export class MessagesService {
  private readonly handleChatMessageJobOptions: JobsOptions;

  constructor(
    @InjectRepository(Message)
    private readonly messageRepo: Repository<Message>,
    @InjectQueue('message')
    private messageQueue: Queue,
  ) {
    this.handleChatMessageJobOptions = {
      attempts: 3,
      backoff: {
        type: 'exponential',
        delay: 5000,
      },
      removeOnComplete: true,
      removeOnFail: false,
    };
  }

  async process(ctx: RequestContext, createMessageDto: CreateMessageDto) {
    try {
      await this.create(ctx, createMessageDto);
      await this.messageQueue.add(
        'handle-chat-message',
        createMessageDto,
        this.handleChatMessageJobOptions,
      );
    } catch (error) {
      console.error('Error processing message:', error);
    }
  }

  async create(ctx: RequestContext, createMessageDto: CreateMessageDto) {
    const message = this.messageRepo.create(createMessageDto);
    return this.messageRepo.save(message);
  }

  async findPaginated(
    ctx: RequestContext,
    chatId: string,
    page: number,
    pageSize: number,
  ): Promise<{ data: Message[]; total: number }> {
    const [data, total] = await this.messageRepo.findAndCount({
      where: { chatId },
      skip: (page - 1) * pageSize,
      take: pageSize,
      order: { timestamp: 'DESC' },
    });
    return { data, total };
  }

  async findOne(ctx: RequestContext, id: string) {
    const message = await this.messageRepo.findOne({ where: { id } });
    if (!message) {
      throw new NotFoundException('Message not found');
    }

    return message;
  }

  async remove(ctx: RequestContext, id: string) {
    try {
      await this.messageRepo.delete(id);
    } catch {
      throw new InternalServerErrorException('Failed to delete message');
    }
  }
}
