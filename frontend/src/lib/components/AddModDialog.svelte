<script lang="ts">
    import * as Dialog from '$lib/components/ui/dialog';
  import * as Select from '$lib/components/ui/select';
  import { Input } from '$lib/components/ui/input';
  import { Label } from '$lib/components/ui/label';
  import { Button } from '$lib/components/ui/button';
  import { uiState } from '$lib/stores/ui.svelte';
  import { modSearchState } from '$lib/stores/mods.svelte';
  import { projectState } from '$lib/stores/project.svelte';
  import * as modService from '$lib/services/mod-service';
  import type { ModSide, ModSearchResult } from '$lib/types/mod';
  import CurseForgeIcon from '$lib/icons/CurseforgeIcon.svelte';
  import ModrinthIcon from '$lib/icons/ModrinthIcon.svelte';
  import { Search, Server, Laptop, ExternalLink, Plus, Loader } from '@lucide/svelte';
  import { BrowserOpenURL } from '../../../wailsjs/runtime';
    import {toast} from "svelte-sonner";

  let selectedSide = $state<ModSide>('both');
  let loadingMods = $state<Record<string, true>>({});
  let selectedVersions = $state<Record<string, string>>({});

  async function handleSearch(event: Event) {
    event.preventDefault();
    await modSearchState.clearResults();
    await modSearchState.search();
  }

  function handlePlatformChange() {
    modSearchState.clearResults();
    selectedVersions = {};
  }

  async function handleAddMod(result: ModSearchResult) {
    const { ID, ClientSide, ServerSide, URL } = result;

    if (modService.requiresSideSelection(modSearchState.platform, ClientSide, ServerSide)) {
      await uiState.openAddModSideDialog(result);
    } else {
      const side = modService.determineModSide(modSearchState.platform, ClientSide, ServerSide);
      loadingMods = {...loadingMods, [ID]: true};
      try {
        await modService.addMod(ID, modSearchState.platform, { URL, Side: side, Version: selectedVersions[ID] || result.Versions[0] });
        await projectState.refreshProject();
        toast.success('Mod added successfully');
      } catch (error) {
        console.error('Failed to add mod:', error);
        toast.error('Failed to add mod', { description: String(error) });
      } finally {
        loadingMods = (({ [ID]: _, ...rest }) => rest)(loadingMods);
      }
    }
  }

  async function handleAddModWithSide() {
    const mod = uiState.modToAdd;
    if (!mod) return;

    loadingMods = {...loadingMods, [mod.ID]: true};
    try {
      await modService.addMod(mod.ID, modSearchState.platform, {
        URL: mod.URL,
        Side: selectedSide,
        Version: selectedVersions[mod.ID] || mod.Versions[0]
      });
      await projectState.refreshProject()
        toast.success('Mod added successfully');
    } catch (error) {
      console.error('Failed to add mod:', error);
        toast.error('Failed to add mod', { description: String(error) });
    } finally {
      loadingMods = (({ [mod.ID]: _, ...rest }) => rest)(loadingMods);
    }

    await uiState.closeAddModSideDialog();
  }

  function isModInstalled(modId: string): boolean {
    return !!projectState.current.mods?.[modId];
  }
</script>

<Dialog.Root bind:open={uiState.addModDialogOpen}>
  <Dialog.Content
    class="max-w-lg max-h-9/12 overflow-y-auto"
    style="scrollbar-width: none; -ms-overflow-style: none"
  >
    <Dialog.Header>
      <Dialog.Title>Add Mod</Dialog.Title>
      <Dialog.Description>
        Search for mods on the selected platform and add them to the project.
      </Dialog.Description>
    </Dialog.Header>

    <div class="grid gap-4">
      <div class="grid gap-3">
        <Label for="platform">Platform</Label>
        <Select.Root type="single" bind:value={modSearchState.platform} onValueChange={handlePlatformChange}>
          <Select.Trigger class="w-[180px]">
            <span class="flex items-center gap-2">
              {#if modSearchState.platform === 'modrinth'}
                <ModrinthIcon />
                Modrinth
              {:else if modSearchState.platform === 'curseforge'}
                <CurseForgeIcon />
                Curseforge
              {:else}
                Select Platform
              {/if}
            </span>
          </Select.Trigger>
          <Select.Content>
            <Select.Item value="modrinth">
              <ModrinthIcon /> Modrinth
            </Select.Item>
            <Select.Item value="curseforge">
              <CurseForgeIcon /> Curseforge
            </Select.Item>
          </Select.Content>
        </Select.Root>
      </div>

      <div class="grid gap-3">
        <Label for="query">Query</Label>
        <form class="flex justify-between items-center" onsubmit={handleSearch}>
          <Input id="query" class="w-full" bind:value={modSearchState.query} autocomplete="off" required />
          <Button class="ml-2 cursor-pointer" disabled={modSearchState.isSearching} type="submit">
            <Search />
          </Button>
        </form>
      </div>

      {#if modSearchState.isSearching}
        <p>Searching...</p>
      {:else if modSearchState.results.length > 0}
        <ul class="mt-2">
          {#each modSearchState.results as result}
            <li class="p-2 border rounded mb-1 flex items-center justify-center">
              <div class="w-full">
                  <div class="flex items-center gap-2">
                      <button class="font-semibold flex items-center gap-1 cursor-pointer" onclick={() => BrowserOpenURL(result.URL)}>
                          {result.Name} <ExternalLink class="w-3 mb-1" />
                      </button>

                      <Select.Root type="single" bind:value={selectedVersions[result.ID]}>
                          <Select.Trigger class="py-1">
                              <span class="truncate max-w-[80px] block">{selectedVersions[result.ID] || result.Versions[0]}</span>
                          </Select.Trigger>
                          <Select.Content>
                              {#each result.Versions as version}
                                  <Select.Item value={version}>{version}</Select.Item>
                              {/each}
                          </Select.Content>
                      </Select.Root>
                  </div>

                <div class="text-sm text-muted-foreground">
                  {result.Description}
                </div>
                <div class="flex items-center gap-4 mt-1">
                  <div class="text-xs">Downloads: {result.Downloads}</div>
                  {#if modSearchState.platform === 'modrinth'}
                    <div class="flex items-center gap-2">
                      <Laptop
                        class={[
                          'w-5',
                          result.ClientSide === 'optional'
                            ? 'text-muted-foreground'
                            : result.ClientSide === 'required'
                              ? ''
                              : 'text-destructive'
                        ]}
                      />
                      <Server
                        class={[
                          'w-4',
                          result.ServerSide === 'optional'
                            ? 'text-muted-foreground'
                            : result.ServerSide === 'required'
                              ? ''
                              : 'text-destructive'
                        ]}
                      />
                    </div>
                  {/if}
                </div>
              </div>

              <Button class="ml-2 cursor-pointer" disabled={isModInstalled(result.ID) || loadingMods[result.ID]} onclick={() => handleAddMod(result)}>
                  {#if loadingMods[result.ID]}
                        <Loader class="w-4 h-4 animate-spin" />
                  {:else}
                      <Plus />
                  {/if}
              </Button>
            </li>
          {/each}
        </ul>
      {:else}
        <p class="text-muted-foreground">No results</p>
      {/if}
    </div>
  </Dialog.Content>
</Dialog.Root>

<Dialog.Root bind:open={uiState.addModSideDialogOpen}>
  <Dialog.Content class="max-w-sm">
    <Dialog.Header>
      <Dialog.Title>Choose Mod Side</Dialog.Title>
      <Dialog.Description>
        Select where this mod should be installed.
      </Dialog.Description>
    </Dialog.Header>

    <div class="grid gap-4">
      <div class="grid gap-3">
        <Label for="side">Side</Label>
        <Select.Root type="single" bind:value={selectedSide}>
          <Select.Trigger class="w-full">
            <span>
              {#if selectedSide === 'client'}
                Client
              {:else if selectedSide === 'server'}
                Server
              {:else if selectedSide === 'both'}
                Both
              {/if}
            </span>
          </Select.Trigger>
          <Select.Content>
            <Select.Item value="client">Client</Select.Item>
            <Select.Item value="server">Server</Select.Item>
            <Select.Item value="both">Both</Select.Item>
          </Select.Content>
        </Select.Root>
      </div>
    </div>

    <Dialog.Footer>
      <Button class="cursor-pointer" onclick={handleAddModWithSide}>
        Add Mod
      </Button>
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>
