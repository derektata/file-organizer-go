# file-organizer

Organize your downloaded files into categorical folders as they are downloaded.

This package was made to be used on macOS machines using Automator and Folder Actions.

## Installation/Setup

```
git clone https://github.com/derektata/file-organizer-go.git
cd file-organizer-go
make
```

## Usage

```
Usage of file-organizer:
  -p, --path string    Path to organize files
  -d, --prepend-date   Prepend the current date to the file name

Examples:
  ./file-organizer -p ~/Downloads
  ./file-organizer -p ~/Downloads -d
```

## License

MIT License