import { Expose, instanceToPlain, plainToInstance } from 'class-transformer';
import { Document } from '../entities/document.entity';

import { ApiProperty } from '@nestjs/swagger';

export class DocumentDto {
  @ApiProperty({ example: 'doc123' })
  @Expose()
  id: string;

  @ApiProperty({ example: 'Project Proposal.pdf' })
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
