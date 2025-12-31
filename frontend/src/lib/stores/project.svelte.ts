import type { Project } from '$lib/types/project';
import { uiState } from './ui.svelte';
import * as projectService from '$lib/services/project-service';

const EMPTY_PROJECT: Project = {
  path: '',
  name: '',
  minecraft: '',
  loader: 'forge',
  mods: null,
};

// Project state management using Svelte 5 runes
class ProjectStore {
    current = $state<Project>(EMPTY_PROJECT);

    get hasProject() { return !!this.current.name; }
    get hasMods() { return !!this.current.mods; }

    selectAndOpenProject = async () => {
        const projectPath = await projectService.selectProjectDirectory();
        if (!projectPath) return;

        this.current.path = projectPath;

        try {
            const projectConfig = await projectService.openProject(projectPath);
            this.current = {
                path: projectPath,
                ...projectConfig
            };
        } catch {
            await uiState.openNewProjectDialog();
        }
    };

    createProject = async (name: string, minecraftVersion: string, modLoader: string) => {
        await projectService.initializeProject(this.current.path, name, minecraftVersion, modLoader);
        const projectConfig = await projectService.openProject(this.current.path);
        this.current = {
            path: this.current.path,
            ...projectConfig
        };
        await uiState.closeNewProjectDialog();
    };

    refreshProject = async () => {
        if (!this.current.path) return;
        try {
        const projectConfig = await projectService.openProject(this.current.path);
            this.current = {
                path: this.current.path,
                ...projectConfig
            };
        } catch {
            await uiState.openNewProjectDialog();
        }
    }
}


export const projectState = new ProjectStore();

