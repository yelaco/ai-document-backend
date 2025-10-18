import {
  Inject,
  Injectable,
  InternalServerErrorException,
  NotFoundException,
  UnauthorizedException,
} from '@nestjs/common';
import { CreateChatDto } from './dto/create-chat.dto';
import { UpdateChatDto } from './dto/update-chat.dto';
import { Repository } from 'typeorm';
import { Chat } from './entities/chat.entity';
import { InjectRepository } from '@nestjs/typeorm';
import { RequestContext } from '../shared/interceptors/request-context.interceptor';
import { endWith, map, Observable, shareReplay } from 'rxjs';
import { AI_SERVICE } from '../ai/ai.constants';
import { buildPrompt } from '../ai/ai.prompts';
import { type AiService } from '../ai/ai.service';
import { AnswerDocumentDto } from '../chats/dto/ask-document.dto';
import { EmbeddingService } from '../embedding/embedding.service';
import { MessagesService } from '../messages/messages.service';
import { ChatResponseProgress, ChatRole } from './chats.enum';
import { CreateMessageDto } from '../messages/dto/create-message.dto';

@Injectable()
export class ChatsService {
  constructor(
    @InjectRepository(Chat)
    private chatRepo: Repository<Chat>,
    private readonly embeddingService: EmbeddingService,
    private readonly messageService: MessagesService,
    @Inject(AI_SERVICE)
    private readonly aiService: AiService,
  ) {}

  async ask(
    ctx: RequestContext,
    chatId: string,
    question: string,
  ): Promise<Observable<AnswerDocumentDto>> {
    const chat = await this.chatRepo.findOne({
      where: { id: chatId, userId: ctx.userId },
      relations: ['document'],
      cache: {
        id: `chat-${chatId}-${ctx.userId}`,
        milliseconds: 30000,
      },
    });
    if (!chat) {
      throw new NotFoundException('Chat not found');
    }

    // Save user message
    this.messageService
      .process(ctx, {
        chatId: chat.id,
        content: question,
        role: ChatRole.USER,
        timestamp: new Date(),
      })
      .catch((err) => {
        console.error('Failed to process message:', err);
      });

    const questionContext = await this.embeddingService.vectorSearchDocument(
      chat.documentId,
      question,
    );

    const chatMessages = await this.embeddingService.vectorSearchChatMessage({
      chatId: chat.id,
      query: question,
      // TODO: pick timestamp of the previous Nth message,
    });

    const answerObs = await this.aiService.answer(
      buildPrompt({
        documentTitle: chat.document.title,
        query: question,
        context: questionContext,
        pastConservations: chatMessages,
      }),
    );

    const replayAnswerObs = answerObs.pipe(
      shareReplay({ bufferSize: 1, refCount: true }),
    );

    const createMessageDto: CreateMessageDto = {
      chatId: chat.id,
      content: '',
      role: ChatRole.ASSISTANT,
      timestamp: new Date(),
    };

    replayAnswerObs.subscribe({
      next: (aiResp) => {
        createMessageDto.content += aiResp.answer;

        if (createMessageDto.content.length > 2000) {
          createMessageDto.timestamp = new Date();

          this.messageService.process(ctx, createMessageDto).catch((err) => {
            console.error('Failed to process message:', err);
          });

          createMessageDto.content = '';
        }
      },
      complete: () => {
        if (createMessageDto.content.length > 0) {
          createMessageDto.timestamp = new Date();

          this.messageService.process(ctx, createMessageDto).catch((err) => {
            console.error('Failed to process message:', err);
          });
        }
      },
      error: (err) => {
        console.error('Error in AI answer observable:', err);
      },
    });

    return replayAnswerObs.pipe(
      map(
        (text) =>
          ({
            text: text.answer,
            status: ChatResponseProgress.IN_PROGRESS,
          }) as AnswerDocumentDto,
      ),
      endWith({
        status: ChatResponseProgress.COMPLETED,
      } as AnswerDocumentDto),
    );
  }

  async create(ctx: RequestContext, createChatDto: CreateChatDto) {
    if (!ctx.userId) {
      throw new UnauthorizedException();
    }
    const chat = this.chatRepo.create({
      title: createChatDto.title,
      documentId: createChatDto.documentId,
      userId: ctx.userId,
    });
    return this.chatRepo.save(chat);
  }

  async findPaginated(
    ctx: RequestContext,
    page: number,
    pageSize: number,
    documentId?: string,
  ): Promise<{ data: Chat[]; total: number }> {
    if (documentId) {
      const exist = await this.chatRepo.exists({
        where: { documentId: documentId, userId: ctx.userId },
      });
      if (!exist) {
        throw new NotFoundException('Document not found');
      }
    }

    const [data, total] = await this.chatRepo.findAndCount({
      where: { userId: ctx.userId, documentId: documentId },
      skip: (page - 1) * pageSize,
      take: pageSize,
      order: { createdAt: 'DESC' },
    });
    return { data, total };
  }

  async update(ctx: RequestContext, id: string, updateChatDto: UpdateChatDto) {
    const chat = await this.findOne(ctx, id);
    if (!chat) {
      throw new NotFoundException('Chat not found');
    }
    chat.title = updateChatDto.title ?? chat.title;
    return this.chatRepo.save(chat);
  }

  async findOne(ctx: RequestContext, id: string): Promise<Chat | null> {
    if (!ctx.userId) {
      throw new UnauthorizedException();
    }

    return this.chatRepo.findOne({
      where: { id: id, userId: ctx.userId },
      relations: ['messages'],
    });
  }

  async remove(ctx: RequestContext, id: string): Promise<void> {
    try {
      await this.chatRepo.delete(id);
    } catch {
      throw new InternalServerErrorException();
    }
  }
}
