# Backend README

## Project Description

This is the backend component of Packsmith, a modpack management tool for Minecraft. It provides a set of APIs and utilities to manage Minecraft modpacks, including searching for mods, downloading, installing, updating, and importing mods from various platforms like CurseForge and Modrinth.

The backend is written in Go and uses Wails for integration with the frontend.

## Architecture

The backend is organized into the following directories:

- `cmd/`: Contains the main application logic and Wails bindings.
- `internal/`: Internal packages for various functionalities.
  - `config/`: Configuration management.
  - `fs/`: File system operations (download, copy, delete).
  - `importer/`: Mod importing from directories.
  - `installer/`: Mod installation to client/server folders.
  - `logger/`: Logging utilities.
  - `sources/`: Integration with mod sources (CurseForge, Modrinth).
  - `updater/`: Mod update checking and applying.
  - `util/`: Utility functions, including worker pools for concurrency.

## Best Practices

### Code Quality
- Use meaningful variable and function names.
- Keep functions small and focused on a single responsibility (DRY, KISS principles).

### Logging
- All functions should log their entry and key operations using `logger.Log.Printf` or `logger.Log.Println`.
- Log errors with context for easier debugging.
- Use structured logging where possible.

### Concurrency
- Use goroutines and channels for concurrent operations.
- Utilize the `WorkerPool` and `WorkerPoolWithError` functions for parallel processing.
- Ensure thread safety with mutexes when accessing shared data.

### Error Handling
- Always check for errors and handle them appropriately.
- Propagate errors up the call stack with context.
- Log errors but do not panic unless absolutely necessary.

### Configuration
- Use the `config` package for all configuration management.
- Validate configuration data before use.

### File System Operations
- Use the `fs` package for all file operations to ensure consistency.