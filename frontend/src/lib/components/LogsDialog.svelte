<script lang="ts">
  import * as AlertDialog from '$lib/components/ui/alert-dialog';
  import { uiState } from '$lib/stores/ui.svelte';
  import { Button } from '$lib/components/ui/button';
</script>

<AlertDialog.Root bind:open={uiState.logsDialogOpen}>
    <AlertDialog.Content class="!h-[90vh] !w-[90vw] !max-w-none flex flex-col">
        <AlertDialog.Header>
            <AlertDialog.Title>Application Logs</AlertDialog.Title>
        </AlertDialog.Header>

        <!-- Основной скролл -->
        <div class="flex-1 overflow-auto bg-background border p-4 rounded font-mono text-sm">
      <pre class="whitespace-pre-wrap break-words">
{uiState.logsContent}
      </pre>
        </div>

        <AlertDialog.Footer>
            <Button variant="outline" onclick={uiState.closeLogsDialog}>
                Close
            </Button>
        </AlertDialog.Footer>
    </AlertDialog.Content>
</AlertDialog.Root>

