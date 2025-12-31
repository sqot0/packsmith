<script lang="ts">
    import './app.css';
    import Menu from '$lib/components/Menu.svelte';
    import ModsTable from "$lib/components/ModsTable.svelte";
    import Footer from "$lib/components/Footer.svelte";

    import NewProjectDialog from "$lib/components/NewProjectDialog.svelte";
    import AddModDialog from "$lib/components/AddModDialog.svelte";
    import UpdateModsDialog from "$lib/components/UpdateModsDialog.svelte";
    import ChangeVersionModDialog from "$lib/components/ChangeVersionModDialog.svelte";
    import LogsDialog from "$lib/components/LogsDialog.svelte";

    import { projectState } from "$lib/stores/project.svelte";
    import { Toaster } from "$lib/components/ui/sonner/index.js";
</script>

<NewProjectDialog />
<AddModDialog />
<UpdateModsDialog />
<ChangeVersionModDialog />
<LogsDialog />
<Toaster expand={true} position="bottom-right" />

<main>
    <Menu />

    {#if projectState.hasProject}
        <ModsTable />
        <Footer />
    {:else}
        <p class="text-muted-foreground w-full h-[calc(100vh-60px)] flex justify-center items-center">Choose project first</p>
    {/if}
</main>