# notepadder

CLI tool to open or activate Notepad with a new tab on Windows.

Usage:
  notepadder [--no-new] [--debug]

Flags:
  --no-new, -n   Do not open a new tab.
  --debug        Print debug output (requires console build).

## Building

To compile the standard GUI-only executable (no console window):

```
go build -ldflags "-H=windowsgui" -o notepadder.exe cmd/notepadder/main.go
```

To compile for debugging (will show console window):

```
go build -o notepadder_debug.exe cmd/notepadder/main.go
```
Then run with `./notepadder_debug.exe --debug`

## Testing

Run unit tests for the Windows wrapper:

```
go test ./pkg/win
``` 