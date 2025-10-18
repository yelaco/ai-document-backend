import { Module } from '@nestjs/common';
import { MessagesService } from './messages.service';
import { MessagesController } from './messages.controller';
import { BullModule } from '@nestjs/bullmq';
import { MessagesProcessor } from './messages.processor';
import { EmbeddingModule } from '../embedding/embedding.module';
import { TypeOrmModule } from '@nestjs/typeorm';
import { Message } from './entities/message.entity';

@Module({
  imports: [
    TypeOrmModule.forFeature([Message]),
    BullModule.registerQueue({ name: 'message' }),
    EmbeddingModule,
  ],
  controllers: [MessagesController],
  providers: [MessagesService, MessagesProcessor],
  exports: [MessagesService],
})
export class MessagesModule {}
