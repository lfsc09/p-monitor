# p-monitor

A lightweight system performance monitor for Linux that runs in the system tray and displays real-time information about CPU, memory, disk, and GPU usage.

## Features

- **System Tray Integration**: Runs in the background with a system tray icon
- **Comprehensive Monitoring**: 
  - Disk usage (total capacity and used space percentage)
  - Memory usage (total capacity and usage percentage)
  - CPU usage percentage and temperature
  - GPU usage and temperature for all available GPUs (NVIDIA, AMD, integrated)
- **Reliable GPU Monitoring**: Uses command-line tools (`nvidia-smi`, `radeontop`) for accurate GPU metrics
- **Real-time CPU Temperature**: Reads actual CPU temperature from thermal sensors
- **Smooth Interface**: Non-flickering system tray menu with optimized updates
- **Configurable**: Adjustable update intervals (1-60 seconds/minutes) and temperature units (Celsius/Fahrenheit)
- **Error Handling**: Graceful handling of missing hardware or drivers
- **Logging**: Comprehensive logging to `~/.p-monitor/logs/`

## System Requirements

- Linux (tested on Ubuntu/Debian)
- Desktop environment with system tray support
- Go 1.19+ (for building from source)

### Optional Dependencies for GPU Monitoring

- `nvidia-smi` (for NVIDIA GPUs)
- `radeontop` (for AMD GPUs)

## Installation

### From GitHub Releases (Recommended)

1. Go to the [Releases page](https://github.com/lfsc09/p-monitor/releases) and download the latest `.deb` file
2. Install the package:
   ```bash
   sudo dpkg -i p-monitor_*_amd64.deb
   sudo apt-get install -f  # Fix any missing dependencies
   ```
3. The application will automatically start and appear in your system tray

### From Source

1. Clone the repository:
   ```bash
   git clone https://github.com/lfsc09/p-monitor.git
   cd p-monitor
   ```

2. Build the package:
   ```bash
   ./build-deb.sh
   ```

3. Install the package:
   ```bash
   sudo dpkg -i p-monitor_1.0.0_amd64.deb
   sudo apt-get install -f  # Fix any missing dependencies
   ```

4. The application will automatically start and appear in your system tray.

### Development Build

1. Install dependencies:
   ```bash
   sudo apt install libgl1-mesa-dev libx11-dev libxcursor-dev libxrandr-dev libxinerama-dev libxi-dev libxxf86vm-dev libglfw3-dev pkg-config
   ```

2. Build and run:
   ```bash
   go build -o p-monitor ./cmd
   ./p-monitor
   ```

## Usage

### System Tray Display

The system tray shows all metrics in the following format:

```
[disk-icon] HDD: 217.97GB (13.1%)
[memory-icon] RAM: 15.55GB (34.5%)
[cpu-icon] CPU: 12.2% (27.8°C)
[gpu-icon] NVIDIA GPU 0: 13.0% 36.0°C
[gpu-icon] AMD GPU 0: 5.2% 42.0°C
[gpu-icon] iGPU 0: 2.1% 35.0°C
```

### Configuration

Right-click the system tray icon to access configuration options:

- **Update Interval**: Change how often metrics are collected (1-60 seconds or minutes)
- **Temperature Unit**: Switch between Celsius and Fahrenheit

### Error Handling

If a component cannot be monitored (e.g., GPU drivers not installed), the system will:
- Display an error icon with "n/a" values
- Log the error to `~/.p-monitor/logs/`
- Continue monitoring other components

## Uninstallation

If the application doesn't work as expected or you want to remove it:

### Remove .deb Package

1. Stop the service:
   ```bash
   systemctl --user stop p-monitor.service
   systemctl --user disable p-monitor.service
   ```

2. Remove the package:
   ```bash
   sudo dpkg -r p-monitor
   ```

3. Clean up configuration and logs (optional):
   ```bash
   rm -rf ~/.p-monitor/
   ```

### Remove from Source Build

1. Stop the application if running:
   ```bash
   pkill p-monitor
   ```

2. Remove the binary:
   ```bash
   rm p-monitor
   ```

3. Clean up configuration and logs (optional):
   ```bash
   rm -rf ~/.p-monitor/
   ```

### Troubleshooting Uninstallation

If you encounter issues during uninstallation:

- **Service won't stop**: `systemctl --user kill p-monitor.service`
- **Package removal fails**: `sudo dpkg --purge p-monitor`
- **Files remain**: Manually remove `/usr/bin/p-monitor` and `/usr/share/p-monitor/`

## Configuration

Configuration is stored in `~/.p-monitor/config.json`:

```json
{
  "update_interval": 5,
  "time_unit": "seconds",
  "temperature_unit": "celsius"
}
```

## Logging

Logs are written to `~/.p-monitor/logs/` with timestamps. Each log file includes:
- System startup/shutdown events
- Configuration changes
- Error conditions
- Debug information

## Architecture

The application follows a clean architecture with separate packages:

- `pkg/monitor/`: System metrics collection
- `pkg/display/`: System tray interface and display logic
- `pkg/config/`: Configuration management
- `pkg/gpu/`: GPU-specific monitoring using command-line tools
- `pkg/types/`: Shared data structures
- `internal/logs/`: Logging system

## GPU Monitoring Details

### NVIDIA GPUs
- Uses `nvidia-smi` for reliable metrics
- Supports multiple NVIDIA GPUs
- Provides usage percentage and temperature

### AMD GPUs
- Uses `radeontop` for AMD GPU metrics
- Provides usage percentage
- Temperature monitoring depends on driver support

### Integrated GPUs
- Detects Intel/AMD integrated graphics
- Basic monitoring (usage monitoring not fully implemented)

## Troubleshooting

### Application doesn't appear in system tray
- Ensure you're running a desktop environment with system tray support
- Check that the application is running: `ps aux | grep p-monitor`
- Check logs: `tail -f ~/.p-monitor/logs/p-monitor_*.log`

### GPU monitoring not working
- Install appropriate drivers and tools:
  - NVIDIA: `sudo apt install nvidia-driver-xxx nvidia-smi`
  - AMD: `sudo apt install radeontop`
- Check if tools are available: `which nvidia-smi` or `which radeontop`

### Permission issues
- Ensure the user has access to system information
- Check that `/proc` and `/sys` filesystems are accessible

### Build errors
- **"cannot find -lXxf86vm"**: Install missing dependency: `sudo apt install libxxf86vm-dev`
- **"Package gl was not found"**: Install OpenGL development libraries: `sudo apt install libgl1-mesa-dev`
- **"X11/Xcursor/Xcursor.h: No such file or directory"**: Install X11 development libraries: `sudo apt install libxcursor-dev libx11-dev`
- **General build failures**: Run the build script which will check and report missing dependencies: `./build-deb.sh`

## Development

### Building from Source

```bash
go mod tidy
go build -o p-monitor ./cmd
```

### Testing Core Functionality

```bash
# Build and run the application to test
go build -o p-monitor ./cmd
./p-monitor
```

### Project Structure

```
p-monitor/
├── cmd/                 # Application entry points
├── pkg/                 # Main packages
│   ├── monitor/         # System monitoring
│   ├── display/         # System tray interface
│   ├── config/          # Configuration management
│   ├── gpu/            # GPU monitoring
│   └── types/          # Shared types
├── internal/           # Internal packages
│   └── logs/           # Logging system
├── assets/             # Icons and resources
├── debian/             # Debian package files
├── .github/workflows/  # GitHub Actions workflows
├── build-deb.sh       # Local build script
└── create-release.sh  # Release creation script
```

### GitHub Workflows

This project includes automated GitHub Actions workflows:

- **Release Build** (`.github/workflows/build-release.yml`): 
  - Triggers on new releases in the `main` branch
  - Automatically builds and uploads `.deb` packages
  - Creates SHA256 checksums for verification

- **Development Build** (`.github/workflows/build-dev.yml`):
  - Triggers on pushes to the `dev` branch
  - Builds development packages for testing
  - Uploads artifacts for download

### Creating Releases

To create a new release:

1. Make sure you're on the `main` branch
2. Run the release script:
   ```bash
   ./create-release.sh v1.0.0
   ```
3. The workflow will automatically build and attach the `.deb` package to the release
