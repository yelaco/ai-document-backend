import { Test, TestingModule } from '@nestjs/testing';
import { DocumentsService } from './documents.service';
import { AI_SERVICE } from '../ai/ai.constants';
import { mockAiService } from '../ai/ai.service.spec';
import { getRepositoryToken } from '@nestjs/typeorm';
import { Document } from './entities/document.entity';
import { EmbeddingService } from 'src/embedding/embedding.service';
import { mockEmbeddingService } from 'src/embedding/embedding.service.spec';

const mockDocumentRepository = {
  save: jest.fn(),
  findOne: jest.fn(),
  find: jest.fn(),
};

describe('DocumentsService', () => {
  let service: DocumentsService;

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
      const result = await service.ask({
        documentId: '1',
        question: 'What is NestJS',
      });
      const possibleStatuses = ['in-progress', 'completed'];

      result.subscribe((event) => {
        expect(possibleStatuses).toContain(event.data['status']);
      });
    });
  });
});
