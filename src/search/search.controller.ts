import { Controller, Get, Query } from '@nestjs/common';
import { SearchService } from './search.service';

@Controller('search')
export class SearchController {
  constructor(private readonly searchService: SearchService) {}

  @Get()
  async executeSearch(@Query('q') query: string) {
    const results = await this.searchService.performSemanticSearch(query);

    return {
      query: query,
      results: results,
    };
  }
}
