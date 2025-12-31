import type { ModSearchResult } from '$lib/types/mod';
import * as modService from '$lib/services/mod-service';
import type { ModPlatform } from '$lib/services/mod-service';
import {toast} from "svelte-sonner";

// Mod search state management using Svelte 5 runes
class ModSearchStore {
  results = $state<ModSearchResult[]>([]);
  isSearching = $state(false);
  query = $state('');
  private _platform = $state<ModPlatform>('modrinth');

  get platform() {
    return this._platform;
  }

    set platform(value: ModPlatform) {
    this._platform = value;
    }

  search = async (searchQuery?: string) => {
    const q = searchQuery ?? this.query;
    if (!q) return;

    this.isSearching = true;
    try {
      this.results = await modService.searchMods(q, this._platform);
    } catch (error) {
      console.error('Failed to search mods:', error);
      toast.error('Failed to search for mod', { description: String(error) });
      this.results = [];
    } finally {
      this.isSearching = false;
    }
  }

  clearResults = async () => {
    this.results = [];
  }
}

export const modSearchState = new ModSearchStore();

