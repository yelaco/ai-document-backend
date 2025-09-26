import { Injectable } from '@nestjs/common';
import { CreateDocumentDto } from './dto/create-document.dto';
import { InjectRepository } from '@nestjs/typeorm';
import { Document } from './entities/document.entity';
import { Repository } from 'typeorm';

@Injectable()
export class DocumentsService {
  constructor(
    @InjectRepository(Document)
    private documentsRepository: Repository<Document>,
  ) {}

  create(createDocumentDto: CreateDocumentDto) {
    return 'This action adds a new document';
  }

  findAll(): Promise<Document[]> {
    return this.documentsRepository.find();
  }

  findOne(id: number): Promise<Document | null> {
    return this.documentsRepository.findOneBy({ id });
  }

  async remove(id: number): Promise<void> {
    await this.documentsRepository.delete(id);
  }
}
