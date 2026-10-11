<script lang="ts">
  import * as Dialog from '$lib/components/ui/dialog';
  import { Button } from '$lib/components/ui/button';
  import { Input } from '$lib/components/ui/input';
  import { Label } from '$lib/components/ui/label';
  import { uiState } from '$lib/stores/ui.svelte';
  import { settingsService } from '$lib/services';
  import type { RecentProject } from '$lib/types/settings';
  import { toast } from 'svelte-sonner';
  import { Loader, Trash2, Settings, Cpu, History } from '@lucide/svelte';

  let maxWorkers = $state(15);
  let recentProjects = $state<RecentProject[]>([]);
  let loading = $state(false);
  let saving = $state(false);

  $effect(() => {
    if (uiState.settingsDialogOpen) {
      loadSettings();
    }
  });

  async function loadSettings() {
    loading = true;
    try {
      const cfg = await settingsService.getAppSettings();
      maxWorkers = cfg.maxWorkers || 15;
      recentProjects = cfg.recentProjects || [];
    } catch (error) {
      console.error('Failed to load settings:', error);
      toast.error('Failed to load settings', { description: String(error) });
    } finally {
      loading = false;
    }
  }

  async function handleSave(event: Event) {
    event.preventDefault();
    if (saving) return;

    if (maxWorkers < 1) {
      toast.error('Max workers must be at least 1');
      return;
    }

    saving = true;
    try {
      await settingsService.saveAppSettings({
        maxWorkers: Number(maxWorkers),
        recentProjects,
      });
      toast.success('Settings saved successfully');
      await uiState.closeSettingsDialog();
    } catch (error) {
      console.error('Failed to save settings:', error);
      toast.error('Failed to save settings', { description: String(error) });
    } finally {
      saving = false;
    }
  }

  async function handleRemoveRecent(path: string) {
    try {
      await settingsService.removeRecentProject(path);
      recentProjects = recentProjects.filter((p) => p.path !== path);
      toast.success('Project removed from history');
    } catch (error) {
      console.error('Failed to remove recent project:', error);
      toast.error('Failed to remove project', { description: String(error) });
    }
  }

  async function handleClearRecent() {
    try {
      await settingsService.clearRecentProjects();
      recentProjects = [];
      toast.success('Recent projects cleared');
    } catch (error) {
      console.error('Failed to clear recent projects:', error);
      toast.error('Failed to clear recent projects', { description: String(error) });
    }
  }
</script>

<Dialog.Root bind:open={uiState.settingsDialogOpen}>
  <Dialog.Content class="sm:max-w-[540px]">
    <Dialog.Header>
      <Dialog.Title class="flex items-center gap-2">
        <Settings class="w-5 h-5" />
        Application Settings
      </Dialog.Title>
      <Dialog.Description>
        Configure global Packsmith preferences and manage recent projects.
      </Dialog.Description>
    </Dialog.Header>

    {#if loading}
      <div class="flex items-center justify-center py-10 text-muted-foreground gap-2">
        <Loader class="w-5 h-5 animate-spin" />
        <span>Loading settings...</span>
      </div>
    {:else}
      <form onsubmit={handleSave} class="space-y-6">
        <div class="space-y-4">
          <div class="space-y-2">
            <Label for="maxWorkers" class="flex items-center gap-2 text-sm font-semibold">
              <Cpu class="w-4 h-4 text-primary" />
              Max Worker Threads
            </Label>
            <Input
              id="maxWorkers"
              type="number"
              min="1"
              max="64"
              bind:value={maxWorkers}
              required
              class="w-full"
            />
            <p class="text-xs text-muted-foreground">
              Maximum concurrent worker goroutines for file downloading and copying in installer and updater (default: 15).
            </p>
          </div>

          <div class="space-y-2 pt-2 border-t">
            <div class="flex items-center justify-between">
              <Label class="flex items-center gap-2 text-sm font-semibold">
                <History class="w-4 h-4 text-primary" />
                Recent Projects ({recentProjects.length})
              </Label>
              {#if recentProjects.length > 0}
                <Button
                  type="button"
                  variant="ghost"
                  size="sm"
                  onclick={handleClearRecent}
                  class="text-xs text-destructive hover:text-destructive hover:bg-destructive/10 h-7 px-2"
                >
                  Clear all
                </Button>
              {/if}
            </div>

            {#if recentProjects.length === 0}
              <p class="text-xs text-muted-foreground py-2">No recent projects yet.</p>
            {:else}
              <div class="max-h-48 overflow-y-auto space-y-1.5 pr-1 border rounded-md p-2 bg-muted/20">
                {#each recentProjects as project}
                  <div class="flex items-center justify-between p-2 rounded bg-background hover:bg-muted/50 border text-xs gap-2">
                    <div class="overflow-hidden min-w-0 flex-1">
                      <p class="font-medium truncate">{project.name || 'Untitled'}</p>
                      <p class="text-muted-foreground truncate font-mono text-[11px]">{project.path}</p>
                    </div>
                    <Button
                      type="button"
                      variant="ghost"
                      size="icon"
                      onclick={() => handleRemoveRecent(project.path)}
                      class="h-7 w-7 text-muted-foreground hover:text-destructive shrink-0"
                      title="Remove from history"
                    >
                      <Trash2 class="w-3.5 h-3.5" />
                    </Button>
                  </div>
                {/each}
              </div>
            {/if}
          </div>
        </div>

        <Dialog.Footer class="gap-2 sm:gap-0">
          <Button
            type="button"
            variant="outline"
            onclick={uiState.closeSettingsDialog}
          >
            Cancel
          </Button>
          <Button
            type="submit"
            disabled={saving || loading}
            class="gap-2 ml-2"
          >
            {#if saving}
              <Loader class="w-4 h-4 animate-spin" />
              Saving...
            {:else}
              Save Changes
            {/if}
          </Button>
        </Dialog.Footer>
      </form>
    {/if}
  </Dialog.Content>
</Dialog.Root>
