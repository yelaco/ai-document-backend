import { Module } from '@nestjs/common';
import { ConfigModule, ConfigService } from '@nestjs/config';
import { GeminiAiService } from './gemini.service';
import { AI_SERVICE } from './ai.constants';
import { AnthropicAiService } from './anthropic.service';

@Module({
  imports: [ConfigModule],
  providers: [
    GeminiAiService,
    AnthropicAiService,
    {
      provide: AI_SERVICE,
      useFactory: (
        configService: ConfigService,
        geminiService: GeminiAiService,
      ) => {
        const provider = configService.get<string>('AI_PROVIDER');

        if (provider === 'GEMINI') {
          return geminiService;
        }
      },
      inject: [ConfigService, GeminiAiService, AnthropicAiService],
    },
  ],
  exports: [AI_SERVICE],
})
export class AiModule {}
