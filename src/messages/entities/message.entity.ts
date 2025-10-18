import { Chat } from 'src/chats/entities/chat.entity';
import {
  Column,
  CreateDateColumn,
  Entity,
  JoinColumn,
  ManyToOne,
  PrimaryGeneratedColumn,
  type Relation,
  UpdateDateColumn,
} from 'typeorm';

@Entity()
export class Message {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column()
  role: string;

  @Column()
  content: string;

  @Column({ type: 'timestamptz' })
  timestamp: Date;

  @CreateDateColumn()
  createdAt: Date;

  @UpdateDateColumn()
  updatedAt: Date;

  @Column({ name: 'chat_id' })
  chatId: string;

  @ManyToOne(() => Chat, (chat) => chat.messages, {
    cascade: true,
    onDelete: 'CASCADE',
  })
  @JoinColumn({ name: 'chat_id' })
  chat: Relation<Chat>;
}
