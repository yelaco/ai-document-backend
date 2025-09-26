import { Module } from '@nestjs/common';
import { CoreModule } from './core/core.module';
import { DocumentsModule } from './documents/documents.module';
import { QueueModule } from './queue/queue.module';
import { SearchModule } from './search/search.module';
import { WebsocketModule } from './websocket/websocket.module';
import { TypeOrmModule } from '@nestjs/typeorm';
import { ConfigModule, ConfigService } from '@nestjs/config';
import { Document } from './documents/entities/document.entity';
import { AiModule } from './ai/ai.module';

@Module({
  imports: [
    CoreModule,
    DocumentsModule,
    QueueModule,
    SearchModule,
    WebsocketModule,
    TypeOrmModule.forRootAsync({
      imports: [ConfigModule],
      inject: [ConfigService],
      useFactory: (configService: ConfigService) => ({
        type: 'postgres',
        host: configService.get<string>('database.host'),
        post: configService.get<number>('database.port'),
        database: configService.get<string>('database.name'),
        username: configService.get<string>('database.user'),
        password: configService.get<string>('database.pass'),
        entities: [Document],
      }),
    }),
    AiModule,
  ],
})
export class AppModule {}
