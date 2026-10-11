package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sqot0/packsmith/backend/internal/logger"
	"github.com/sqot0/packsmith/backend/internal/sources"
	"github.com/sqot0/packsmith/backend/internal/updater"
)

const mcpInstructions = `Packsmith is an open-source Minecraft modpack manager supporting Forge, Fabric, NeoForge, and Quilt mod loaders.

Recommended AI workflow:
1. Project Loading: Call 'open_project' with the directory path containing packsmith.json, or supply 'project_path' directly in other tool calls. Call 'get_project_config' to view the Minecraft version, mod loader, and current list of mods.
2. Mod Discovery: Call 'search_mods' with query and platform ('modrinth' or 'curseforge') to find mods matching the modpack's Minecraft version and loader.
3. Adding Mods: Call 'add_mod' with mod_id, platform, version, and execution side ('client', 'server', or 'both'). Local .jar files can be added via 'add_local_mod'.
4. Mod Configuration: Use 'change_mod_side', 'change_mod_locked', or 'change_mod_version' to adjust mod settings.
5. Updates: Call 'check_mods_updates' to see available updates, and 'update_mods' to apply them.
6. Build & Distribution: Call 'install_mods' to clean and prepare 'client' and 'server' distribution folders with all configured mods.`

func (a *App) ensureProject(projectPath string) error {
	if projectPath != "" && projectPath != a.ProjectPath {
		_, err := a.OpenProject(projectPath)
		return err
	}
	if a.ProjectPath == "" {
		return fmt.Errorf("no project is currently open; specify 'project_path' or call 'open_project' first")
	}
	return nil
}

func (a *App) NewMCPServer() *server.MCPServer {
	logger.Init()

	s := server.NewMCPServer(
		"Packsmith",
		"1.0.0",
		server.WithInstructions(mcpInstructions),
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
	)

	a.registerProjectTools(s)
	a.registerModTools(s)
	a.registerSettingsTools(s)
	a.registerResources(s)
	a.registerPrompts(s)

	return s
}

func (a *App) ServeMCP() error {
	s := a.NewMCPServer()
	return server.ServeStdio(s)
}

func (a *App) registerProjectTools(s *server.MCPServer) {
	// Tool: open_project
	s.AddTool(
		mcp.NewTool(
			"open_project",
			mcp.WithDescription("Open an existing Packsmith Minecraft modpack project directory"),
			mcp.WithString("project_path", mcp.Required(), mcp.Description("Path to modpack project directory containing packsmith.json")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			path, err := request.RequireString("project_path")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			cfg, err := a.OpenProject(path)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to open project: %v", err)), nil
			}
			data, _ := json.MarshalIndent(cfg, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// Tool: get_project_config
	s.AddTool(
		mcp.NewTool(
			"get_project_config",
			mcp.WithDescription("Get configuration and installed mods for the active modpack project (or specified project_path)"),
			mcp.WithString("project_path", mcp.Description("Optional path to modpack project. If omitted, uses active project")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			path := request.GetString("project_path", "")
			if err := a.ensureProject(path); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			cfg, err := a.loadConfig()
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to load project config: %v", err)), nil
			}
			data, _ := json.MarshalIndent(cfg, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// Tool: initialize_project
	s.AddTool(
		mcp.NewTool(
			"initialize_project",
			mcp.WithDescription("Initialize a new Packsmith Minecraft modpack project in a given directory"),
			mcp.WithString("project_path", mcp.Required(), mcp.Description("Directory path where project should be created")),
			mcp.WithString("name", mcp.Required(), mcp.Description("Name of the modpack project")),
			mcp.WithString("minecraft", mcp.Required(), mcp.Description("Minecraft version (e.g. '1.20.1', '1.21.1')")),
			mcp.WithString("loader", mcp.Required(), mcp.Description("Mod loader ('forge', 'fabric', 'neoforge', 'quilt')"), mcp.Enum("forge", "fabric", "neoforge", "quilt")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			path, err := request.RequireString("project_path")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			name, err := request.RequireString("name")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			mc, err := request.RequireString("minecraft")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			loader, err := request.RequireString("loader")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			if err := a.InitializeProject(path, name, mc, loader); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to initialize project: %v", err)), nil
			}
			cfg, err := a.OpenProject(path)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Project initialized, but failed to open: %v", err)), nil
			}
			data, _ := json.MarshalIndent(cfg, "", "  ")
			return mcp.NewToolResultText(fmt.Sprintf("Project initialized and opened successfully:\n%s", string(data))), nil
		},
	)

	// Tool: get_logs
	s.AddTool(
		mcp.NewTool(
			"get_logs",
			mcp.WithDescription("Get recent Packsmith application log entries"),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			logs, err := a.GetLogs()
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to read logs: %v", err)), nil
			}
			return mcp.NewToolResultText(logs), nil
		},
	)
}

func (a *App) registerSettingsTools(s *server.MCPServer) {
	// Tool: get_app_settings
	s.AddTool(
		mcp.NewTool(
			"get_app_settings",
			mcp.WithDescription("Get Packsmith global application settings including maxWorkers and recent projects"),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			cfg, err := a.GetAppSettings()
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to get app settings: %v", err)), nil
			}
			data, _ := json.MarshalIndent(cfg, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// Tool: save_app_settings
	s.AddTool(
		mcp.NewTool(
			"save_app_settings",
			mcp.WithDescription("Update Packsmith global application settings"),
			mcp.WithNumber("max_workers", mcp.Description("Maximum worker goroutines for concurrent file operations (e.g. 15)")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			cfg, err := a.GetAppSettings()
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to load current settings: %v", err)), nil
			}
			type savePayload struct {
				MaxWorkers int `json:"max_workers"`
			}
			var payload savePayload
			_ = request.BindArguments(&payload)
			if payload.MaxWorkers > 0 {
				cfg.MaxWorkers = payload.MaxWorkers
			}
			if err := a.SaveAppSettings(cfg); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to save app settings: %v", err)), nil
			}
			data, _ := json.MarshalIndent(cfg, "", "  ")
			return mcp.NewToolResultText(fmt.Sprintf("Settings saved successfully:\n%s", string(data))), nil
		},
	)

	// Tool: get_recent_projects
	s.AddTool(
		mcp.NewTool(
			"get_recent_projects",
			mcp.WithDescription("Get the list of recently opened or initialized modpack projects"),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			projects, err := a.GetRecentProjects()
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to get recent projects: %v", err)), nil
			}
			data, _ := json.MarshalIndent(projects, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// Tool: clear_recent_projects
	s.AddTool(
		mcp.NewTool(
			"clear_recent_projects",
			mcp.WithDescription("Clear all recent modpack projects history"),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			if err := a.ClearRecentProjects(); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to clear recent projects: %v", err)), nil
			}
			return mcp.NewToolResultText("Recent projects cleared successfully"), nil
		},
	)
}

func (a *App) registerModTools(s *server.MCPServer) {
	// Tool: search_mods
	s.AddTool(
		mcp.NewTool(
			"search_mods",
			mcp.WithDescription("Search for Minecraft mods on Modrinth or CurseForge compatible with current project's Minecraft version and mod loader"),
			mcp.WithString("query", mcp.Required(), mcp.Description("Search terms or mod name (e.g. 'jei', 'sodium', 'waystones')")),
			mcp.WithString("platform", mcp.Required(), mcp.Description("Platform to search on ('modrinth' or 'curseforge')"), mcp.Enum("modrinth", "curseforge")),
			mcp.WithString("project_path", mcp.Description("Optional path to modpack project. If omitted, uses active project")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			path := request.GetString("project_path", "")
			if err := a.ensureProject(path); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			query, err := request.RequireString("query")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			platform, err := request.RequireString("platform")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			results, err := a.SearchMods(query, strings.ToLower(platform))
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to search mods: %v", err)), nil
			}
			data, _ := json.MarshalIndent(results, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// Tool: get_mod_versions
	s.AddTool(
		mcp.NewTool(
			"get_mod_versions",
			mcp.WithDescription("Get compatible versions for a specific mod in the active modpack project"),
			mcp.WithString("mod_id", mcp.Required(), mcp.Description("Mod ID or slug (e.g. 'sodium', 'jei')")),
			mcp.WithString("project_path", mcp.Description("Optional path to modpack project. If omitted, uses active project")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			path := request.GetString("project_path", "")
			if err := a.ensureProject(path); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			modID, err := request.RequireString("mod_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			versions, err := a.GetModVersions(modID)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to get versions for mod %s: %v", modID, err)), nil
			}
			data, _ := json.MarshalIndent(versions, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// Tool: add_mod
	s.AddTool(
		mcp.NewTool(
			"add_mod",
			mcp.WithDescription("Download and add a mod from Modrinth or CurseForge to the modpack project"),
			mcp.WithString("mod_id", mcp.Required(), mcp.Description("Mod ID or slug (e.g. 'jei', 'sodium')")),
			mcp.WithString("platform", mcp.Required(), mcp.Description("Platform source ('modrinth' or 'curseforge')"), mcp.Enum("modrinth", "curseforge")),
			mcp.WithString("version", mcp.Required(), mcp.Description("Specific version ID or string of the mod")),
			mcp.WithString("side", mcp.Required(), mcp.Description("Execution environment side ('client', 'server', or 'both')"), mcp.Enum("client", "server", "both")),
			mcp.WithString("url", mcp.Description("Optional web URL for the mod")),
			mcp.WithString("project_path", mcp.Description("Optional path to modpack project. If omitted, uses active project")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			path := request.GetString("project_path", "")
			if err := a.ensureProject(path); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			modID, err := request.RequireString("mod_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			platform, err := request.RequireString("platform")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			version, err := request.RequireString("version")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			side, err := request.RequireString("side")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			url := request.GetString("url", "")

			meta := sources.ModMetaData{
				URL:     url,
				Side:    side,
				Version: version,
			}
			if err := a.AddMod(modID, strings.ToLower(platform), meta); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to add mod %s: %v", modID, err)), nil
			}
			return mcp.NewToolResultText(fmt.Sprintf("Mod '%s' (version %s, side %s) added successfully", modID, version, side)), nil
		},
	)

	// Tool: add_local_mod
	s.AddTool(
		mcp.NewTool(
			"add_local_mod",
			mcp.WithDescription("Add a local .jar mod file from disk to the modpack project"),
			mcp.WithString("file_path", mcp.Required(), mcp.Description("Path to local .jar mod file")),
			mcp.WithString("side", mcp.Required(), mcp.Description("Execution environment side ('client', 'server', or 'both')"), mcp.Enum("client", "server", "both")),
			mcp.WithString("project_path", mcp.Description("Optional path to modpack project. If omitted, uses active project")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			path := request.GetString("project_path", "")
			if err := a.ensureProject(path); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			filePath, err := request.RequireString("file_path")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			side, err := request.RequireString("side")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			if err := a.AddLocalMod(filePath, side); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to add local mod from %s: %v", filePath, err)), nil
			}
			return mcp.NewToolResultText(fmt.Sprintf("Local mod '%s' added successfully with side '%s'", filePath, side)), nil
		},
	)

	// Tool: remove_mod
	s.AddTool(
		mcp.NewTool(
			"remove_mod",
			mcp.WithDescription("Remove a mod from the modpack project and delete its cached file"),
			mcp.WithString("mod_id", mcp.Required(), mcp.Description("Mod ID to remove")),
			mcp.WithString("project_path", mcp.Description("Optional path to modpack project. If omitted, uses active project")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			path := request.GetString("project_path", "")
			if err := a.ensureProject(path); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			modID, err := request.RequireString("mod_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			if err := a.RemoveMod(modID); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to remove mod %s: %v", modID, err)), nil
			}
			return mcp.NewToolResultText(fmt.Sprintf("Mod '%s' removed successfully", modID)), nil
		},
	)

	// Tool: change_mod_side
	s.AddTool(
		mcp.NewTool(
			"change_mod_side",
			mcp.WithDescription("Change execution side ('client', 'server', or 'both') of an existing mod"),
			mcp.WithString("mod_id", mcp.Required(), mcp.Description("Mod ID to update")),
			mcp.WithString("side", mcp.Required(), mcp.Description("New execution side ('client', 'server', or 'both')"), mcp.Enum("client", "server", "both")),
			mcp.WithString("project_path", mcp.Description("Optional path to modpack project. If omitted, uses active project")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			path := request.GetString("project_path", "")
			if err := a.ensureProject(path); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			modID, err := request.RequireString("mod_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			side, err := request.RequireString("side")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			if err := a.ChangeModSide(modID, side); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to change side for %s: %v", modID, err)), nil
			}
			return mcp.NewToolResultText(fmt.Sprintf("Side for mod '%s' changed to '%s'", modID, side)), nil
		},
	)

	// Tool: change_mod_locked
	s.AddTool(
		mcp.NewTool(
			"change_mod_locked",
			mcp.WithDescription("Lock or unlock a mod to exclude or include it in automatic updates"),
			mcp.WithString("mod_id", mcp.Required(), mcp.Description("Mod ID to lock/unlock")),
			mcp.WithBoolean("locked", mcp.Required(), mcp.Description("True to lock (prevent updates), false to unlock")),
			mcp.WithString("project_path", mcp.Description("Optional path to modpack project. If omitted, uses active project")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			path := request.GetString("project_path", "")
			if err := a.ensureProject(path); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			modID, err := request.RequireString("mod_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			locked, err := request.RequireBool("locked")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			if err := a.ChangeModLocked(modID, locked); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to update lock status for %s: %v", modID, err)), nil
			}
			return mcp.NewToolResultText(fmt.Sprintf("Lock status for mod '%s' set to %t", modID, locked)), nil
		},
	)

	// Tool: change_mod_version
	s.AddTool(
		mcp.NewTool(
			"change_mod_version",
			mcp.WithDescription("Change the version of an installed mod (downloads new version and updates config)"),
			mcp.WithString("mod_id", mcp.Required(), mcp.Description("Mod ID to modify")),
			mcp.WithString("version", mcp.Required(), mcp.Description("Target version string to install")),
			mcp.WithString("project_path", mcp.Description("Optional path to modpack project. If omitted, uses active project")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			path := request.GetString("project_path", "")
			if err := a.ensureProject(path); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			modID, err := request.RequireString("mod_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			version, err := request.RequireString("version")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			if err := a.ChangeModVersion(modID, version); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to change version for %s: %v", modID, err)), nil
			}
			return mcp.NewToolResultText(fmt.Sprintf("Mod '%s' version successfully updated to '%s'", modID, version)), nil
		},
	)

	// Tool: check_mods_updates
	s.AddTool(
		mcp.NewTool(
			"check_mods_updates",
			mcp.WithDescription("Check for updates for mods in the modpack project"),
			mcp.WithArray("mod_ids", mcp.Description("Optional list of specific mod IDs to check. If omitted or empty, checks all unlocked mods"), mcp.WithStringItems()),
			mcp.WithString("project_path", mcp.Description("Optional path to modpack project. If omitted, uses active project")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			path := request.GetString("project_path", "")
			if err := a.ensureProject(path); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			modIDs := request.GetStringSlice("mod_ids", nil)
			if len(modIDs) == 0 {
				cfg, err := a.loadConfig()
				if err != nil {
					return mcp.NewToolResultError(fmt.Sprintf("Failed to load project config: %v", err)), nil
				}
				for id, m := range cfg.Mods {
					if !m.Locked && m.Source != "file" && m.Source != "" {
						modIDs = append(modIDs, id)
					}
				}
			}

			updates, err := a.CheckModsUpdates(modIDs)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to check mod updates: %v", err)), nil
			}
			data, _ := json.MarshalIndent(updates, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	// Tool: update_mods
	s.AddTool(
		mcp.NewTool(
			"update_mods",
			mcp.WithDescription("Update mods in the project. If specific mods are omitted, automatically checks and updates all eligible mods"),
			mcp.WithArray("mods", mcp.Description("Optional list of mods to update, each object with ModId, Version, URL")),
			mcp.WithString("project_path", mcp.Description("Optional path to modpack project. If omitted, uses active project")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			path := request.GetString("project_path", "")
			if err := a.ensureProject(path); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			type updatePayload struct {
				Mods []updater.ModToUpdate `json:"mods"`
			}
			var payload updatePayload
			_ = request.BindArguments(&payload)

			modsToUpdate := payload.Mods
			if len(modsToUpdate) == 0 {
				cfg, err := a.loadConfig()
				if err != nil {
					return mcp.NewToolResultError(fmt.Sprintf("Failed to load project config: %v", err)), nil
				}
				var modIDs []string
				for id, m := range cfg.Mods {
					if !m.Locked && m.Source != "file" && m.Source != "" {
						modIDs = append(modIDs, id)
					}
				}
				if len(modIDs) == 0 {
					return mcp.NewToolResultText("No mods eligible for update in this project"), nil
				}
				foundUpdates, err := a.CheckModsUpdates(modIDs)
				if err != nil {
					return mcp.NewToolResultError(fmt.Sprintf("Failed checking updates: %v", err)), nil
				}
				modsToUpdate = foundUpdates
			}

			if len(modsToUpdate) == 0 {
				return mcp.NewToolResultText("All mods are already up to date"), nil
			}

			if err := a.UpdateMods(modsToUpdate); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to update mods: %v", err)), nil
			}
			data, _ := json.MarshalIndent(modsToUpdate, "", "  ")
			return mcp.NewToolResultText(fmt.Sprintf("Updated %d mod(s) successfully:\n%s", len(modsToUpdate), string(data))), nil
		},
	)

	// Tool: install_mods
	s.AddTool(
		mcp.NewTool(
			"install_mods",
			mcp.WithDescription("Build and install the modpack: cleans and builds 'client' and 'server' folders and copies/downloads all configured mods based on their side configuration"),
			mcp.WithString("project_path", mcp.Description("Optional path to modpack project. If omitted, uses active project")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			path := request.GetString("project_path", "")
			if err := a.ensureProject(path); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			if err := a.InstallMods(); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to install mods: %v", err)), nil
			}
			return mcp.NewToolResultText("Modpack build completed: mods installed to 'client' and 'server' directories"), nil
		},
	)
}

func (a *App) registerResources(s *server.MCPServer) {
	// Resource: Current Project Configuration
	s.AddResource(
		mcp.NewResource(
			"packsmith://project/config",
			"Active Project Configuration",
			mcp.WithResourceDescription("The packsmith.json configuration for the currently open modpack project"),
			mcp.WithMIMEType("application/json"),
		),
		func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			if a.ProjectPath == "" {
				return nil, fmt.Errorf("no project is currently open")
			}
			cfg, err := a.loadConfig()
			if err != nil {
				return nil, err
			}
			data, err := json.MarshalIndent(cfg, "", "  ")
			if err != nil {
				return nil, err
			}
			return []mcp.ResourceContents{
				mcp.TextResourceContents{
					URI:      request.Params.URI,
					MIMEType: "application/json",
					Text:     string(data),
				},
			}, nil
		},
	)

	// Resource: Application Logs
	s.AddResource(
		mcp.NewResource(
			"packsmith://logs",
			"Application Logs",
			mcp.WithResourceDescription("Current application logs from app.log"),
			mcp.WithMIMEType("text/plain"),
		),
		func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			logs, err := a.GetLogs()
			if err != nil {
				return nil, err
			}
			return []mcp.ResourceContents{
				mcp.TextResourceContents{
					URI:      request.Params.URI,
					MIMEType: "text/plain",
					Text:     logs,
				},
			}, nil
		},
	)
}

func (a *App) registerPrompts(s *server.MCPServer) {
	// Prompt: modpack_summary
	s.AddPrompt(
		mcp.NewPrompt(
			"modpack_summary",
			mcp.WithPromptDescription("Summarize the current modpack project configuration, loader, Minecraft version, and mod list"),
		),
		func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
			var promptText string
			if a.ProjectPath == "" {
				promptText = "No modpack project is currently open in Packsmith. Please open a project first with open_project."
			} else {
				cfg, err := a.loadConfig()
				if err != nil {
					promptText = fmt.Sprintf("Error loading current project config: %v", err)
				} else {
					data, _ := json.MarshalIndent(cfg, "", "  ")
					promptText = fmt.Sprintf("Analyze and summarize this Packsmith modpack project:\n\n%s\n\nProvide: 1) Overview of Minecraft version and loader, 2) Breakdown of mods by side (client, server, both), 3) Any locked mods or potential improvements.", string(data))
				}
			}

			return mcp.NewGetPromptResult(
				"Summary of current modpack",
				[]mcp.PromptMessage{
					mcp.NewPromptMessage(mcp.RoleUser, mcp.NewTextContent(promptText)),
				},
			), nil
		},
	)

	// Prompt: diagnose_modpack
	s.AddPrompt(
		mcp.NewPrompt(
			"diagnose_modpack",
			mcp.WithPromptDescription("Check modpack configuration health, missing sides, and update status"),
		),
		func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
			promptText := "Please diagnose the active Packsmith modpack. Check if all mods have valid sides assigned, check for available updates using check_mods_updates, and suggest optimizations for client/server distribution."
			return mcp.NewGetPromptResult(
				"Modpack health diagnostic prompt",
				[]mcp.PromptMessage{
					mcp.NewPromptMessage(mcp.RoleUser, mcp.NewTextContent(promptText)),
				},
			), nil
		},
	)
}
