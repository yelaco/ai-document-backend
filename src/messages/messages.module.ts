import { Module } from '@nestjs/common';
import { MessagesService } from './messages.service';
import { MessagesController } from './messages.controller';
import { BullModule } from '@nestjs/bullmq';
import { MessagesProcessor } from './messages.processor';

@Module({
  imports: [BullModule.registerQueue({ name: 'message' })],
  controllers: [MessagesController],
  providers: [MessagesService, MessagesProcessor],
})
export class MessagesModule {}
