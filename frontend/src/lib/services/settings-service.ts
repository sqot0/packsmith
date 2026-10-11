import {
  GetAppSettings,
  SaveAppSettings,
  GetRecentProjects,
  RemoveRecentProject,
  ClearRecentProjects
} from '$backend';
import type { AppSettings, RecentProject } from '$lib/types/settings';
import { appconfig } from '../../../wailsjs/go/models';

export async function getAppSettings(): Promise<AppSettings> {
  const settings = await GetAppSettings();
  return {
    maxWorkers: settings.maxWorkers || 15,
    recentProjects: (settings.recentProjects || []).map(p => ({
      path: p.path,
      name: p.name,
      lastOpened: p.lastOpened
    }))
  };
}

export async function saveAppSettings(settings: AppSettings): Promise<void> {
  const model = new appconfig.AppSettings({
    maxWorkers: settings.maxWorkers,
    recentProjects: (settings.recentProjects || []).map(p => new appconfig.RecentProject(p))
  });
  await SaveAppSettings(model);
}

export async function getRecentProjects(): Promise<RecentProject[]> {
  const projects = await GetRecentProjects();
  return (projects || []).map(p => ({
    path: p.path,
    name: p.name,
    lastOpened: p.lastOpened
  }));
}

export async function removeRecentProject(path: string): Promise<void> {
  await RemoveRecentProject(path);
}

export async function clearRecentProjects(): Promise<void> {
  await ClearRecentProjects();
}
