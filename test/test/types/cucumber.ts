import { setWorldConstructor } from "@cucumber/cucumber";
import dotenv from "dotenv";

type DraftRequest = {
  url?: string;
  options: {
    method?: string;
    headers: Record<string, string>;
    body?: string | FormData;
  };
};
type ViolationResult = {
  status: number;
  body?: Record<string, unknown>;
};

export class ApiWorld {
  private static instance: ApiWorld;

  baseUrl: string;
  lastResponse?: Response;
  jsonBody?: Record<string, unknown>;

  request: DraftRequest;
  invalidPayloads?: unknown[];
  violationResponses: ViolationResult[];
  savedPayload: Record<string, unknown>;
  lastSchemaName?: string;

  private constructor() {
    dotenv.config();

    this.baseUrl = process.env["API_URL"] || "[https://example.com](https://example.com)";
    this.request = { options: { headers: {} } };
    this.violationResponses = [];
    this.savedPayload = {};
  }

  // Singleton getter
  static getInstance() {
    if (!ApiWorld.instance) {
      ApiWorld.instance = new ApiWorld();
    }
    return ApiWorld.instance;
  }

  resetRequest() {
    this.request = { options: { headers: {} } };
    this.invalidPayloads = [];
    this.violationResponses = [];
    this.savedPayload = {};
  }
}

setWorldConstructor(() => ApiWorld.getInstance());

