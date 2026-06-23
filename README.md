# uopsy

USB forensics and analysis tool. Scans connected USB devices and reports vendor ID, product ID, serial number, and connection status.

## Supported Platforms

| Platform | Method |
|---|---|
| macOS | `system_profiler SPUSBDataType -json` |
| Linux | `/sys/bus/usb/devices/` sysfs |
| Windows | `Get-PnpDevice` via PowerShell |

## Output

```
 [+] Scanning USB devices...
 [CONNECTED] USB-C Digital AV Multiport Adapter
   ├─ VID:    1A40
   ├─ PID:    0801
   └─ Serial: Unknown

 [CONNECTED] USB Keyboard
   ├─ VID:    046A
   ├─ PID:    0011
   └─ Serial: 12345678
```

## Usage

```
uopsy scan
```

## Install

```
go install github.com/AdwaithSaiju/uopsy@latest
```

## Build

```
git clone https://github.com/AdwaithSaiju/uopsy.git
cd uopsy
go build -o uopsy .
```
