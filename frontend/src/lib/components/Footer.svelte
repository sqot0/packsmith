<script lang="ts">
  import { projectState } from '$lib/stores';
  import { uiState } from '$lib/stores/ui.svelte';
  import {Button} from "$lib/components/ui/button";

  let modCount = $derived(projectState.current.mods ? Object.keys(projectState.current.mods).length : 0);
</script>

<div class="fixed bg-background bottom-0 right-0 p-2 border-t w-full flex justify-between items-center">
    <Button variant="link" size="sm" class="cursor-pointer" onclick={uiState.openLogsDialog}>
        Logs
    </Button>

  <div class="flex w-full items-center justify-end gap-4 text-sm">
    <div>Mods: {modCount}</div>
    <div>Minecraft: {projectState.current.minecraft || 'N/A'}</div>
    <div>Loader: {projectState.current.loader}</div>
  </div>
</div>
