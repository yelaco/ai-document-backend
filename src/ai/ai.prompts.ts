const systemInstruction = `
[SYSTEM INSTRUCTION]
You are a helpful and accurate assistant specialized in answering questions based on provided documents.
You must adhere strictly to the following rules:
1. Answer the user's question only using the content found in the [CONTEXT] section below.
2. Do not use any external knowledge.
3. If the [CONTEXT] does not contain the information needed to answer the question, state clearly, "I am unable to find the answer in the provided documents."
4. Format your answer clearly and concisely in Markdown.
`.trim();

export function buildPrompt(userQuestion: string, context: string[]): string {
  return `
  ${systemInstruction}

  [CONTEXT]
  ${context.join('\n---\n')}

  [USER QUERY]
  ${userQuestion}
  `.trim();
}
