import { User } from '../../users/entities/user.entity';
import { Document } from '../../documents/entities/document.entity';
import {
  Column,
  CreateDateColumn,
  Entity,
  JoinColumn,
  ManyToOne,
  OneToMany,
  PrimaryGeneratedColumn,
  UpdateDateColumn,
  type Relation,
} from 'typeorm';
import { Message } from 'src/messages/entities/message.entity';

@Entity()
export class Chat {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column({ nullable: true })
  title?: string;

  @CreateDateColumn()
  createdAt: Date;

  @UpdateDateColumn()
  updatedAt: Date;

  @Column({ name: 'document_id' })
  documentId: string;

  @ManyToOne(() => Document, (document) => document.chats, {
    cascade: true,
    onDelete: 'CASCADE',
  })
  @JoinColumn({ name: 'document_id' })
  document: Relation<Document>;

  @Column({ name: 'user_id' })
  userId: string;

  @ManyToOne(() => User, (user) => user.chats, {
    cascade: true,
    onDelete: 'CASCADE',
  })
  @JoinColumn({ name: 'user_id' })
  user: Relation<User>;

  @OneToMany(() => Message, (message) => message.chat)
  messages: Relation<Message>[];
}
