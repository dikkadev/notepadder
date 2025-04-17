# notepadder

CLI tool to open or activate Notepad with a new tab on Windows.

Usage:
  notepadder [--no-new] 

## Building

To compile the GUI-only executable without a console window:

```
go build -ldflags "-H=windowsgui" -o notepadder.exe cmd/notepadder/main.go
```

## Testing

Run unit tests for the Windows wrapper:

```
go test ./pkg/win
``` 