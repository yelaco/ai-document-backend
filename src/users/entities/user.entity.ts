import { Exclude } from 'class-transformer';
import { Document } from '../../documents/entities/document.entity';
import {
  Column,
  CreateDateColumn,
  Entity,
  OneToMany,
  PrimaryGeneratedColumn,
  UpdateDateColumn,
  type Relation,
} from 'typeorm';
import { Chat } from '../../chats/entities/chat.entity';

@Entity()
export class User {
  @PrimaryGeneratedColumn('uuid')
  id: number;

  @Column()
  email: string;

  @Exclude({ toPlainOnly: true })
  @Column()
  passwordHash: string;

  @Column()
  fullName: string;

  @CreateDateColumn()
  createdAt: Date;

  @UpdateDateColumn()
  updatedAt: Date;

  @OneToMany(() => Document, (document) => document.user)
  documents: Relation<Document>[];

  @OneToMany(() => Chat, (chat) => chat.user)
  chats: Relation<Chat>[];
}
