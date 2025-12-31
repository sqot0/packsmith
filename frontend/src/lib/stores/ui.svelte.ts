import type { ModSearchResult } from '$lib/types/mod';
import { projectService } from '$lib/services';

// Dialog state management using Svelte 5 runes
class UIStore {
  addModDialogOpen = $state(false);
  newProjectDialogOpen = $state(false);
  addModSideDialogOpen = $state(false);
  updateModsDialogOpen = $state(false);
  changeVersionModDialogOpen = $state(false);
  logsDialogOpen = $state(false);
  importModsDialogOpen = $state(false);
  changeModSideDialogOpen = $state(false);
  logsContent = $state('');
  modToAdd = $state<ModSearchResult | null>(null);
  selectedModId = $state<string | null>(null);
  modToChangeSide = $state<string | null>(null);

  openAddModDialog = async () => {
    this.addModDialogOpen = true;
  }

  openNewProjectDialog = async () => {
      this.newProjectDialogOpen = true;
  }

  closeNewProjectDialog = async () => {
    this.newProjectDialogOpen = false;
  }

  openAddModSideDialog = async (mod: ModSearchResult) => {
    this.modToAdd = mod;
    this.addModSideDialogOpen = true;
  }

  closeAddModSideDialog = async () => {
      this.addModSideDialogOpen = false;
    this.modToAdd = null;
  }

  openUpdateModsDialog = async () => {
    this.updateModsDialogOpen = true;
  }

  closeUpdateModsDialog = async () => {
    this.updateModsDialogOpen = false;
  }

  openChangeVersionModDialog = async (modId: string) => {
    this.selectedModId = modId;
    this.changeVersionModDialogOpen = true;
  }

  closeChangeVersionModDialog = async () => {
      this.selectedModId = null;
      this.changeVersionModDialogOpen = false;
  }

  openChangeModSideDialog = async (modId: string) => {
    this.modToChangeSide = modId;
    this.changeModSideDialogOpen = true;
  }

  closeChangeModSideDialog = async () => {
    this.changeModSideDialogOpen = false;
    this.modToChangeSide = null;
  }

  openLogsDialog = async () => {
    this.logsContent = await projectService.getLogs();
    this.logsDialogOpen = true;
  }

  closeLogsDialog = async () => {
    this.logsDialogOpen = false;
    this.logsContent = '';
  }
}

export const uiState = new UIStore();
