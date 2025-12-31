import type { Mod } from './mod';

export type LoaderType = "forge" | "fabric" | "neoforge" | "quilt";

export interface ProjectConfig {
  name: string;
  minecraft: string;
  loader: LoaderType;
  mods: Record<string, Mod>;
}

export interface Project {
  path: string;
  name: string;
  minecraft: string;
  loader: LoaderType;
  mods: Record<string, Mod> | null;
}
