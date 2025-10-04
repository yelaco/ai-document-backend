import { of } from 'rxjs';

export const mockAiService = {
  generateEmbedding: jest.fn((texts: string[]) =>
    Promise.resolve(Array<number[]>(texts.length).fill([1, 1, 1])),
  ),
  answer: jest.fn((question: string) => Promise.resolve(of(question))),
};
