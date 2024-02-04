# File Synchronization Tool

File Synchronization Tool is a command-line application designed to synchronize files and directories from a specified source to a destination. It also provides real-time monitoring of a specified source folder for file operations such as creating, updating, and deleting files.

## Features

- **Synchronization:** Copies files and directories from a source location to a specified destination, ensuring both locations stay up-to-date.
- **Real-time Monitoring:** Watches the specified source folder for file operations, allowing immediate synchronization of changes.
- **Create, Update, Delete Operations:** Detects and handles file creation, modification, and deletion efficiently.

## Installation

To use the File Synchronization Tool, follow these steps:

1. **Clone the Repository:**
   ```
   git clone https://github.com/prajwalscodestack/FileSyncTool.git
   ```

2. **Build the Application:**
   ```
   cd filesynctool
   go build
   ```

## Usage

1. **Synchronize Files and Directories:**
   ```
   ./filesynctool -source <source_path> -destination <destination_path>
   ```
   Replace `<source_path>` with the path to the source directory and `<destination_path>` with the path to the destination directory.
   

## Contributing

Contributions are welcome! If you have any ideas for improvements, bug fixes, or new features, feel free to open an issue or submit a pull request.

## License

This project is licensed under the [MIT License](LICENSE).
