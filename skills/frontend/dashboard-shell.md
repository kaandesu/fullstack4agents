# Skill: Dashboard Shell Layout

A persistent sidebar layout for authenticated pages. Uses the shadcn-vue `Sidebar` component.

Copy the layout to `frontend/app/layouts/dashboard.vue`. All dashboard pages then use it via `definePageMeta({ layout: 'dashboard' })`.

---

## Layout file

```vue
<!-- app/layouts/dashboard.vue -->
<template>
  <SidebarProvider>
    <AppSidebar />
    <SidebarInset>
      <header class="flex h-14 shrink-0 items-center gap-2 border-b px-4">
        <SidebarTrigger class="-ml-1" />
        <Separator orientation="vertical" class="h-4" />
        <!-- Page-specific header content via named slot -->
        <slot name="header" />
      </header>
      <main class="flex flex-1 flex-col gap-4 p-4">
        <slot />
      </main>
    </SidebarInset>
  </SidebarProvider>
</template>
```

---

## Sidebar component

Create `app/components/AppSidebar.vue` with your nav items as an array — never hardcode repeated links in the template:

```vue
<!-- app/components/AppSidebar.vue -->
<template>
  <Sidebar>
    <SidebarHeader class="p-4">
      <span class="font-semibold">My App</span>
    </SidebarHeader>

    <SidebarContent>
      <SidebarGroup>
        <SidebarGroupLabel>Navigation</SidebarGroupLabel>
        <SidebarGroupContent>
          <SidebarMenu>
            <SidebarMenuItem v-for="item in navItems" :key="item.href">
              <SidebarMenuButton as-child :is-active="route.path === item.href">
                <NuxtLink :to="item.href">
                  <Icon :name="item.icon" class="size-4" />
                  <span>{{ item.label }}</span>
                </NuxtLink>
              </SidebarMenuButton>
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarGroupContent>
      </SidebarGroup>
    </SidebarContent>

    <SidebarFooter class="p-4">
      <Button variant="ghost" class="w-full justify-start gap-2" @click="auth.logout()">
        <Icon name="lucide:log-out" class="size-4" />
        Sign out
      </Button>
    </SidebarFooter>
  </Sidebar>
</template>

<script setup lang="ts">
const route = useRoute()
const auth = useAuthManager()

// Define nav items as data — never repeat <SidebarMenuItem> blocks manually
const navItems = [
  { href: '/dashboard', label: 'Dashboard', icon: 'lucide:layout-dashboard' },
  { href: '/dashboard/posts', label: 'Posts', icon: 'lucide:file-text' },
  { href: '/dashboard/settings', label: 'Settings', icon: 'lucide:settings' },
]
</script>
```

---

## Using the layout in a page

```vue
<!-- app/pages/dashboard/index.vue -->
<template>
  <template #header>
    <h1 class="text-sm font-medium">Dashboard</h1>
  </template>

  <div>Page content here</div>
</template>

<script setup lang="ts">
definePageMeta({ layout: 'dashboard' })
</script>
```
