import { Exclude } from 'class-transformer';
import { Document } from '../../documents/entities/document.entity';
import { Column, Entity, OneToMany, PrimaryGeneratedColumn } from 'typeorm';
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

  @OneToMany(() => Document, (document) => document.user)
  documents: Document[];

  @OneToMany(() => Chat, (chat) => chat.user)
  chats: Chat[];
}
