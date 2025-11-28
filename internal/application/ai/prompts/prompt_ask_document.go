package prompts

import "strings"

const systemInstructionAskDocument = `
[SYSTEM INSTRUCTION]
You are a helpful and accurate assistant specialized in answering questions based on provided documents.
You must adhere strictly to the following rules:
1. Answer the user's question using the content found in the [CONTEXT] section below and external knowledge.
2. If the document information is wrong or contradicts common knowledge, you should point out where it's wrong.
3. If the question does not related at all to the [CONTEXT],
state clearly, "I am unable to find the answer in the provided documents." and explain more on why you can't answer.
4. Format your answer clearly and concisely in Markdown.
`

func BuildAskDocumentPrompt(context []string, question string, pastConversations []string, metadata map[string]string) string {
	sb := strings.Builder{}
	sb.WriteString(systemInstructionAskDocument)
	sb.WriteString("\n[CONTEXT]\n")
	for _, doc := range context {
		sb.WriteString(doc)
		sb.WriteString("\n---\n")
	}
	if len(pastConversations) > 0 {
		sb.WriteString("\n[PAST CONVERSATIONS]\n")
		for _, conv := range pastConversations {
			sb.WriteString(conv)
			sb.WriteString("\n")
		}
	}
	sb.WriteString("\n[METADATA]\n")
	for key, value := range metadata {
		sb.WriteString(key)
		sb.WriteString(": ")
		sb.WriteString(value)
		sb.WriteString("\n")
	}
	sb.WriteString("\n[USER QUESTION]\n")
	sb.WriteString(question)
	return sb.String()
}
