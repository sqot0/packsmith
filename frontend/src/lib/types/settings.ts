export interface RecentProject {
  path: string;
  name: string;
  lastOpened: string;
}

export interface AppSettings {
  maxWorkers: number;
  recentProjects: RecentProject[];
}
