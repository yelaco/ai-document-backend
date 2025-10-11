import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
import { ConfigService } from '@nestjs/config';
import CustomValidationPipe from './shared/pipes/validation.pipe';
import { RequestContextInterceptor } from './shared/interceptors/request-context.interceptor';
import { LoggingInterceptor } from './shared/interceptors/logging.interceptor';

async function bootstrap() {
  const app = await NestFactory.create(AppModule);
  const configService = app.get(ConfigService);

  app.enableCors();
  app.useGlobalPipes(new CustomValidationPipe());
  app.useGlobalInterceptors(
    new RequestContextInterceptor(),
    new LoggingInterceptor(),
  );
  app.setGlobalPrefix('api');

  const port = configService.get<number>('port', 3000);

  await app.listen(port);
}
bootstrap();
