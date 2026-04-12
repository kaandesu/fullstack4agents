# Skill: Auth Pages

Login and register pages using shadcn-vue Form, VeeValidate, Zod, and `useAuthManager`.

Copy to `frontend/app/pages/auth/`. Both pages use `definePageMeta({ layout: false })` — they are standalone, outside the main app shell.

---

## Login page

```vue
<!-- app/pages/auth/sign-in.vue -->
<template>
  <div class="flex min-h-screen items-center justify-center p-4">
    <Card class="w-full max-w-sm">
      <CardHeader>
        <CardTitle>Sign in</CardTitle>
        <CardDescription>Enter your credentials to continue.</CardDescription>
      </CardHeader>
      <CardContent>
        <form class="space-y-4" @submit="onSubmit">
          <FormField v-slot="{ componentField }" name="email">
            <FormItem>
              <FormLabel>Email</FormLabel>
              <FormControl>
                <Input type="email" placeholder="you@example.com" v-bind="componentField" />
              </FormControl>
              <FormMessage />
            </FormItem>
          </FormField>

          <FormField v-slot="{ componentField }" name="password">
            <FormItem>
              <FormLabel>Password</FormLabel>
              <FormControl>
                <Input type="password" v-bind="componentField" />
              </FormControl>
              <FormMessage />
            </FormItem>
          </FormField>

          <Button type="submit" class="w-full" :disabled="isSubmitting || auth.loading">
            <span v-if="isSubmitting || auth.loading">Signing in…</span>
            <span v-else>Sign in</span>
          </Button>
        </form>
      </CardContent>
      <CardFooter class="justify-center text-sm text-muted-foreground">
        No account?
        <NuxtLink to="/auth/sign-up" class="ml-1 underline">Sign up</NuxtLink>
      </CardFooter>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { z } from 'zod'

definePageMeta({ layout: false })

const auth = useAuthManager()

const schema = toTypedSchema(z.object({
  email: z.string().email(),
  password: z.string().min(8),
}))

const { handleSubmit, isSubmitting } = useForm({ validationSchema: schema })

const onSubmit = handleSubmit((values) => {
  auth.signIn({
    email: values.email,
    password: values.password,
    onSuccess: () => navigateTo('/dashboard'),
  })
})
</script>
```

---

## Register page

```vue
<!-- app/pages/auth/sign-up.vue -->
<template>
  <div class="flex min-h-screen items-center justify-center p-4">
    <Card class="w-full max-w-sm">
      <CardHeader>
        <CardTitle>Create account</CardTitle>
        <CardDescription>Fill in the details below to get started.</CardDescription>
      </CardHeader>
      <CardContent>
        <form class="space-y-4" @submit="onSubmit">
          <FormField v-slot="{ componentField }" name="name">
            <FormItem>
              <FormLabel>Name</FormLabel>
              <FormControl>
                <Input placeholder="Your name" v-bind="componentField" />
              </FormControl>
              <FormMessage />
            </FormItem>
          </FormField>

          <FormField v-slot="{ componentField }" name="email">
            <FormItem>
              <FormLabel>Email</FormLabel>
              <FormControl>
                <Input type="email" placeholder="you@example.com" v-bind="componentField" />
              </FormControl>
              <FormMessage />
            </FormItem>
          </FormField>

          <FormField v-slot="{ componentField }" name="password">
            <FormItem>
              <FormLabel>Password</FormLabel>
              <FormControl>
                <Input type="password" v-bind="componentField" />
              </FormControl>
              <FormMessage />
            </FormItem>
          </FormField>

          <Button type="submit" class="w-full" :disabled="isSubmitting || auth.loading">
            <span v-if="isSubmitting || auth.loading">Creating account…</span>
            <span v-else>Create account</span>
          </Button>
        </form>
      </CardContent>
      <CardFooter class="justify-center text-sm text-muted-foreground">
        Already have an account?
        <NuxtLink to="/auth/sign-in" class="ml-1 underline">Sign in</NuxtLink>
      </CardFooter>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { z } from 'zod'

definePageMeta({ layout: false })

const auth = useAuthManager()

const schema = toTypedSchema(z.object({
  name: z.string().min(2),
  email: z.string().email(),
  password: z.string().min(8),
}))

const { handleSubmit, isSubmitting } = useForm({ validationSchema: schema })

const onSubmit = handleSubmit((values) => {
  auth.signUp({
    name: values.name,
    email: values.email,
    password: values.password,
    onSuccess: () => navigateTo('/dashboard'),
  })
})
</script>
```
