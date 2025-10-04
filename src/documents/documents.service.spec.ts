import { Test, TestingModule } from '@nestjs/testing';
import { DocumentsService } from './documents.service';
import { AI_SERVICE } from '../ai/ai.constants';
import { mockAiService } from '../ai/ai.service.spec';
import { getRepositoryToken } from '@nestjs/typeorm';
import { Document } from './entities/document.entity';
import { CHROMA_CLIENT } from './documents.constants';

const mockDocumentRepository = {
  save: jest.fn(),
  findOne: jest.fn(),
  find: jest.fn(),
};

export const mockCollection = {
  add: jest.fn().mockResolvedValue(true),
  query: jest.fn().mockResolvedValue({
    ids: [['doc1']],
    documents: [['The mocked document result']],
  }),
};

export const mockChromaClient = {
  constructor: jest.fn(),
  getOrCreateCollection: jest.fn().mockResolvedValue(mockCollection),
  getCollection: jest.fn().mockResolvedValue(mockCollection),
  heartbeat: jest.fn().mockResolvedValue(1),
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
          provide: CHROMA_CLIENT,
          useValue: mockChromaClient,
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
