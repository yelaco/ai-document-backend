import {
  Inject,
  Injectable,
  NotFoundException,
  UnauthorizedException,
} from '@nestjs/common';
import { CreateChatDto } from './dto/create-chat.dto';
import { UpdateChatDto } from './dto/update-chat.dto';
import { Repository } from 'typeorm';
import { Chat } from './entities/chat.entity';
import { InjectRepository } from '@nestjs/typeorm';
import { RequestContext } from '../shared/interceptors/request-context.interceptor';
import { endWith, map, Observable, tap } from 'rxjs';
import { AI_SERVICE } from '../ai/ai.constants';
import { buildPrompt } from '../ai/ai.prompts';
import { type AiService } from '../ai/ai.service';
import { AnswerDocumentDto } from '../chats/dto/ask-document.dto';
import { EmbeddingService } from '../embedding/embedding.service';
import { MessagesService } from '../messages/messages.service';
import { ChatResponseProgress, ChatRole } from './chats.enum';

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
    const chat = await this.findOne(ctx, chatId);
    if (!chat) {
      throw new NotFoundException('Chat not found');
    }

    // Save user message
    this.messageService
      .process({
        chatId: chat.id,
        content: question,
        role: ChatRole.USER,
      })
      .catch((err) => {
        console.error('Failed to process message:', err);
      });

    const questionContext = await this.embeddingService.vectorSearchDocument(
      chat.documentId,
      question,
    );

    const observableAnswer = await this.aiService.answer(
      buildPrompt(question, questionContext),
    );

    return observableAnswer.pipe(
      tap((aiResp) => {
        this.messageService
          .process({
            chatId: chat.id,
            content: aiResp.answer,
            role: ChatRole.ASSISTANT,
          })
          .catch((err) => {
            console.error('Failed to process message:', err);
          });
      }),
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
  ): Promise<{ data: Chat[]; total: number }> {
    const [data, total] = await this.chatRepo.findAndCount({
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
    });
  }

  async remove(ctx: RequestContext, id: string): Promise<void> {
    await this.chatRepo.delete(id);
  }
}
