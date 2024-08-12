# file-organizer

<img src="./.github/image.png" width="600px">

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
Usage of ./file-organizer:
  -c, --config string      The path to the configuration file (default "~/.config/file-organizer/config.json")
  -d, --directory string   Path to organize files
      --prepend-date       Prepend the current date to the file name

Examples:
  ./file-organizer -d ~/Downloads
  ./file-organizer -d ~/Downloads --prepend-date
```

## License

MIT License