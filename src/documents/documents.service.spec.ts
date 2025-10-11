import { Test, TestingModule } from '@nestjs/testing';
import { DocumentsService } from './documents.service';
import { AI_SERVICE } from '../ai/ai.constants';
import { mockAiService } from '../ai/ai.service.spec';
import { getRepositoryToken } from '@nestjs/typeorm';
import { Document } from './entities/document.entity';
import { EmbeddingService } from '../embedding/embedding.service';
import { mockEmbeddingService } from '../embedding/embedding.service.spec';
import { RequestContext } from '../shared/interceptors/request-context.interceptor';

const mockDocumentRepository = {
  save: jest.fn(),
  findOne: jest.fn(),
  find: jest.fn(),
};

describe('DocumentsService', () => {
  let service: DocumentsService;

  const ctx = {
    userId: 'fd200dc2-09a4-4f0e-9553-f5a73e7b04f8',
  } as RequestContext;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      providers: [
        DocumentsService,
        {
          provide: getRepositoryToken(Document),
          useValue: mockDocumentRepository,
        },
        {
          provide: AI_SERVICE,
          useValue: mockAiService,
        },
        {
          provide: EmbeddingService,
          useValue: mockEmbeddingService,
        },
      ],
    }).compile();

    service = module.get<DocumentsService>(DocumentsService);
  });

  it('should be defined', () => {
    expect(service).toBeDefined();
  });

  describe('ask', () => {
    it('should return an observable', async () => {
      const result = await service.ask(ctx, '1', 'What is NestJS');
      const possibleStatuses = ['in-progress', 'completed'];

      result.subscribe((event) => {
        expect(possibleStatuses).toContain(event.data['status']);
      });
    });
  });
});
