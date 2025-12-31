import {
    SearchMods,
    AddMod,
    RemoveMod,
    CheckModsUpdates,
    UpdateMods,
    InstallMods,
    ChangeModSide,
    ChangeModLocked,
    ChangeModVersion, GetModVersions,
} from '$backend';
import type {ModSearchResult, ModSide, AddModOptions, ModUpdateInfo} from '$lib/types/mod';

export type ModPlatform = 'modrinth' | 'curseforge';

/**
 * Searches for mods on the specified platform
 */
export async function searchMods(
  query: string,
  platform: ModPlatform
): Promise<ModSearchResult[]> {
  const results = await SearchMods(query, platform);
  return results ?? [];
}

/**
 * Adds a mod to the current project
 */
export async function addMod(
  modId: string,
  platform: ModPlatform,
  options: AddModOptions
): Promise<void> {
  await AddMod(modId, platform, options);
}

export async function removeMod(
    modId: string,
): Promise<void> {
    await RemoveMod(modId);
}

export async function checkModsUpdates(modIds: string[]): Promise<ModUpdateInfo[]> {
    return await CheckModsUpdates(modIds);
}

export async function updateMods(modsToUpdate: ModUpdateInfo[]): Promise<void> {
    return await UpdateMods(modsToUpdate);
}

export async function installMods(): Promise<void> {
    return await InstallMods();
}

export async function changeModSide(modId: string, side: ModSide): Promise<void> {
    return await ChangeModSide(modId, side);
}

export async function changeModLocked(modId: string, locked: boolean): Promise<void> {
    return await ChangeModLocked(modId, locked);
}

export async function getModVersions(modId: string): Promise<string[]> {
    return await GetModVersions(modId)
}

export async function changeModVersion(modId: string, version: string): Promise<void> {
    return await ChangeModVersion(modId, version);
}

/**
 * Determines the mod side based on platform-specific requirements
 */
export function determineModSide(
  platform: ModPlatform,
  clientSide?: string,
  serverSide?: string
): ModSide {
  if (platform !== 'modrinth') {
    return 'both';
  }

  if (clientSide === 'required' && serverSide === 'unsupported') {
    return 'client';
  }
  
  if (serverSide === 'required' && clientSide === 'unsupported') {
    return 'server';
  }
  
  return 'both';
}

/**
 * Checks if a mod requires side selection (for Modrinth mods with optional sides)
 */
export function requiresSideSelection(
  platform: ModPlatform,
  clientSide?: string,
  serverSide?: string
): boolean {
  return platform === 'modrinth' && 
    (clientSide === 'optional' || serverSide === 'optional');
}
