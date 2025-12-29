import {
  When, Then
} from "@cucumber/cucumber";
import { payloads } from "../payloads";
import { resolvePayload } from "../support/helper";
import type { ApiWorld } from "../types/cucumber";

import crypto from "crypto";
import fs from "fs";
import path from "path";

// Add payload
When(
  "I use formdata {word} payload",
  async function(this: ApiWorld, payloadName: string) {
    const payload = payloads[payloadName];
    if (!payload) {
      throw new Error(`Payload '${payloadName}' not found`);
    }

    const resolved = resolvePayload(payload.positive, this);

    const form = new FormData();

    for (const [key, value] of Object.entries(resolved)) {
      if (key === "path") continue;
      form.append(key, String(value));
    }

    const filePath = path.resolve(
      __dirname,
      "../test-images",
      resolved["path"] as string
    );

    if (!fs.existsSync(filePath)) {
      throw new Error(`File not found at: ${filePath}`);
    }

    const buffer = fs.readFileSync(filePath);
    const blob = new Blob([buffer], { type: "image/jpeg" });

    form.append("file", blob, path.basename(filePath));

    const headers = { ...(this.request.options.headers || {}) };
    delete headers["Content-Type"];
    delete headers["content-type"];

    this.request.options = {
      ...this.request.options,
      body: form,
      headers,
    };
  }
);

// Add violations
When(
  "I inject formdata {word} violations",
  async function(this: ApiWorld, entity: string) {
    const payload = payloads[entity];
    if (!payload) throw new Error(`Violation '${entity}' not found`);
    this.invalidPayloads = payload.negative;
  },
);

function hashBuffer(buf: Buffer): string {
  return crypto.createHash("sha256").update(buf).digest("hex");
}

function hashFile(path: string): string {
  return crypto.createHash("sha256").update(fs.readFileSync(path)).digest("hex");
}

Then(
  "the file should be equal as {word}",
  async function(this: ApiWorld, payloadName: string) {
    if (!this.lastResponse) {
      throw new Error("lastResponse not found");
    }

    // ambil binary dari response GET file
    const responseBuffer = Buffer.from(
      await this.lastResponse.arrayBuffer()
    );

    const payload = payloads[payloadName];
    if (!payload) {
      throw new Error(`Payload '${payloadName}' not found`);
    }

    const resolved = resolvePayload(payload.positive, this);

    const form = new FormData();

    for (const [key, value] of Object.entries(resolved)) {
      if (key === "path") continue;
      form.append(key, String(value));
    }

    const filePath = path.resolve(
      __dirname,
      "../test-images",
      resolved["path"] as string
    );

    if (!fs.existsSync(filePath)) {
      throw new Error(`File not found at: ${filePath}`);
    }

    const localHash = hashFile(filePath);
    const responseHash = hashBuffer(responseBuffer);

    if (localHash !== responseHash) {
      throw new Error(
        `File is not the same`
      );
    }
  }
);
