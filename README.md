# file-organizer

<img src="./.github/image.png" width="600px">

File Organizer is a tool designed to help you keep your downloaded files in order. It automatically moves files into categorized folders based on their file extensions or MIME types, making it easy to find what you need when you need it. This tool is particularly useful for macOS users and is designed to work with Automator and Folder Actions.

## Features

- **Automatic File Organization:** Automatically categorizes and moves files into subdirectories based on their file extension or MIME type.
- **Date Prepending:** Optionally prepend the current date to filenames for better sorting.
- **Dry Run Mode:** Preview changes without moving any files.
- **Customizable:** Easily customize the categorization rules via a JSON configuration file.

## Installation/Setup

To install and set up the file organizer:

```bash
git clone https://github.com/derektata/file-organizer-go.git
cd file-organizer-go
make
```

## Usage

```
Usage of file-organizer:
  -c, --config string      The path to the configuration file (default "~/.config/file-organizer/config.json")
  -d, --directory string   The path to the directory to organize
      --dry-run            Print the actions without moving any files
      --prepend-date       Prepend the current date to the filenames when moving files
      
Examples:
  ./file-organizer -d ~/Downloads
  ./file-organizer -d ~/Downloads --prepend-date
  ./file-organizer -d ~/Downloads --dry-run
```

## License

MIT License