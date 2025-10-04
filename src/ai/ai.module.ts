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
      inject: [ConfigService, GeminiAiService, AnthropicAiService],
      useFactory: (
        configService: ConfigService,
        geminiAiService: GeminiAiService,
        anthropicAiService: AnthropicAiService,
      ) => {
        const aiServiceType = configService.get<string>('ai.serviceType');

        if (aiServiceType === 'anthropic') {
          return anthropicAiService;
        }

        return geminiAiService;
      },
    },
  ],
  exports: [AI_SERVICE],
})
export class AiModule {}
