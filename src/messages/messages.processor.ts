import { OnWorkerEvent, Processor, WorkerHost } from '@nestjs/bullmq';
import { Job } from 'bullmq';
import { CreateMessageDto } from './dto/create-message.dto';

@Processor('messages')
export class MessagesProcessor extends WorkerHost {
  async process(job: Job): Promise<any> {
    switch (job.name) {
      case 'handle-chat-message': {
        const message = job.data as CreateMessageDto;
        console.log('Processing message:', message);

        await new Promise((resolve) => setTimeout(resolve, 5000));
        return 'Message processed successfully';
      }
      default:
        throw new Error('Unknown job type');
    }
  }

  @OnWorkerEvent('completed')
  onCompleted(job: Job) {
    console.log(
      `Processing job ${job.id} of type ${job.name} with data ${job.data}...`,
    );
  }

  @OnWorkerEvent('failed')
  onFailed(job: Job, error: Error) {
    console.log(
      `Processing job ${job.id} of type ${job.name} with data ${job.data}...`,
    );
  }
}
