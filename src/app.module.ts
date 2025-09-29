import { Module } from '@nestjs/common';
import { CoreModule } from './core/core.module';
import { DocumentsModule } from './documents/documents.module';
import { QueueModule } from './queue/queue.module';
import { SearchModule } from './search/search.module';
import { WebsocketModule } from './websocket/websocket.module';
import { AiModule } from './ai/ai.module';
import { LoggingInterceptor } from './shared/interceptors/logging.interceptor';

@Module({
  imports: [
    CoreModule,
    DocumentsModule,
    QueueModule,
    SearchModule,
    WebsocketModule,
    AiModule,
  ],
  providers: [
    {
      provide: 'APP_INTERCEPTOR',
      useClass: LoggingInterceptor,
    },
  ],
})
export class AppModule {}
