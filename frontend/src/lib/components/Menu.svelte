<script lang="ts">
  import { onMount } from 'svelte';
  import * as Menubar from '$lib/components/ui/menubar';
  import { Folder, Download, CloudSync, FilePlusCorner, Loader, Settings, History } from '@lucide/svelte';
  import { projectState } from '$lib/stores/project.svelte';
  import * as modService from '$lib/services/mod-service';
  import { settingsService } from '$lib/services';
  import { uiState } from '$lib/stores/ui.svelte';
  import type { RecentProject } from '$lib/types/settings';
  import { toast } from "svelte-sonner";

  let isInstalling = $state(false);
  let recentProjects = $state<RecentProject[]>([]);

  async function loadRecentProjects() {
    try {
      recentProjects = await settingsService.getRecentProjects();
    } catch (err) {
      console.error('Failed to load recent projects:', err);
    }
  }

  onMount(() => {
    loadRecentProjects();
  });

  $effect(() => {
    if (projectState.current.path) {
      loadRecentProjects();
    }
  });

  async function handleOpenRecent(path: string) {
    try {
      await projectState.openRecentProject(path);
      toast.success("Project opened successfully");
      await loadRecentProjects();
    } catch (error) {
      console.error('Failed to open recent project:', error);
      toast.error('Failed to open project', { description: String(error) });
    }
  }

  async function handleInstallMods() {
    if (isInstalling) return;
    isInstalling = true;
    try {
      await modService.installMods();
      toast.success("Mods installed successfully");
    } catch (error) {
      console.error('Failed to install mods:', error);
      toast.error('Failed to install mods', { description: String(error) });
    } finally {
      isInstalling = false;
    }
  }
</script>

<header class="w-full mb-5">
  <Menubar.Root class="rounded-none">
    <Menubar.Menu>
      <Menubar.Trigger class="cursor-pointer">
        <Folder size="16" class="mr-1" />
        Open
      </Menubar.Trigger>
      <Menubar.Content>
        <Menubar.Item onclick={projectState.selectAndOpenProject} class="cursor-pointer">
          <Folder size="16" class="mr-2" />
          Browse folder...
        </Menubar.Item>
        <Menubar.Sub>
          <Menubar.SubTrigger class="cursor-pointer">
            <History size="16" class="mr-2" />
            Open Recent
          </Menubar.SubTrigger>
          <Menubar.SubContent class="w-72">
            {#if recentProjects.length === 0}
              <Menubar.Item disabled>No recent projects</Menubar.Item>
            {:else}
              {#each recentProjects as project}
                <Menubar.Item
                  onclick={() => handleOpenRecent(project.path)}
                  class="cursor-pointer flex flex-col items-start gap-0.5 py-1.5"
                >
                  <span class="font-medium truncate max-w-full">{project.name || 'Untitled'}</span>
                  <span class="text-xs text-muted-foreground truncate max-w-full font-mono">{project.path}</span>
                </Menubar.Item>
              {/each}
            {/if}
          </Menubar.SubContent>
        </Menubar.Sub>
      </Menubar.Content>
    </Menubar.Menu>

    <Menubar.Menu>
      <Menubar.Trigger disabled={!projectState.hasProject} onclick={uiState.openAddModDialog}>
        <FilePlusCorner size="16" class="mr-1" />
        Add mod
      </Menubar.Trigger>
    </Menubar.Menu>

    <Menubar.Menu>
      <Menubar.Trigger disabled={!projectState.hasProject} onclick={uiState.openUpdateModsDialog}>
        <CloudSync size="16" class="mr-1" />
        Update mods
      </Menubar.Trigger>
    </Menubar.Menu>

    <Menubar.Menu>
      <Menubar.Trigger disabled={!projectState.hasProject || isInstalling} onclick={handleInstallMods}>
        {#if isInstalling}
          <Loader size="16" class="mr-1 animate-spin" />
          Installing mods...
        {:else}
          <Download size="16" class="mr-1" />
          Install mods
        {/if}
      </Menubar.Trigger>
    </Menubar.Menu>

    <Menubar.Menu>
      <Menubar.Trigger onclick={uiState.openSettingsDialog} class="cursor-pointer">
        <Settings size="16" class="mr-1" />
        Settings
      </Menubar.Trigger>
    </Menubar.Menu>
  </Menubar.Root>
</header>
