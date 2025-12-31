<script lang="ts">
    import * as Dialog from '$lib/components/ui/dialog';
  import * as Select from '$lib/components/ui/select';
  import { Label } from '$lib/components/ui/label';
  import { Button } from '$lib/components/ui/button';
  import { uiState } from '$lib/stores/ui.svelte';
  import { projectState } from '$lib/stores/project.svelte';
  import * as modService from '$lib/services/mod-service';
  import {toast} from "svelte-sonner";

  let searchResult = $state<string[] | null>(null);
  let loading = $state(false);
  let selectedVersion = $state('');

  $effect(() => {
    if (uiState.changeVersionModDialogOpen && uiState.selectedModId) {
      searchMod();
      console.log(uiState.changeVersionModDialogOpen, uiState.selectedModId);
    } else {
      searchResult = null;
      selectedVersion = '';
    }
  });

  async function searchMod() {
    if (!uiState.selectedModId) return;
    const mod = projectState.current.mods?.[uiState.selectedModId];
    if (!mod) return;
    loading = true;
    try {
      const results = await modService.getModVersions(uiState.selectedModId);
      if (results.length > 0) {
        searchResult = results;
        selectedVersion = mod.version;
      } else {
        searchResult = null;
      }
    } catch (error) {
      console.error('Failed to search for mod:', error);
      toast.error('Failed to search for mod', { description: String(error) });
    } finally {
      loading = false;
    }
  }

  async function handleChangeVersion() {
    if (!uiState.selectedModId || !selectedVersion) return;
    try {
      await modService.changeModVersion(uiState.selectedModId, selectedVersion);
      await projectState.refreshProject();
      toast.success('Mod version changed successfully');
    } catch (error) {
      console.error('Failed to change mod version:', error);
      toast.error('Failed to change mod version', { description: String(error) });
    }
    await uiState.closeChangeVersionModDialog();
  }

  function getCurrentVersion(): string {
    if (!uiState.selectedModId) return '';
    return projectState.current.mods?.[uiState.selectedModId]?.version || '';
  }
</script>

<Dialog.Root bind:open={uiState.changeVersionModDialogOpen}>
  <Dialog.Content class="max-w-lg max-h-9/12 overflow-y-auto">
    <Dialog.Header>
      <Dialog.Title>Change Mod Version</Dialog.Title>
      <Dialog.Description>
        Select a different version for the mod.
      </Dialog.Description>
    </Dialog.Header>

    <div class="grid gap-4">
      {#if loading}
        <p>Searching for mod...</p>
      {:else if searchResult}
        <div class="grid gap-3">
          <Label for="mod-version">Select Mod Version</Label>
          <Select.Root
            type="single"
            bind:value={selectedVersion}
          >
            <Select.Trigger class="w-full">
              <span>{selectedVersion || 'Select a version'}</span>
            </Select.Trigger>
            <Select.Content>
              {#each searchResult as version}
                <Select.Item value={version}>
                  {version}
                </Select.Item>
              {/each}
            </Select.Content>
          </Select.Root>
        </div>
      {:else}
        <p class="text-muted-foreground">Mod not found</p>
      {/if}
    </div>

    <Dialog.Footer>
        <Button
                disabled={!searchResult || selectedVersion === getCurrentVersion() || loading}
                onclick={handleChangeVersion}
                class="cursor-pointer"
        >
            Change to Selected Version
        </Button>
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>
