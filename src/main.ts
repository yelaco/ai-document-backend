import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
import { ConfigService } from '@nestjs/config';
import CustomValidationPipe from './shared/pipes/validation.pipe';
import { RequestContextInterceptor } from './shared/interceptors/request-context.interceptor';
import { LoggingInterceptor } from './shared/interceptors/logging.interceptor';
import { DocumentBuilder, SwaggerModule } from '@nestjs/swagger';

async function bootstrap() {
  const app = await NestFactory.create(AppModule);
  const configService = app.get(ConfigService);

  // Global middleware and settings
  app.enableCors();
  app.useGlobalPipes(new CustomValidationPipe());
  app.useGlobalInterceptors(
    new RequestContextInterceptor(),
    new LoggingInterceptor(),
  );
  app.setGlobalPrefix('api');

  // Swagger setup
  const config = new DocumentBuilder()
    .setTitle('AI Document API')
    .setDescription('API documentation for the AI Document application')
    .setVersion('1.0')
    .addTag('ai-document')
    .build();
  const documentFactory = () => SwaggerModule.createDocument(app, config);
  SwaggerModule.setup('swagger', app, documentFactory);

  // Start the application
  await app.listen(configService.get<number>('port', 3000));
}
bootstrap();
