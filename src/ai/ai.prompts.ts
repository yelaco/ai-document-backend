import { PromptDto } from './dto/prompt.dto';

const systemInstruction = `
[SYSTEM INSTRUCTION]
You are a helpful and accurate assistant specialized in answering questions based on provided documents.
You must adhere strictly to the following rules:
1. Answer the user's question using the content found in the [CONTEXT] section below and external knowledge.
2. If the document information is wrong or contradicts common knowledge, you should point out where it's wrong.
3. If the question does not related at all to the [CONTEXT],
state clearly, "I am unable to find the answer in the provided documents." and explain more on why you can't answer.
4. Format your answer clearly and concisely in Markdown.
`.trim();

export function buildPrompt(promptDto: PromptDto): string {
  return `
  ${systemInstruction}

  [DOCUMENT_METADATA]
  - Title: ${promptDto.documentTitle}

  [CONSERVATION_HISTORY]
  ${promptDto.pastConservations ? promptDto.pastConservations.join('\n---\n') : 'No prior conversations.'}

  [CONTEXT]
  ${promptDto.context.join('\n---\n')}

  [USER QUERY]
  ${promptDto.query}
  `.trim();
}
