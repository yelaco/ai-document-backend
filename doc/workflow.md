# Workflow

The full workflow of an application using **Retrieval-Augmented Generation (RAG)** with **Gemini** (for embeddings and generation) and **pgvector** (as the vector database) is divided into two main pipelines: the **Indexing Pipeline** (data ingestion) and the **Query Processing Pipeline** (runtime search and response generation).

## Indexing Pipeline (Data Ingestion)

This pipeline converts your source documents into a searchable format and stores them in the database.

1.  **Load Documents** 📄: Source documents (e.g., PDFs, text files, HTML) are loaded into the system.
2.  **Split into Chunks** ✂️: Large documents are broken down into smaller, manageable, and semantically coherent **text chunks** (e.g., 200-800 tokens with some overlap) to fit within the Language Model's (LLM's) context window and improve retrieval accuracy.
3.  **Generate Embeddings** 🧠: The **Gemini embedding model** (e.g., `text-embedding-004`) converts each text chunk into a high-dimensional **vector** (an embedding) that captures its semantic meaning.
4.  **Store in Vector Database** 💾: The text chunks, along with their corresponding **vector embeddings** and **metadata** (like source document name), are stored in **PostgreSQL** using the **pgvector** extension. An index (like HNSW) is typically created on the vector column for fast similarity search.

---

## Query Processing Pipeline (Runtime)

This pipeline executes when a user submits a question and retrieves the relevant context to generate an answer.

1.  **User Query** ❓: A user asks a natural language question.
2.  **Embed Query** 🔬: The user's question is converted into a **query vector** using the _exact same_ **Gemini embedding model** used in the Indexing Pipeline to ensure the vectors exist in the same semantic space.
3.  **Retrieve Relevant Chunks** 🔎: The query vector is used to perform a **vector similarity search** against the stored document embeddings in **pgvector**. This search quickly identifies and retrieves the **Top-K** (e.g., 3-5) most semantically similar text chunks from the knowledge base.
4.  **Construct Prompt** ✍️: The retrieved text chunks are combined with the original user question into a single, comprehensive **prompt**. A system instruction is typically added to tell the **Gemini LLM** to answer the question _only_ based on the provided context.
5.  **Generate Response** ✅: The constructed prompt is sent to the **Gemini LLM** (e.g., `gemini-1.5-flash`), which uses the provided context to generate a factual, grounded answer.
6.  **Deliver Answer** 🗣️: The final answer is returned to the user, often including **source attribution** pointing back to the original documents from the retrieved metadata.

This video provides an end-to-end guide on building a RAG pipeline, which covers all the core steps of the app's workflow. [Build a RAG Pipeline with Gemini Embeddings and Vector Search- A Deep Dive (Full Code)](https://medium.com/google-cloud/build-a-rag-pipeline-with-gemini-embeddings-and-vector-search-a-deep-dive-full-code-bcd521ad9e1c).
