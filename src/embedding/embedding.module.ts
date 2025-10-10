import { Module } from '@nestjs/common';
import { EmbeddingService } from './embedding.service';
import { ConfigModule, ConfigService } from '@nestjs/config';
import {
  CHROMA_CLIENT,
  DOCUMENT_EMBEDDER,
  EMBEDDING_MODEL,
  QUERY_EMBEDDER,
} from './embedding.constant';
import { GoogleGeminiEmbeddingFunction } from '@chroma-core/google-gemini';
import { ChromaClient } from 'chromadb';
import { EmbeddingTaskType } from './embedding.enum';

@Module({
  imports: [ConfigModule],
  providers: [
    EmbeddingService,
    {
      provide: CHROMA_CLIENT,
      inject: [ConfigService],
      useFactory: (configService: ConfigService) => {
        return new ChromaClient({
          host: configService.get<string>('vectorDatabase.host'),
          port: configService.get<number>('vectorDatabase.port'),
        });
      },
    },
    {
      provide: DOCUMENT_EMBEDDER,
      useValue: new GoogleGeminiEmbeddingFunction({
        apiKey: process.env.GEMINI_API_KEY,
        modelName: EMBEDDING_MODEL,
        taskType: EmbeddingTaskType.RETRIEVAL_DOCUMENT,
      }),
    },
    {
      provide: QUERY_EMBEDDER,
      useValue: new GoogleGeminiEmbeddingFunction({
        apiKey: process.env.GEMINI_API_KEY,
        modelName: EMBEDDING_MODEL,
        taskType: EmbeddingTaskType.RETRIEVAL_QUERY,
      }),
    },
  ],
  exports: [EmbeddingService],
})
export class EmbeddingModule {}
