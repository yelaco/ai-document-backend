import { Column, PrimaryGeneratedColumn } from 'typeorm';

export class Message {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column()
  role: string;

  @Column()
  content: string;

  @Column({ name: 'chat_id' })
  chatId: string;

  @Column({ type: 'timestamptz', default: () => 'CURRENT_TIMESTAMP' })
  createdAt: Date;
}
