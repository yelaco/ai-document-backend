import { Test, TestingModule } from '@nestjs/testing';
import { ChatsService } from './chats.service';
import { RequestContext } from '../shared/interceptors/request-context.interceptor';
import { mockAiService } from '../ai/ai.service.spec';
import { AI_SERVICE } from '../ai/ai.constants';
import { getRepositoryToken } from '@nestjs/typeorm';
import { Chat } from './entities/chat.entity';
import { EmbeddingService } from '../embedding/embedding.service';
import { mockEmbeddingService } from '../embedding/embedding.service.spec';

describe('ChatsService', () => {
  let service: ChatsService;

  const ctx = {
    userId: 'fd200dc2-09a4-4f0e-9553-f5a73e7b04f8',
  } as RequestContext;

  const mockChat = {
    id: '3df4b453-d7d3-4b5d-87cf-4466220dc8b2',
    title: 'Hello World',
    documentId: '43c3d501-6071-47ff-ac79-1272e366657f',
  } as Chat;

  const mockChatRepository = {
    create: jest.fn().mockImplementation(() => mockChat),
    save: jest.fn().mockImplementation(() => Promise.resolve(mockChat)),
    findOne: jest.fn().mockImplementation(() => mockChat),
    find: jest.fn(),
  };

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      providers: [
        ChatsService,
        {
          provide: getRepositoryToken(Chat),
          useValue: mockChatRepository,
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

    service = module.get<ChatsService>(ChatsService);
  });

  it('should be defined', () => {
    expect(service).toBeDefined();
  });

  describe('ask', () => {
    it('should return an observable', async () => {
      const chat = await service.create(ctx, {
        title: mockChat.title,
        documentId: mockChat.documentId,
      });
      const result = await service.ask(ctx, chat.id, 'What is NestJS');
      const possibleStatuses = ['in-progress', 'completed'];

      result.subscribe((event) => {
        expect(possibleStatuses).toContain(event.status);
      });
    });
  });
});
