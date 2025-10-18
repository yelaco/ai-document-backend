export class PromptDto {
  documentTitle: string;
  query: string;
  context: string[];
  pastConservations?: string[];
}
