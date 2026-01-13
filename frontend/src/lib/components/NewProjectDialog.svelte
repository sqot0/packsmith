<script lang="ts">
  import * as AlertDialog from '$lib/components/ui/alert-dialog';
  import * as Select from "$lib/components/ui/select";
  import { Input } from '$lib/components/ui/input';
  import { Label } from '$lib/components/ui/label';
  import { uiState } from '$lib/stores/ui.svelte';
  import { projectState } from '$lib/stores/project.svelte';
  import {toast} from "svelte-sonner";
  import ForgeLogo from "$lib/assets/forge.png"
  import FabricLogo from "$lib/assets/fabric.png"
  import NeoForgeLogo from "$lib/assets/neoforge.png"
  import QuiltLogo from "$lib/assets/quilt.png"

  type ModLoader = 'forge' | 'fabric' | 'neoforge' | 'quilt';

  let projectName = $state('');
  let minecraftVersion = $state('');
  let modLoader = $state<ModLoader>('forge');

  let modLoaders = {
      "forge": {
          "name": "Forge",
          "logo": ForgeLogo
      },
        "fabric": {
            "name": "Fabric",
            "logo": FabricLogo
        },
        "neoforge": {
            "name": "NeoForge",
            "logo": NeoForgeLogo
        },
        "quilt": {
            "name": "Quilt",
            "logo": QuiltLogo
        }

  }

  async function handleSubmit(event: Event) {
    event.preventDefault();
    try {
      await projectState.createProject(projectName, minecraftVersion, modLoader);
      toast.success('Project created successfully');

      projectName = '';
      minecraftVersion = '';
        modLoader = 'forge';
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

                      <span>
                          <img src={modLoaders[modLoader].logo} alt="{modLoaders[modLoader].name} Logo" class="inline h-5 mr-1">
                          {modLoaders[modLoader].name}
                      </span>
                  </Select.Trigger>
                  <Select.Content>
                      {#each Object.entries(modLoaders) as [key, loader]}
                          <Select.Item value={key}>
                              <img src={loader.logo} alt="{loader.name} Logo" class="inline h-5">
                              {loader.name}
                          </Select.Item>
                      {/each}
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
