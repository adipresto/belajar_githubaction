import z from "zod";

export const AuthSchema = z.object({
  email: z.email(),
  token: z.string().min(1),
});
