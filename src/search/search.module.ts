import { Module } from '@nestjs/common';
import { AiModule } from 'src/ai/ai.module';
import { SearchService } from './search.service';
import { SearchController } from './search.controller';

@Module({
  imports: [AiModule],
  controllers: [SearchController],
  providers: [SearchService],
})
export class SearchModule {}
