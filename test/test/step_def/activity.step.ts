import { Given, Then } from "@cucumber/cucumber";
import { ApiWorld } from "../types/cucumber";
import { schemaRegistry } from "../schemas";
import get from "lodash.get";
import assert from "node:assert";
import { payloads } from "../payloads";
import { resolvePayload } from "../support/helper";
import { ActivityTypes_E } from "../schemas/activity.schema";
import { generateRandomNumber } from "../helper/generator";

Given("I prepare about 100 {word} to inject and send right away", async function(this: ApiWorld, payloadName: string) {
  const payload = payloads[payloadName];
  if (!payload) throw new Error(`Payload '${payloadName}' not found`);

  const resolved = resolvePayload(payload.positive, this);

  this.request.options.headers["Content-Type"] = "application/json";
  for (let typeIndex = 0; typeIndex < 4; typeIndex++) {
    for (let itemIndex = 0; itemIndex < 25; itemIndex++) {
      resolved["activityType"] = ActivityTypes_E[typeIndex];
      resolved["durationInMinutes"] = generateRandomNumber(5, 120)
      this.request.options.body = JSON.stringify(resolved);

      if (!this.request.url) {
        throw new Error(`Request url not found`);
      }
      this.lastResponse = await fetch(this.request.url, this.request.options);

      this.jsonBody = (await this.lastResponse
        ?.clone()
        .json()
        .catch(() => undefined)) as Record<string, unknown> | undefined;
    }
  }
});

Then(
  "doneAt must be the result of adding the durationInMinutes to createdAt",
  function(this: ApiWorld) {
    if (!this.jsonBody)
      throw new Error("jsonBody is empty. Did you validate schema first?");

    const doneAtRaw = get(this.jsonBody, "doneAt") as string;
    const createdAtRaw = get(this.jsonBody, "createdAt") as string;
    const duration = get(this.jsonBody, "durationInMinutes") as number;

    if (!doneAtRaw || !createdAtRaw || typeof duration !== "number") {
      throw new Error("doneAt, createdAt, or durationInMinutes missing in response JSON");
    }

    // Parse dates
    const actual = new Date(doneAtRaw);
    const created = new Date(createdAtRaw);

    // Normalize createdAt to avoid drift (sec & ms)
    created.setSeconds(0, 0);

    // Calculate expected
    const expectedDoneAt = new Date(created.getTime() + duration * 60_000);
    expectedDoneAt.setSeconds(0, 0);

    // Normalize actual doneAt
    actual.setSeconds(0, 0);

    const expectedISO = expectedDoneAt.toISOString();
    const actualISO = actual.toISOString();

    assert.strictEqual(
      actualISO,
      expectedISO,
      `Expected doneAt to equal ${expectedISO}, but got ${actualISO}`
    );
  }
);

Then(
  "the value of {word} are between {int} and {int}",
  async function(this: ApiWorld, paramName: string, from: number, to: number) {
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

    const actuals = result.data as Record<string, any>[];

    const allInRange = actuals.every(a => a[paramName] >= from && a[paramName] <= to);

    assert.strictEqual(
      allInRange,
      true,
      `Expected '${paramName}' are between '${from}' to '${to}'.`,
    );
  },
);

Then(
  "the value of {word} is {word}",
  async function(this: ApiWorld, paramName: number, actual: number) {
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

    const actuals = result.data as Record<string, any>[];

    const allInRange = actuals.every(a => a[paramName] == actual);

    assert.strictEqual(
      allInRange,
      true,
      `Expected '${paramName} values are '${actual}'.`,
    );
  },
);
