import { Module } from '@nestjs/common';
import { DocumentsModule } from './documents/documents.module';
import { AiModule } from './ai/ai.module';
import { LoggingInterceptor } from './shared/interceptors/logging.interceptor';
import { TypeOrmModule } from '@nestjs/typeorm';
import { ConfigModule, ConfigService } from '@nestjs/config';
import { Document } from './documents/entities/document.entity';
import { UsersModule } from './users/users.module';
import { AuthModule } from './auth/auth.module';

@Module({
  imports: [
    ConfigModule,
    AiModule,
    DocumentsModule,
    UsersModule,
    AuthModule,
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
        synchronize: configService.get<string>('appEnv') === 'development',
        logging: configService.get<string>('appEnv') === 'development',
        entities: [Document],
      }),
    }),
  ],
  providers: [
    {
      provide: 'APP_INTERCEPTOR',
      useClass: LoggingInterceptor,
    },
  ],
})
export class AppModule {}
