import { 
  SelectProjectDirectory, 
  OpenProject, 
  InitializeProject, GetLogs
} from '$backend';
import type {LoaderType, ProjectConfig} from '$lib/types/project';

/**
 * Opens a directory selection dialog and returns the selected path
 */
export async function selectProjectDirectory(): Promise<string> {
  return await SelectProjectDirectory();
}

export async function getLogs(): Promise<string> {
    return await GetLogs()
}

/**
 * Opens an existing project from the given path
 */
export async function openProject(projectPath: string): Promise<ProjectConfig> {
  const resp = await OpenProject(projectPath);

    return {
        name: resp.name,
        minecraft: resp.minecraft,
        loader: resp.loader as LoaderType,
        mods: resp.mods,
    };
}

/**
 * Initializes a new project at the given path
 */
export async function initializeProject(
  projectPath: string,
  name: string,
  minecraftVersion: string,
  modLoader: string
): Promise<void> {
  await InitializeProject(projectPath, name, minecraftVersion, modLoader);
}
