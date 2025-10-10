import { Module } from '@nestjs/common';
import { ConfigModule, ConfigService } from '@nestjs/config';
import { GeminiAiService } from './gemini.service';
import { AI_SERVICE, GEMINI_AI } from './ai.constants';
import { AnthropicAiService } from './anthropic.service';
import { GoogleGenAI } from '@google/genai';

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
    {
      provide: GEMINI_AI,
      inject: [ConfigService],
      useFactory: (configService: ConfigService) => {
        return new GoogleGenAI({
          apiKey: configService.get<string>('ai.geminiApiKey'),
        });
      },
    },
  ],
  exports: [AI_SERVICE],
})
export class AiModule {}
