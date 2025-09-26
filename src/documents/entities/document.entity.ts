import { Column, PrimaryGeneratedColumn } from 'typeorm';

export class Document {
  @PrimaryGeneratedColumn()
  id: number;

  @Column()
  title: string;

  @Column('text')
  content: string;

  @Column('jsonb')
  vector: number[];
}
