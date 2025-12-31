<script lang="ts">
  import * as AlertDialog from '$lib/components/ui/alert-dialog';
  import * as Select from "$lib/components/ui/select";
  import { Input } from '$lib/components/ui/input';
  import { Label } from '$lib/components/ui/label';
  import { uiState } from '$lib/stores/ui.svelte';
  import { projectState } from '$lib/stores/project.svelte';
  import {toast} from "svelte-sonner";

  let projectName = $state('');
  let minecraftVersion = $state('');
  let modLoader = $state('forge'); // e.g., 'Fabric', 'Forge', etc

  async function handleSubmit(event: Event) {
    event.preventDefault();
    try {
      await projectState.createProject(projectName, minecraftVersion, modLoader);
      toast.success('Project created successfully');

      projectName = '';
      minecraftVersion = '';
        modLoader = '';
    } catch (error) {
      toast.error('Failed to create project', { description: String(error) });
    }
  }
</script>

<AlertDialog.Root bind:open={uiState.newProjectDialogOpen}>
  <AlertDialog.Content>
    <AlertDialog.Header>
      <AlertDialog.Title>Create project</AlertDialog.Title>
      <AlertDialog.Description>
        Project config not found. Please initialize project.
      </AlertDialog.Description>
    </AlertDialog.Header>

    <form onsubmit={handleSubmit}>
      <div class="grid gap-4">
        <div class="grid gap-3">
          <Label for="name">Name</Label>
          <Input id="name" bind:value={projectName} autocomplete="off" required />
        </div>

        <div class="grid gap-3">
          <Label for="version">Minecraft Version</Label>
          <Input id="version" bind:value={minecraftVersion} autocomplete="off" required />
        </div>

          <div class="grid gap-3">
              <Label for="loader">Mod Loader</Label>
              <Select.Root
                      type="single"
                      bind:value={modLoader}
              >
                  <Select.Trigger class="w-full">
                      <span>{modLoader.charAt(0).toUpperCase() + modLoader.slice(1)}</span>
                  </Select.Trigger>
                  <Select.Content>
                      <Select.Item value={"forge"}>
                          Forge
                      </Select.Item>
                      <Select.Item value={"fabric"}>
                          Fabric
                      </Select.Item>
                  </Select.Content>
              </Select.Root>
          </div>
      </div>

      <AlertDialog.Footer class="mt-3">
        <AlertDialog.Cancel>Cancel</AlertDialog.Cancel>
        <AlertDialog.Action type="submit">Create</AlertDialog.Action>
      </AlertDialog.Footer>
    </form>
  </AlertDialog.Content>
</AlertDialog.Root>
