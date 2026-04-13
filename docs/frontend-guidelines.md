# Frontend Guidelines

Stack: Nuxt 4 · shadcn-vue · Tailwind CSS 4 · Pinia · VeeValidate + Zod · @vueuse/motion · PocketBase JS SDK

---

## Auto-imports

Never manually import these — Nuxt auto-imports them:

- Vue: `ref`, `computed`, `watch`, `watchEffect`, `onMounted`, `storeToRefs`, etc.
- Nuxt: `navigateTo`, `useRoute`, `useRouter`, `useFetch`, `definePageMeta`, `useColorMode`, `useRuntimeConfig`
- VueUse: all composables from `@vueuse/core`
- Components: everything under `app/components/` (including `ui/`)
- Utils: all files under `app/utils/`
- Composables: all files under `app/composables/`

**Exception: Always explicitly import Pinia stores** even though they're auto-imported. Edge cases exist where auto-import doesn't resolve correctly.

```ts
// Wrong
import { ref } from "vue";
import { Button } from "@/components/ui/button";

// Correct — explicit import for stores
import { useAuthManager } from "~/stores/auth-manager";

const count = ref(0);
const auth = useAuthManager();
// <Button /> works directly in template
```

---

## Component naming

Component name = directory path + filename, duplicates removed.

```
app/components/base/foo/Button.vue  →  <BaseFooButton />
app/components/AppSidebar.vue       →  <AppSidebar />
app/components/ui/Card.vue          →  <Card /> (shadcn prefix is "")
```

---

## When to extract a component

**Extract into a component when any of these are true:**

- The template block is used in more than one place
- A section has its own local state or logic
- The template grows past ~80 lines
- A block represents a distinct UI concept (a card, a list item, a form)

Keep pages thin. Pages should be mostly composition — a page assembles components, not raw HTML.

```vue
<!-- Wrong — page doing too much inline work -->
<template>
  <div>
    <div v-for="post in posts" class="border rounded p-4 ...">
      <h2>{{ post.title }}</h2>
      <p>{{ post.content }}</p>
      <button @click="deletePost(post.id)">Delete</button>
    </div>
  </div>
</template>

<!-- Correct — extract the list item -->
<template>
  <div>
    <PostCard
      v-for="post in posts"
      :key="post.id"
      :post="post"
      @delete="deletePost"
    />
  </div>
</template>
```

---

## Use arrays + v-for instead of repeating template blocks

Anytime you write the same element more than twice, define the data as an array and render with `v-for`.

```vue
<!-- Wrong — three nearly identical blocks -->
<template>
  <div>
    <NuxtLink to="/dashboard">Dashboard</NuxtLink>
    <NuxtLink to="/posts">Posts</NuxtLink>
    <NuxtLink to="/settings">Settings</NuxtLink>
  </div>
</template>

<!-- Correct -->
<template>
  <div>
    <NuxtLink v-for="item in navItems" :key="item.href" :to="item.href">
      {{ item.label }}
    </NuxtLink>
  </div>
</template>

<script setup lang="ts">
const navItems = [
  { href: "/dashboard", label: "Dashboard" },
  { href: "/posts", label: "Posts" },
  { href: "/settings", label: "Settings" },
];
</script>
```

This applies to: nav links, stat cards, tab lists, badge sets, feature lists, form option groups — anything structural that repeats.

---

## Use the component library — don't reimplement

Before writing a custom component, check if shadcn-vue already has it.

| Need          | Use                                            |
| ------------- | ---------------------------------------------- |
| Button        | `<Button>`                                     |
| Modal/overlay | `<Dialog>`                                     |
| Dropdown      | `<DropdownMenu>`                               |
| Side panel    | `<Sheet>`                                      |
| Notification  | `createToast({ ... })()` via `vue-sonner`      |
| Table         | `<Table>` + TanStack (see data-table skill)    |
| Form field    | `<FormField>` + `<FormItem>` + `<FormMessage>` |
| Loading       | `<Skeleton>`                                   |
| Tooltip       | `<Tooltip>`                                    |
| Tabs          | `<Tabs>`                                       |
| Accordion     | `<Accordion>`                                  |

---

## Layouts

Create layouts in `app/layouts/`. Apply per-page with `definePageMeta`.

```vue
<!-- app/layouts/dashboard.vue -->
<template>
  <div class="flex h-screen">
    <AppSidebar />
    <main class="flex-1 overflow-auto p-6"><slot /></main>
  </div>
</template>
```

```ts
// app/pages/dashboard/index.vue
definePageMeta({ layout: "dashboard" });
```

**Named templates and slots:**

- Named templates (`<template #header>`, `<template #footer>`) should NOT be root nodes
- If using named slot templates, wrap them in `<NuxtLayout>` component
- For pages rendering inside a layout, use `definePageMeta({ layout: '---' })` instead of `<template #name>` slots (MUCH preferred approach)

```vue
<!-- Good — using definePageMeta -->
<script setup lang="ts">
definePageMeta({ layout: "dashboard" });
</script>
<template>
  <div>Page content</div>
</template>

<!-- Also okay — explicit NuxtLayout wrapper for named slots -->
<template>
  <NuxtLayout name="dashboard">
    <template #header>
      <Header />
    </template>
  </NuxtLayout>
</template>
```

- Default layout (`default.vue`) applies when no layout is specified
- `definePageMeta({ layout: false })` for standalone pages (auth, landing)
- See `skills/frontend/dashboard-shell.md` for a full sidebar layout

---

## Pages and routing

```
app/pages/index.vue           →  /
app/pages/dashboard/index.vue →  /dashboard
app/pages/users/[id].vue      →  /users/:id
```

---

## State management — store managers

Every feature gets a `*-manager.ts` Pinia store in `app/stores/`. The template includes `auth-manager.ts` as the base pattern.

```ts
// app/stores/posts-manager.ts
export const usePostsManager = defineStore("PostsManager", () => {
  const auth = useAuthManager();
  const config = useRuntimeConfig();
  const posts = ref<Post[]>([]);

  const fetchPosts = async () => {
    const res = await auth.API(() =>
      $fetch<{ data: Post[] }>("/api/posts", {
        baseURL: config.public.apiUrl as string,
        headers: { Authorization: `Bearer ${auth.token}` },
      }),
    );
    posts.value = res.data ?? [];
  };

  return { posts, fetchPosts };
});
```

See `skills/frontend/auth-manager.md` for the full pattern.

---

## API wrapper — `auth.API()`

**Every Gin API call must be wrapped in `auth.API(fn)`.**

This handles 401 responses by refreshing the PocketBase session and retrying automatically. Use it at the call site, not inside helper definitions:

```ts
// Correct — wrap the call
const data = await auth.API(() => $fetch('/api/posts', {
  baseURL: config.public.apiUrl,
  headers: { Authorization: `Bearer ${auth.token}` },
}))

// Wrong — wrapping the definition (doesn't help, retry needs to re-execute the fetch)
const fetchPosts = auth.API(async () => { ... }) // ← no
```

**Don't use `API()` for PocketBase SDK calls** — the SDK handles retries itself.

---

## Forms

Always use VeeValidate + Zod. Never write manual validation or conditional error display.

```vue
<script setup lang="ts">
import { toTypedSchema } from "@vee-validate/zod";
import { useForm } from "vee-validate";
import { z } from "zod";

const schema = toTypedSchema(
  z.object({
    email: z.string().email(),
    name: z.string().min(2),
  }),
);
const { handleSubmit, isSubmitting } = useForm({ validationSchema: schema });
const onSubmit = handleSubmit(async (values) => {
  /* ... */
});
</script>
```

Use `<FormField>`, `<FormItem>`, `<FormLabel>`, `<FormControl>`, `<FormMessage>` for layout.  
When a form has more than 3 fields or is reused, extract it into a `*Form.vue` component.

---

## Auth middleware

The template includes `app/middleware/auth.ts` which:

- Redirects unauthenticated users to `/auth/sign-in`
- Redirects authenticated users away from `/auth/*`
- Tries to refresh the PocketBase session on first protected-page load

To mark a page as protected (default behavior), do nothing — the middleware applies globally.  
To opt a route out, add it to `publicRoutes` in `middleware/auth.ts`.

---

## PocketBase SDK

Use the `usePocketBase()` composable (singleton):

```ts
const pb = usePocketBase();

// Direct collection access
const posts = await pb.collection("posts").getFullList({ sort: "-created" });

// Auth
await pb.collection("users").authWithPassword(email, password);
```

For Gin API calls use `$fetch` + `auth.API()` as shown above.

---

## Animations

Use `v-motion` for entrance animations. Keep them subtle — one direction, short duration.

```vue
<div v-motion :initial="{ opacity: 0, y: 20 }" :enter="{ opacity: 1, y: 0 }">
  Content
</div>
```

---

## Dark mode

Use `useColorMode()` — never hardcode `dark:` classes manually in new components.

```ts
const colorMode = useColorMode();
colorMode.preference = "dark";
```

---

## Toasts

Toasts are managed by `useToastManagerStore` (`app/stores/toast-manager.ts`). Always use the `createToast` util (`app/utils/create-toast.ts`) which delegates to the store — never call `toast` from vue-sonner directly.

`createToast` returns a function. Call it immediately or pass it as a callback:

```ts
// Immediate
createToast({ message: "Saved!", type: "success" })();

// As callback — no extra wrapper needed
auth.signIn({
  email,
  password,
  onSuccess: createToast({ message: "Welcome back", type: "success" }),
});
```

Full options (`ToastCreateParam`):

```ts
createToast({
  message: "Post deleted",
  type: "error", // 'success' | 'error' | 'info' | 'warning' | 'loading' | 'default'
  toastOps: {
    description: "Could not reach the server",
    duration: 5000,
    position: "bottom-right",
    action: {
      label: "Retry",
      onClick: () => retryFn(),
    },
  },
})();
```
