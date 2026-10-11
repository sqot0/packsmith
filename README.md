<p align="center">
    <img src="https://i.imgur.com/cbSsyW9.png" alt="Packsmith - Smith modpack easy" />
</p>

Packsmith is an open source application for managing Minecraft modpacks that allows you:

- Manage modpacks with Forge, Fabric, Neoforge, or Quilt mod loaders.
- Search and download mods from Modrinth and CurseForge.
- Separate mods for client and server.
- Update, change versions, or delete mods from your modpack.
- Configure global application settings and adjust concurrent worker threads for file downloads and installs.
- Quick switch between projects via the Recent Projects history.
- Automate and interact with your modpacks via the built-in Model Context Protocol (MCP) server.

> [!WARNING]
> Please keep in mind that Packsmith is still under active development and may contain bugs or incomplete features

# Requirements
- Windows 10 or later
- Access to the internet and [Modrinth](https://modrinth.com)/[CurseForge](https://www.curseforge.com/) for downloading mods

# Quick Start

1. Download the latest release from [Releases](https://github.com/sqot0/packsmith/releases).
2. Run the executable:
   - packsmith.exe
3. Create a new project or open an existing one (or select from **Open > Open Recent**).
4. Add mods, manage versions, and configure client/server separation.
5. Customize global preferences in **Settings** (such as worker thread concurrency).

# Application Settings & Recent Projects

Packsmith persists global user configuration in standard system configuration directory:
- **Windows**: `%APPDATA%\packsmith\config.json`
- **Linux/macOS**: `~/.config/packsmith/config.json`

### Configuration Options
- `maxWorkers`: Maximum concurrent worker goroutines for mod downloads and file operations in the installer and updater (default: `15`).
- `recentProjects`: Automatically maintained history of recently opened or created modpacks.

### In the UI
- **Open Recent**: Accessible under the **Open** menu in the top menubar to reopen previous projects in a single click.
- **Settings Dialog**: Open via **Settings** in the top menubar to adjust concurrency or manage/clear project history.

# Building from Source

To build Packsmith from source, follow these steps:

1. Install Wails: Follow the installation guide at [https://wails.io/docs/gettingstarted/installation](https://wails.io/docs/gettingstarted/installation).
2. Install pnpm globally: Run `npm install -g pnpm`.
3. Build the application: Run `wails build`.
4. The executable will be located in `build/bin`.

# Model Context Protocol (MCP) Server

Packsmith provides a full Model Context Protocol (MCP) layer using [`github.com/mark3labs/mcp-go`](https://mcp-go.dev), allowing AI assistants (Claude Desktop, Cursor, Antigravity, etc.) to manage modpacks, search mods, inspect updates, configure settings, and build modpacks.

### Running the MCP Server

**Option 1: Using the main Packsmith executable**
```bash
packsmith.exe --mcp
# Or specify an initial project directory:
packsmith.exe --mcp "C:/path/to/modpack"
```

**Option 2: Using the standalone MCP CLI entrypoint**
```bash
go run ./cmd/mcp
# Or build a standalone binary:
go build -o packsmith-mcp.exe ./cmd/mcp
./packsmith-mcp.exe --project "C:/path/to/modpack"
```

### Configuration for MCP Clients

Add Packsmith to your MCP client configuration (e.g. `claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "packsmith": {
      "command": "C:/path/to/packsmith.exe",
      "args": ["--mcp"]
    }
  }
}
```

> [!TIP]
> You can also specify an initial project directory directly in `args`:
> ```json
> "args": ["--mcp", "C:/path/to/your/modpack"]
> ```
> If using the standalone binary `packsmith-mcp.exe`:
> ```json
> {
>   "mcpServers": {
>     "packsmith": {
>       "command": "C:/path/to/packsmith-mcp.exe",
>       "args": ["--project", "C:/path/to/your/modpack"]
>     }
>   }
> }
> ```

### Available MCP Tools

| Tool | Description |
|------|-------------|
| `open_project` | Open an existing modpack project directory |
| `get_project_config` | Retrieve the project configuration and installed mods |
| `initialize_project` | Initialize a new project with name, Minecraft version, and loader |
| `search_mods` | Search for compatible mods on Modrinth or CurseForge |
| `get_mod_versions` | List compatible versions for a specific mod |
| `add_mod` | Download and add a mod to the project |
| `add_local_mod` | Add a local `.jar` file from disk to the modpack |
| `remove_mod` | Remove a mod and delete its cached file |
| `change_mod_side` | Set execution side (`client`, `server`, `both`) |
| `change_mod_locked` | Lock or unlock a mod from automatic updates |
| `change_mod_version` | Switch a mod to a different version |
| `check_mods_updates` | Check for updates across installed mods |
| `update_mods` | Update specified or all eligible mods |
| `install_mods` | Build/install client and server modpack distribution directories |
| `get_app_settings` | Get global application settings (`maxWorkers`, recent projects) |
| `save_app_settings` | Update global application settings |
| `get_recent_projects` | List recently opened or initialized modpack projects |
| `clear_recent_projects` | Clear recent modpack projects history |
| `get_logs` | Retrieve recent application logs |

### Resources & Prompts

- **Resources**: `packsmith://project/config` (active `packsmith.json`), `packsmith://logs` (`app.log`)
- **Prompts**: `modpack_summary` (modpack audit and summary), `diagnose_modpack` (health check & update inspection)

# Contributing

Packsmith is open-source and welcomes contributions!
- Report bugs or request features on GitHub Issues.
- Contribute code via Pull Requests.
- Help improve documentation or add screenshots.
