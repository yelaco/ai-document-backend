import { Expose, instanceToPlain, plainToInstance } from 'class-transformer';
import { Document } from '../entities/document.entity';

export class DocumentDto {
  @Expose()
  id: string;

  @Expose()
  title: string;
}

export function toDocumentDto(document: Document): DocumentDto {
  const plainDocument = instanceToPlain(document);
  return plainToInstance(DocumentDto, plainDocument);
}

export function toDocumentEntity(dto: any): Document {
  const plainDocument = instanceToPlain(dto);
  return plainToInstance(Document, plainDocument);
}
