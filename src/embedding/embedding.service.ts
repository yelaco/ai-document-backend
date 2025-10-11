import { Inject, Injectable, OnModuleInit } from '@nestjs/common';
import {
  ChromaClient,
  ChromaNotFoundError,
  type EmbeddingFunction,
} from 'chromadb';
import {
  CHROMA_CLIENT,
  DOCUMENT_EMBEDDER,
  QUERY_EMBEDDER,
} from './embedding.constant';
import { RecursiveCharacterTextSplitter } from '@langchain/textsplitters';

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
        name: 'documents',
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
        name: 'documents',
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
}
