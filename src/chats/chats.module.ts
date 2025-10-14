import { Module } from '@nestjs/common';
import { ChatsService } from './chats.service';
import { ChatsController } from './chats.controller';
import { ConfigModule } from '@nestjs/config';
import { TypeOrmModule } from '@nestjs/typeorm';
import { Chat } from './entities/chat.entity';

@Module({
  imports: [ConfigModule, TypeOrmModule.forFeature([Chat])],
  controllers: [ChatsController],
  providers: [ChatsService],
})
export class ChatsModule {}
