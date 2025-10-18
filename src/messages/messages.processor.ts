import { Processor, WorkerHost } from '@nestjs/bullmq';
import { Job } from 'bullmq';
import { CreateMessageDto } from './dto/create-message.dto';
import { EmbeddingService } from '../embedding/embedding.service';

@Processor('message')
export class MessagesProcessor extends WorkerHost {
  constructor(private readonly embeddingService: EmbeddingService) {
    super();
  }

  async process(job: Job): Promise<any> {
    switch (job.name) {
      case 'handle-chat-message': {
        const message = job.data as CreateMessageDto;
        try {
          await this.embeddingService.embedChatMessage({
            chatId: message.chatId,
            content: message.content,
            timestamp: message.timestamp,
            role: message.role,
          });
        } catch {
          console.error('Failed to embed message:', message);
        }
        break;
      }
      default:
        throw new Error('Unknown job type');
    }
  }

  // @OnWorkerEvent('active')
  // onActive(job: Job) {
  //   console.log(
  //     `Processing job ${job.id} of type ${job.name} with data ${job.data}...`,
  //   );
  // }
  //
  // @OnWorkerEvent('completed')
  // onCompleted(job: Job) {
  //   console.log(
  //     `Processed job ${job.id} of type ${job.name} with data ${job.data}...`,
  //   );
  // }
  //
  // @OnWorkerEvent('failed')
  // onFailed(job: Job, error: Error) {
  //   console.log(
  //     `Processing job ${job.id} of type ${job.name} with data ${job.data}: ${error.message}`,
  //   );
  // }
}
