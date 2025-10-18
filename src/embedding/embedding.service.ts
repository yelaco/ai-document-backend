import { Inject, Injectable, OnModuleInit } from '@nestjs/common';
import {
  ChromaClient,
  ChromaNotFoundError,
  type EmbeddingFunction,
} from 'chromadb';
import {
  CHAT_MESSAGE_SEARCH_TIME_IN_MINUTES,
  CHROMA_CLIENT,
  COLLECTION_NAME_CHAT_MESSAGES,
  COLLECTION_NAME_DOCUMENTS,
  DOCUMENT_EMBEDDER,
  QUERY_EMBEDDER,
} from './embedding.constant';
import { RecursiveCharacterTextSplitter } from '@langchain/textsplitters';
import { EmbedMessageDto } from './dto/embed-message.dto';
import { ChatRole } from '../chats/chats.enum';
import { SearchMessageDto } from './dto/search-message.dto';

@Injectable()
export class EmbeddingService implements OnModuleInit {
  constructor(
    @Inject(CHROMA_CLIENT)
    private readonly chroma: ChromaClient,
    @Inject(DOCUMENT_EMBEDDER)
    private readonly documentEmbedder: EmbeddingFunction,
    @Inject(QUERY_EMBEDDER)
    private readonly queryEmbedder: EmbeddingFunction,
  ) {}

  async onModuleInit() {
    await this.chroma.heartbeat();
  }

  async embedDocument(documentId: string, text: string) {
    const splitter = new RecursiveCharacterTextSplitter({
      chunkSize: 500,
      chunkOverlap: 100,
    });

    const chunks = await splitter.splitText(text);
    try {
      const collection = await this.chroma.getOrCreateCollection({
        name: COLLECTION_NAME_DOCUMENTS,
        metadata: { 'hnsw:space': 'cosine' },
        embeddingFunction: this.documentEmbedder,
      });

      await collection.add({
        ids: chunks.map((_, idx) => `doc_chunk_${Date.now()}_${idx}`),
        documents: chunks,
        metadatas: chunks.map(() => ({
          documentId: documentId,
          created: new Date().toISOString(),
        })),
      });
    } catch (error) {
      if (error instanceof Error) {
        throw Error('Failed to store embeddings in ChromaDB: ' + error.message);
      }
    }
  }

  async vectorSearchDocument(documentId: string, query: string, topK = 10) {
    try {
      const collection = await this.chroma.getCollection({
        name: COLLECTION_NAME_DOCUMENTS,
        embeddingFunction: this.queryEmbedder,
      });
      const result = await collection.query({
        queryTexts: [query],
        nResults: topK,
        include: ['documents'],
        where: { documentId: documentId },
      });
      const retrievedChunks = result.documents[0].filter((doc) => doc !== null);
      return retrievedChunks;
    } catch (error) {
      if (error instanceof ChromaNotFoundError) {
        return [];
      }
      throw error;
    }
  }

  async embedChatMessage(options: EmbedMessageDto) {
    const splitter = new RecursiveCharacterTextSplitter({
      chunkSize: 500,
      chunkOverlap: 100,
    });

    const chunks = await splitter.splitText(options.content);
    try {
      const collection = await this.chroma.getOrCreateCollection({
        name: COLLECTION_NAME_CHAT_MESSAGES,
        metadata: { 'hnsw:space': 'cosine' },
        embeddingFunction: this.documentEmbedder,
      });

      await collection.add({
        ids: chunks.map((_, idx) => `chat_chunk_${Date.now()}_${idx}`),
        documents: chunks,
        metadatas: chunks.map(() => ({
          chatId: options.chatId,
          role: options.role,
          timestamp: new Date(options.timestamp).getTime(),
        })),
      });
    } catch (error) {
      if (error instanceof Error) {
        throw Error('Failed to store embeddings in ChromaDB: ' + error.message);
      }
    }
  }

  async vectorSearchChatMessage(searchMessageDto: SearchMessageDto) {
    try {
      const collection = await this.chroma.getCollection({
        name: COLLECTION_NAME_CHAT_MESSAGES,
        embeddingFunction: this.queryEmbedder,
      });
      const result = await collection.query({
        queryTexts: [searchMessageDto.query],
        nResults: searchMessageDto.topK || 15,
        include: ['documents'],
        where: {
          $and: [
            { chatId: searchMessageDto.chatId },
            {
              role: {
                $in: [ChatRole.USER, ChatRole.ASSISTANT],
              },
            },
            {
              timestamp: {
                $gte: searchMessageDto.timestamp
                  ? new Date(searchMessageDto.timestamp).getTime()
                  : new Date().getTime() - CHAT_MESSAGE_SEARCH_TIME_IN_MINUTES,
              },
            }, // messages in the last 30 minutes
          ],
        },
      });
      const retrievedChunks = result.documents[0].filter((doc) => doc !== null);
      return retrievedChunks;
    } catch (error) {
      if (error instanceof ChromaNotFoundError) {
        return [];
      }
      throw error;
    }
  }
}
