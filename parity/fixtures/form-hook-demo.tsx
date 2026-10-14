"use client"

import { zodResolver } from "@/node_modules/@hookform/resolvers/zod"
import { useForm } from "@/node_modules/react-hook-form"
import { z } from "@/node_modules/zod"

import { Button } from "@/registry/new-york-v4/ui/button"
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/registry/new-york-v4/ui/form"
import { Input } from "@/registry/new-york-v4/ui/input"

// Port fixture: the Form component (form.tsx) with react-hook-form and a
// zod schema: label, control, description and message state after submit
// and revalidation on change.
const FormSchema = z.object({
  username: z
    .string()
    .min(2, { message: "Username must be at least 2 characters." })
    .max(10, { message: "Username must be at most 10 characters." }),
  bio: z.string().max(20, { message: "Bio must be at most 20 characters." }),
})

export default function FormHookDemo() {
  const form = useForm<z.infer<typeof FormSchema>>({
    resolver: zodResolver(FormSchema),
    defaultValues: { username: "", bio: "" },
  })

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(() => {})} className="w-2/3 space-y-6">
        <FormField
          control={form.control}
          name="username"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Username</FormLabel>
              <FormControl>
                <Input placeholder="shadcn" {...field} />
              </FormControl>
              <FormDescription>This is your public display name.</FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="bio"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Bio</FormLabel>
              <FormControl>
                <Input placeholder="Tell us about you" {...field} />
              </FormControl>
              <FormMessage>Optional, up to 20 characters.</FormMessage>
            </FormItem>
          )}
        />
        <Button type="submit">Submit</Button>
      </form>
    </Form>
  )
}
