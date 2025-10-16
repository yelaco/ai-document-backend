import { Module } from '@nestjs/common';
import { ChatsService } from './chats.service';
import { ChatsController } from './chats.controller';
import { ConfigModule } from '@nestjs/config';
import { TypeOrmModule } from '@nestjs/typeorm';
import { Chat } from './entities/chat.entity';
import { AiModule } from '../ai/ai.module';
import { EmbeddingModule } from '../embedding/embedding.module';
import { MessagesModule } from '../messages/messages.module';

@Module({
  imports: [
    ConfigModule,
    TypeOrmModule.forFeature([Chat]),
    AiModule,
    EmbeddingModule,
    MessagesModule,
  ],
  controllers: [ChatsController],
  providers: [ChatsService],
})
export class ChatsModule {}
