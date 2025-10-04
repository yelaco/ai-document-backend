import { Module } from '@nestjs/common';
import { DocumentsService } from './documents.service';
import { DocumentsController } from './documents.controller';
import { TypeOrmModule } from '@nestjs/typeorm';
import { Document } from './entities/document.entity';
import { AiModule } from '../ai/ai.module';
import { ConfigModule, ConfigService } from '@nestjs/config';
import { ChromaClient } from 'chromadb';
import { CHROMA_CLIENT } from './documents.constants';

@Module({
  imports: [ConfigModule, TypeOrmModule.forFeature([Document]), AiModule],
  controllers: [DocumentsController],
  providers: [
    DocumentsService,
    {
      provide: CHROMA_CLIENT,
      useFactory: (configService: ConfigService) => {
        return new ChromaClient({
          host: configService.get<string>('vectorDatabase.host'),
          port: configService.get<number>('vectorDatabase.port'),
        });
      },
      inject: [ConfigService],
    },
  ],
})
export class DocumentsModule {}
