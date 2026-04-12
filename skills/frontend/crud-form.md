# Skill: CRUD Form

A create/edit form pattern using VeeValidate + Zod + shadcn-vue Form components.

When a form grows beyond 2–3 fields or is reused in multiple places, extract it into `app/components/<Feature>Form.vue`.

---

## Component

```vue
<!-- app/components/PostForm.vue — adapt to your model -->
<template>
  <form class="space-y-4" @submit="onSubmit">
    <FormField v-slot="{ componentField }" name="title">
      <FormItem>
        <FormLabel>Title</FormLabel>
        <FormControl>
          <Input placeholder="Post title" v-bind="componentField" />
        </FormControl>
        <FormMessage />
      </FormItem>
    </FormField>

    <FormField v-slot="{ componentField }" name="content">
      <FormItem>
        <FormLabel>Content</FormLabel>
        <FormControl>
          <Textarea placeholder="Write something…" v-bind="componentField" />
        </FormControl>
        <FormMessage />
      </FormItem>
    </FormField>

    <!-- Add more FormField blocks for additional fields -->

    <div class="flex justify-end gap-2">
      <Button type="button" variant="outline" @click="$emit('cancel')">
        Cancel
      </Button>
      <Button type="submit" :disabled="isSubmitting">
        {{ props.id ? 'Save changes' : 'Create' }}
      </Button>
    </div>
  </form>
</template>

<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { z } from 'zod'

const props = defineProps<{ id?: string }>()
const emit = defineEmits<{
  cancel: []
  saved: [record: unknown]
}>()

const manager = usePostsManager() // replace with your manager store

const schema = toTypedSchema(z.object({
  title: z.string().min(1, 'Required'),
  content: z.string().min(1, 'Required'),
}))

const { handleSubmit, isSubmitting, setValues } = useForm({ validationSchema: schema })

// Pre-fill when editing
if (props.id) {
  const existing = manager.posts.find((p) => p.id === props.id)
  if (existing) setValues({ title: existing.title, content: existing.content })
}

const onSubmit = handleSubmit(async (values) => {
  if (props.id) {
    await manager.updatePost(props.id, values, { onSuccess: () => emit('saved', values) })
  } else {
    await manager.createPost(values, { onSuccess: () => emit('saved', values) })
  }
})
</script>
```

---

## Using the form in a page

```vue
<template>
  <Dialog v-model:open="showForm">
    <DialogContent>
      <DialogHeader>
        <DialogTitle>{{ editId ? 'Edit Post' : 'New Post' }}</DialogTitle>
      </DialogHeader>
      <PostForm
        :id="editId"
        @saved="showForm = false"
        @cancel="showForm = false"
      />
    </DialogContent>
  </Dialog>
</template>

<script setup lang="ts">
const showForm = ref(false)
const editId = ref<string | undefined>()

const openEdit = (id: string) => {
  editId.value = id
  showForm.value = true
}
</script>
```
