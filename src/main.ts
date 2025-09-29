import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
import { ConfigService } from '@nestjs/config';
import CustomValidationPipe from './shared/pipes/validation.pipe';

async function bootstrap() {
  const app = await NestFactory.create(AppModule);
  const configService = app.get(ConfigService);

  app.enableCors();
  app.useGlobalPipes(new CustomValidationPipe());
  app.setGlobalPrefix('api');

  const port = configService.get<number>('port', 3000);

  await app.listen(port);
}
bootstrap();
