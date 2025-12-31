<script lang="ts">
  import * as Dialog from '$lib/components/ui/dialog';
  import { Button } from '$lib/components/ui/button';
  import { uiState } from '$lib/stores/ui.svelte';
  import { projectState } from '$lib/stores/project.svelte';
  import * as modService from '$lib/services/mod-service';
  import type { ModUpdateInfo } from '$lib/types/mod';
  import {Checkbox} from "$lib/components/ui/checkbox";
  import {toast} from "svelte-sonner";

  let updates = $state<ModUpdateInfo[]>([]);
  let loading = $state(false);
  let selectedUpdates = $state<Set<number>>(new Set());

  $effect(() => {
    if (uiState.updateModsDialogOpen) {
      checkForUpdates();
    } else {
      updates = [];
      selectedUpdates = new Set();
    }
  });

  async function checkForUpdates() {
    if (!projectState.current.mods) return;
    loading = true;
    try {
      const modIds = Object.keys(projectState.current.mods);
      updates = await modService.checkModsUpdates(modIds);
      updates.map((update, i) => {
          selectedUpdates.add(i)
          selectedUpdates = new Set(selectedUpdates);
      })
    } catch (error) {
      console.error('Failed to check for updates:', error);
      toast.error('Failed to check for updates', { description: String(error) });
    } finally {
      loading = false;
    }
  }

  async function handleUpdate() {
      const selectedMods = Array.from(selectedUpdates).map(i => updates[i]);
      try {
          await modService.updateMods(selectedMods);
          await projectState.refreshProject()
          toast.success('Mods updated successfully');
      } catch (error) {
          console.error('Failed to update mods:', error);
          toast.error('Failed to update mods', { description: String(error) });
      }
      await uiState.closeUpdateModsDialog();
  }

  function getModName(modId: string): string {
    const mod = projectState.current.mods?.[modId];
    return mod?.filename || modId;
  }

  function getCurrentVersion(modId: string): string {
    return projectState.current.mods?.[modId]?.version || '';
  }
</script>

<Dialog.Root bind:open={uiState.updateModsDialogOpen}>
  <Dialog.Content class="max-w-lg max-h-9/12 overflow-y-auto">
    <Dialog.Header>
      <Dialog.Title>Update Mods</Dialog.Title>
      <Dialog.Description>
        Check for updates to installed mods and select which ones to update.
      </Dialog.Description>
    </Dialog.Header>

    <div class="grid gap-4">
      {#if loading}
        <p>Checking for updates...</p>
      {:else if updates.length === 0}
        <p class="text-muted-foreground">Everything is up to date</p>
      {:else}
        <ul class="space-y-2">
          {#each updates as update, i}
            <li class="flex items-center space-x-2 p-2 border rounded">
              <Checkbox
                id={update.ModId}
                checked={selectedUpdates.has(i)}
                onCheckedChange={(checked) => {
                  if (checked) {
                    selectedUpdates.add(i);
                  } else {
                    selectedUpdates.delete(i);
                  }
                  selectedUpdates = new Set(selectedUpdates);
                }}
              />
              <label for={update.ModId} class="flex-1 cursor-pointer">
                <span class="font-semibold">{getModName(update.ModId)}</span>
                <span class="text-sm text-muted-foreground">
                  {getCurrentVersion(update.ModId)} → {update.Version}
                </span>
              </label>
            </li>
          {/each}
        </ul>
      {/if}
    </div>

    <Dialog.Footer>
        <Button
          disabled={selectedUpdates.size === 0}
          onclick={handleUpdate}
          class="cursor-pointer"
        >
          Update Selected ({selectedUpdates.size})
        </Button>
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>
