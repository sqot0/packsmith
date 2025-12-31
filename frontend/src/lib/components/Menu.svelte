<script lang="ts">
  import * as Menubar from '$lib/components/ui/menubar';
  import { Folder, Download, CloudSync, FilePlusCorner } from '@lucide/svelte';
  import { projectState } from '$lib/stores/project.svelte';
    import * as modService from '$lib/services/mod-service';
  import { uiState } from '$lib/stores/ui.svelte';
  import {toast} from "svelte-sonner";
</script>

<header class="w-full mb-5">
  <Menubar.Root class="rounded-none">
    <Menubar.Menu>
      <Menubar.Trigger onclick={projectState.selectAndOpenProject}>
        <Folder size="16" class="mr-1" />
        Open
      </Menubar.Trigger>
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
      <Menubar.Trigger disabled={!projectState.hasProject} onclick={async () => {
          try {
                await modService.installMods()
                toast.success("Mods installed successfully");
          } catch (error) {
                console.error('Failed to install mods:', error);
                toast.error('Failed to install mods', { description: String(error) });
          }
      }}>
        <Download size="16" class="mr-1" />
        Install mods
      </Menubar.Trigger>
    </Menubar.Menu>

<!--    <Menubar.Menu>-->
<!--      <Menubar.Trigger-->
<!--        disabled={!projectState.hasProject}-->
<!--        class="data-[state=open]:bg-accent data-[state=open]:text-accent-foreground"-->
<!--      >-->
<!--        <Settings size="16" class="mr-1" />-->
<!--        Settings-->
<!--      </Menubar.Trigger>-->
<!--      <Menubar.Content>-->
<!--        <Menubar.Item inset>Change Config</Menubar.Item>-->
<!--      </Menubar.Content>-->
<!--    </Menubar.Menu>-->
  </Menubar.Root>
</header>
