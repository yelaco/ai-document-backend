import { User } from '../../users/entities/user.entity';
import { Document } from '../../documents/entities/document.entity';
import {
  Column,
  Entity,
  JoinColumn,
  ManyToOne,
  PrimaryGeneratedColumn,
} from 'typeorm';

@Entity()
export class Chat {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column()
  title: string;

  @Column({ type: 'timestamptz', default: () => 'CURRENT_TIMESTAMP' })
  createdAt: Date;

  @ManyToOne(() => Document, (document) => document.chats, {
    cascade: true,
    onDelete: 'CASCADE',
  })
  @JoinColumn({ name: 'document_id' })
  document: typeof Document;

  @Column({ name: 'document_id' })
  documentId: string;

  @Column({ name: 'user_id' })
  userId: string;

  @ManyToOne(() => User, (user) => user.chats, {
    cascade: true,
    onDelete: 'CASCADE',
  })
  @JoinColumn({ name: 'user_id' })
  user: typeof User;
}
