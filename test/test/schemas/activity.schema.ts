import z from "zod";

export const ActivityTypes = {
  Yoga: 4,
  Stretching: 4,
  Cycling: 8,
  Swimming: 8,
  Dancing: 8,
  Hiking: 10,
  Running: 10,
  HIIT: 10,
  JumpRope: 10,
} as const;

export const ActivityTypes_E = Object.keys(ActivityTypes) as Array<keyof typeof ActivityTypes>;
export type ActivityType_E = (typeof ActivityTypes_E)[number];

export const ActivityRequestSchema = z.object({
  activityType: z.enum(ActivityTypes_E),
  doneAt: z.iso.datetime(),
  durationInMinutes: z.number().min(1)
});

export const ActivitySchema = z.object({
  activityId: z.string(),
  activityType: z.enum([
    "Yoga",
    "Stretching",
    "Cycling",
    "Swimming",
    "Dancing",
    "Hiking",
    "Running",
    "HIIT",
    "JumpRope"]),
  doneAt: z.iso.datetime(),
  durationInMinutes: z.number(),
  caloriesBurned: z.number(),
  createdAt: z.iso.datetime(),
  updatedAt: z.iso.datetime()
});

export const ActivitySchemas = z.array(ActivitySchema);
