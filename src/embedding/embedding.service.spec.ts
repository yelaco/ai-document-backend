import { Test, TestingModule } from '@nestjs/testing';
import { EmbeddingService } from './embedding.service';
import {
  CHROMA_CLIENT,
  DOCUMENT_EMBEDDER,
  QUERY_EMBEDDER,
} from './embedding.constant';

const mockCollection = {
  add: jest.fn().mockResolvedValue(true),
  query: jest.fn().mockResolvedValue({
    ids: [['doc1']],
    documents: [['The mocked document result']],
  }),
};

const mockChromaClient = {
  constructor: jest.fn(),
  getOrCreateCollection: jest.fn().mockResolvedValue(mockCollection),
  getCollection: jest.fn().mockResolvedValue(mockCollection),
  heartbeat: jest.fn().mockResolvedValue(1),
};

export const mockEmbeddingService = {
  embedDocument: jest.fn(),
  vectorSearchDocument: jest.fn(() => ['The mocked document result']),
};

const mockEmbedder = { embedDocuments: jest.fn() };

describe('EmbeddingService', () => {
  let service: EmbeddingService;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      providers: [
        EmbeddingService,
        {
          provide: CHROMA_CLIENT,
          useValue: mockChromaClient,
        },
        {
          provide: DOCUMENT_EMBEDDER,
          useValue: mockEmbedder,
        },
        {
          provide: QUERY_EMBEDDER,
          useValue: mockEmbedder,
        },
      ],
    }).compile();

    service = module.get<EmbeddingService>(EmbeddingService);
  });

  it('should be defined', () => {
    expect(service).toBeDefined();
  });
});
