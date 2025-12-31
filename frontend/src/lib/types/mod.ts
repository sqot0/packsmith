export type ModSide = "client" | "server" | "both";

export interface Mod {
  source: string;
  side: string;
  version: string;
  url: string;
  filename: string;
  locked: boolean;
}

export interface ModSearchResult {
  ID: string;
  Name: string;
  Description: string;
  URL: string;
  Downloads: string;
  ClientSide: string;
  ServerSide: string;
  Versions: string[];
}

export interface AddModOptions {
  URL: string;
  Side: ModSide;
  Version: string;
}

export interface ModUpdateInfo {
    ModId: string;
    Version: string;
    URL: string;
}