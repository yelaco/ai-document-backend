import { Module } from '@nestjs/common';
import { DocumentsModule } from './documents/documents.module';
import { AiModule } from './ai/ai.module';
import { TypeOrmModule } from '@nestjs/typeorm';
import { ConfigModule, ConfigService } from '@nestjs/config';
import { Document } from './documents/entities/document.entity';
import { UsersModule } from './users/users.module';
import { AuthModule } from './auth/auth.module';
import { User } from './users/entities/user.entity';
import { EmbeddingModule } from './embedding/embedding.module';
import Joi from 'joi';
import envConfig from './config/env.config';

@Module({
  imports: [
    ConfigModule.forRoot({
      envFilePath: [`.env`],
      load: [envConfig],
      cache: true,
      validationSchema: Joi.object({
        APP_PORT: Joi.number().port().default(3000),
        DB_HOST: Joi.string().required(),
        DB_PORT: Joi.number().port().default(5432),
        DB_NAME: Joi.string().required(),
        DB_USER: Joi.string().required(),
        DB_PASS: Joi.string().optional(),
      }),
    }),
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
        entities: [Document, User],
      }),
    }),
    AiModule,
    DocumentsModule,
    UsersModule,
    AuthModule,
    EmbeddingModule,
  ],
})
export class AppModule {}
