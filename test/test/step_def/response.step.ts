import assert from "node:assert";
import { Given, Then } from "@cucumber/cucumber";
import get from "lodash.get";
import z from "zod";
import { isSchemaAvailable, schemaRegistry } from "../schemas";
import type { ApiWorld } from "../types/cucumber";

Given("{word} exists", function(this: ApiWorld, name: string) {
  const inSavedPayload = this.savedPayload[name];

  if (inSavedPayload === undefined) {
    throw new Error(`Expected '${name}' to exist in tokens or savedPayload`);
  }
});

Then(
  "the response body matches the {word} schema",
  async function(this: ApiWorld, schemaName: string) {
    if (!this.lastResponse) {
      throw new Error("No response found in World");
    }

    if (!isSchemaAvailable(schemaName)) {
      throw new Error(`Schema ${schemaName} not found in World`);
    }

    const schema = schemaRegistry[schemaName];
    if (!schema) {
      throw new Error(`Schema '${schemaName}' not found in SchemaRegistry`);
    }

    const body = await this.lastResponse.clone().json();
    const result = schema.safeParse(body);

    if (!result.success) {
      throw new Error(
        `Response body does not match schema '${schemaName}':\n` +
        JSON.stringify(z.treeifyError(result.error), null, 2),
      );
    }

    this.jsonBody = result.data;
    this.lastSchemaName = schemaName;
  },
);

Then(
  "the response should be {int}",
  async function(this: ApiWorld, status: number) {
    if (!this.lastResponse) {
      throw new Error("No response found in World");
    }
    assert.strictEqual(this.lastResponse.status, status);
  },
);

Then(
  "the violation responses should be {int}",
  function(this: ApiWorld, expectedStatus: number) {
    if (!this.violationResponses || this.violationResponses.length === 0) {
      throw new Error("No violation responses found in World");
    }

    for (const { status } of this.violationResponses) {
      assert.strictEqual(
        status,
        expectedStatus,
        `Expected status ${expectedStatus}, but got ${status}`,
      );
    }
  },
);

Then(
  "{string} value is equal to {string}",
  function(this: ApiWorld, path: string, expected: string) {
    if (!this.jsonBody)
      throw new Error("jsonBody is empty. Did you validate schema first?");

    const actual = get(this.jsonBody, path);

    if (Array.isArray(actual)) {
      assert.deepStrictEqual(
        actual.map((d) => d.toString()),
        Array(actual.length).fill(expected),
        `Expected all values at path '${path}' to equal '${expected}', but got ${JSON.stringify(actual)}`,
      );
      return;
    }
    assert.strictEqual(
      actual,
      expected,
      `Expected value at path '${path}' to equal '${expected}', but got '${actual}'`,
    );
  },
);

Then(
  "{string} value is equal to {int}",
  function(this: ApiWorld, path: string, expected: number) {
    if (!this.jsonBody)
      throw new Error("jsonBody is empty. Did you validate schema first?");

    const actual = get(this.jsonBody, path);

    assert.strictEqual(
      actual,
      expected,
      `Expected value at path '${path}' to equal '${expected}', but got '${actual}'`,
    );
  },
);

Then(
  "{string} value is equal to saved {word}",
  function(this: ApiWorld, path: string, name: string) {
    if (!this.jsonBody)
      throw new Error("jsonBody is empty. Did you validate schema first?");

    const actual = get(this.jsonBody, path);

    const expected = this.savedPayload[name];

    if (expected === undefined) {
      throw new Error(`Expected '${name}' to exist in tokens or savedPayload`);
    }

    assert.strictEqual(
      actual,
      expected,
      `Expected value at path '${path}' to equal '${expected}', but got '${actual}'`,
    );
  },
);

Then(
  "{string} value is null",
  function(this: ApiWorld, path: string) {
    if (!this.jsonBody)
      throw new Error("jsonBody is empty. Did you validate schema first?");

    const actual = get(this.jsonBody, path);

    assert.strictEqual(
      actual,
      null,
      `Expected value at path '${path}' to null, but got '${actual}'`,
    );
  },
);

Then(
  "the fetched data total are {int} datas", async function(this: ApiWorld, expectedCount: number) {
    if (!this.lastSchemaName)
      throw new Error("No schema name found in World");

    if (!this.lastResponse)
      throw new Error("No response found in World");

    const schema = schemaRegistry[this.lastSchemaName];
    if (!schema) {
      throw new Error(`Schema '${this.lastSchemaName}' not found in SchemaRegistry`);
    }

    const body = await this.lastResponse.clone().json();
    const result = schema.safeParse(body);

    const actualCount = (result.data as any[]).length;

    assert.strictEqual(
      expectedCount,
      actualCount,
      `Expected data count of '${this.lastSchemaName}' is '${expectedCount}', but got '${actualCount}'`,
    );
  },
);
