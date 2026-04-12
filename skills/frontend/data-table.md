# Skill: Data Table

A reusable TanStack Table component with sorting and pagination. Uses shadcn-vue Table components.

Create `app/components/DataTable.vue` and use it from any page.

---

## Component

```vue
<!-- app/components/DataTable.vue -->
<template>
  <div class="space-y-4">
    <div class="rounded-md border">
      <Table>
        <TableHeader>
          <TableRow v-for="headerGroup in table.getHeaderGroups()" :key="headerGroup.id">
            <TableHead
              v-for="header in headerGroup.headers"
              :key="header.id"
              :class="header.column.getCanSort() ? 'cursor-pointer select-none' : ''"
              @click="header.column.getToggleSortingHandler()?.($event)"
            >
              <FlexRender :render="header.column.columnDef.header" :props="header.getContext()" />
              <Icon v-if="header.column.getIsSorted() === 'asc'" name="lucide:arrow-up" class="ml-1 inline size-3" />
              <Icon v-else-if="header.column.getIsSorted() === 'desc'" name="lucide:arrow-down" class="ml-1 inline size-3" />
            </TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <template v-if="table.getRowModel().rows.length">
            <TableRow v-for="row in table.getRowModel().rows" :key="row.id">
              <TableCell v-for="cell in row.getVisibleCells()" :key="cell.id">
                <FlexRender :render="cell.column.columnDef.cell" :props="cell.getContext()" />
              </TableCell>
            </TableRow>
          </template>
          <TableRow v-else>
            <TableCell :colspan="columns.length" class="h-24 text-center text-muted-foreground">
              No results.
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </div>

    <div class="flex items-center justify-between text-sm text-muted-foreground">
      <span>Page {{ table.getState().pagination.pageIndex + 1 }} of {{ table.getPageCount() }}</span>
      <div class="flex gap-2">
        <Button variant="outline" size="sm" :disabled="!table.getCanPreviousPage()" @click="table.previousPage()">
          Previous
        </Button>
        <Button variant="outline" size="sm" :disabled="!table.getCanNextPage()" @click="table.nextPage()">
          Next
        </Button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts" generic="TData">
import {
  FlexRender,
  useVueTable,
  getCoreRowModel,
  getPaginationRowModel,
  getSortedRowModel,
} from '@tanstack/vue-table'
import type { ColumnDef, SortingState } from '@tanstack/vue-table'

const props = defineProps<{
  data: TData[]
  columns: ColumnDef<TData>[]
  pageSize?: number
}>()

const sorting = ref<SortingState>([])

const table = useVueTable({
  get data() { return props.data },
  get columns() { return props.columns },
  initialState: { pagination: { pageSize: props.pageSize ?? 10 } },
  getCoreRowModel: getCoreRowModel(),
  getPaginationRowModel: getPaginationRowModel(),
  getSortedRowModel: getSortedRowModel(),
  state: { get sorting() { return sorting.value } },
  onSortingChange: (updater) => {
    sorting.value = typeof updater === 'function' ? updater(sorting.value) : updater
  },
})
</script>
```

---

## Usage in a page

Define columns as a typed array — never hardcode table rows in the template:

```vue
<template>
  <DataTable :data="manager.posts" :columns="columns" />
</template>

<script setup lang="ts">
import type { ColumnDef } from '@tanstack/vue-table'
import type { Post } from '~/stores/posts-manager'

definePageMeta({ layout: 'dashboard' })

const manager = usePostsManager()
onMounted(() => manager.fetchPosts())

const columns: ColumnDef<Post>[] = [
  { accessorKey: 'title', header: 'Title' },
  { accessorKey: 'created', header: 'Created' },
  {
    id: 'actions',
    cell: ({ row }) => h('button', { onClick: () => manager.deletePost(row.original.id) }, 'Delete'),
  },
]
</script>
```
