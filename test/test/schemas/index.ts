import type z from "zod";
import { AuthSchema } from "./auth.schema";
import { UserSchema } from "./user.schema";
import { ActivitySchema, ActivitySchemas } from "./activity.schema";

// biome-ignore lint/suspicious/noExplicitAny: we can't determine what ObjectSchema output will be
type ObjectSchema = z.ZodType<Record<string, unknown>, any, any>;
type SchemaRegistry = Record<string, ObjectSchema>;

export const schemaRegistry: SchemaRegistry = {
  user: UserSchema,
  auth: AuthSchema,
  activity: ActivitySchema,
  activities: ActivitySchemas
};

export function isSchemaAvailable(
  key: string,
): key is keyof typeof schemaRegistry {
  return key in schemaRegistry;
}
