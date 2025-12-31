<script lang="ts">
  import * as Table from '$lib/components/ui/table';
  import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
  import * as Dialog from '$lib/components/ui/dialog';
  import * as Select from '$lib/components/ui/select';
  import { Label } from '$lib/components/ui/label';
  import { Input } from '$lib/components/ui/input';
  import { ExternalLink, Laptop, Server, Trash, Lock, LockOpen, CircleArrowUp, ListStart, Ellipsis } from '@lucide/svelte';
  import CurseforgeIcon from '$lib/icons/CurseforgeIcon.svelte';
  import ModrinthIcon from '$lib/icons/ModrinthIcon.svelte';
  import { BrowserOpenURL } from '$runtime';
  import { projectState } from '$lib/stores/project.svelte';
  import { modService } from "$lib/services";
  import { uiState } from '$lib/stores/ui.svelte';
  import type { ModSide } from '$lib/types/mod';
  import {toast} from "svelte-sonner";
  import {Button} from "$lib/components/ui/button";

  let selectedSide = $state<ModSide>('both');
  let searchQuery = $state('');

  const modsArray = $derived(
    projectState.current.mods ? Object.entries(projectState.current.mods) : []
  );

  const filteredMods = $derived(
    modsArray.filter(([id, mod]) =>
      id.toLowerCase().includes(searchQuery.toLowerCase()) ||
      (mod?.filename?.toLowerCase().includes(searchQuery.toLowerCase()) ?? false)
    )
  );

  function openChangeSideDialog(id: string, side: string) {
    selectedSide = side as ModSide;
    uiState.openChangeModSideDialog(id);
  }

  async function handleChangeSide() {
    try {
      await modService.changeModSide(uiState.modToChangeSide!, selectedSide);
      await projectState.refreshProject();
      toast.success('Mod side changed successfully');
      await uiState.closeChangeModSideDialog();
    } catch (error) {
      toast.error('Failed to change mod side', { description: String(error) });
    }
  }
</script>

<div class="px-4 mb-4">
  <Input placeholder="Search mods..." bind:value={searchQuery} />
</div>

<Table.Root class="px-4 table-scroll mb-10">
  <Table.Header>
    <Table.Row>
      <Table.Head class="pl-10">Mod</Table.Head>
      <Table.Head>Filename</Table.Head>
      <Table.Head>Version</Table.Head>
      <Table.Head>Side</Table.Head>
    </Table.Row>
  </Table.Header>

  {#if filteredMods.length === 0}
    <p class="text-muted-foreground mt-3 ml-2">There are no mods yet</p>
  {/if}

  <Table.Body>
    {#each filteredMods as [id, mod]}
      <Table.Row>
        <Table.Cell onclick={() => BrowserOpenURL(mod?.source)}>
          <div class="flex items-center gap-1 cursor-pointer hover:underline">
            {#if mod?.source?.startsWith('https://www.curseforge.com')}
              <CurseforgeIcon class="w-5 mr-2" />
            {:else if mod?.source?.startsWith('https://modrinth.com')}
              <ModrinthIcon class="w-5 h-5 mr-2" />
            {/if}
            <span>{id}</span>
            <ExternalLink class="w-3 mb-1" />
          </div>
        </Table.Cell>

        <Table.Cell>{mod?.filename}</Table.Cell>
        <Table.Cell>
            {#if mod?.locked}
                <Lock class="w-4 inline-block mr-1" />
            {/if}
            {mod?.version}
        </Table.Cell>
        <Table.Cell>
          <div class="flex items-center justify-between pr-5 gap-2">
            <div class="flex items-center gap-2 h-full">
              {#if mod?.side === 'client'}
                <Laptop class="w-5" />
              {:else if mod?.side === 'server'}
                <Server class="w-5" />
              {:else}
                <Laptop class="w-5" />
                <Server class="w-5" />
              {/if}
            </div>

            <DropdownMenu.Root>
              <DropdownMenu.Trigger>
                <Ellipsis class="w-5" />
              </DropdownMenu.Trigger>
              <DropdownMenu.Content>
                <DropdownMenu.Group>
                  <DropdownMenu.Item onclick={() => uiState.openChangeVersionModDialog(id)}>
                    <CircleArrowUp class="w-4" /> Change Version
                  </DropdownMenu.Item>
                    <DropdownMenu.Item onclick={() => openChangeSideDialog(id, mod?.side || 'both')}>
                      <ListStart class="w-4" />Change Side
                    </DropdownMenu.Item>
                  <DropdownMenu.Item onclick={async () => {
                      try {
                        await modService.changeModLocked(id, !mod?.locked)
                        await projectState.refreshProject()
                        toast.success('Mod lock changed successfully');
                      } catch (error) {
                        toast.error('Failed to change mod lock', { description: String(error) });
                      }
                   }}>
                      {#if !mod?.locked}
                          <Lock class="w-4" />
                            Lock Version
                        {:else}
                            <LockOpen class="w-4" />
                                Unlock Version
                        {/if}
                  </DropdownMenu.Item>
                  <DropdownMenu.Item onclick={async () => {
                      try {
                        await modService.removeMod(id)
                        await projectState.refreshProject()
                        toast.success('Mod removed successfully');
                      } catch (error) {
                        toast.error('Failed to remove mod', { description: String(error) });
                      }
                  }}>
                    <Trash class="w-4" /> Delete
                  </DropdownMenu.Item>
                </DropdownMenu.Group>
              </DropdownMenu.Content>
            </DropdownMenu.Root>
          </div>
        </Table.Cell>
      </Table.Row>
    {/each}
  </Table.Body>
</Table.Root>

<Dialog.Root bind:open={uiState.changeModSideDialogOpen}>
  <Dialog.Content class="max-w-sm">
    <Dialog.Header>
      <Dialog.Title>Change Mod Side</Dialog.Title>
      <Dialog.Description>
        Select where this mod should be installed.
      </Dialog.Description>
    </Dialog.Header>

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

    <Dialog.Footer>
      <Button class="cursor-pointer" onclick={handleChangeSide}>
        Change Side
      </Button>
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>
