import { Injectable } from '@nestjs/common';
import { CreateMessageDto } from './dto/create-message.dto';
import { UpdateMessageDto } from './dto/update-message.dto';
import { InjectQueue } from '@nestjs/bullmq';
import { JobsOptions, Queue } from 'bullmq';

@Injectable()
export class MessagesService {
  private readonly handleChatMessageJobOptions: JobsOptions;

  constructor(@InjectQueue('messages') private messageQueue: Queue) {
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

  async process(createMessageDto: CreateMessageDto) {
    await this.messageQueue.add(
      'handle-chat-message',
      createMessageDto,
      this.handleChatMessageJobOptions,
    );
    console.log('Added job to queue:', createMessageDto);
  }

  create(createMessageDto: CreateMessageDto) {
    return 'This action adds a new message';
  }

  findAll() {
    return `This action returns all messages`;
  }

  findOne(id: number) {
    return `This action returns a #${id} message`;
  }

  update(id: number, updateMessageDto: UpdateMessageDto) {
    return `This action updates a #${id} message`;
  }

  remove(id: number) {
    return `This action removes a #${id} message`;
  }
}
