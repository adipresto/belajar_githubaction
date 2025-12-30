import { z } from "zod";

export const PreferenceTypes = [
  "CARDIO",
  "WEIGHT",
] as const;

export const WeightUnitTypes = [
  "KG",
  "LBS",
] as const;

export const HeightUnitTypes = [
  "CM",
  "INCH",
] as const;

export const UserSchema = z.object({
  preference: z.enum([
    "CARDIO",
    "WEIGHT"
  ]).nullable(),
  weightUnit: z.enum([
    "KG",
    "LBS"
  ]).nullable(),
  heightUnit: z.enum([
    "CM",
    "INCH"
  ]).nullable(),
  weight: z.number().nullable(),
  height: z.number().nullable(),
  email: z.email(),
  name: z.string().nullable(),
  imageUri: z.string().nullable(),
});
