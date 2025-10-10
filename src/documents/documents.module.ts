import { Module } from '@nestjs/common';
import { DocumentsService } from './documents.service';
import { DocumentsController } from './documents.controller';
import { TypeOrmModule } from '@nestjs/typeorm';
import { Document } from './entities/document.entity';
import { AiModule } from '../ai/ai.module';
import { ConfigModule } from '@nestjs/config';
import { EmbeddingModule } from 'src/embedding/embedding.module';

@Module({
  imports: [
    ConfigModule,
    TypeOrmModule.forFeature([Document]),
    AiModule,
    EmbeddingModule,
  ],
  controllers: [DocumentsController],
  providers: [DocumentsService],
})
export class DocumentsModule {}
