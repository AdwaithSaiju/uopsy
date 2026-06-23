# uopsy

USB forensics tool. Scans connected and previously connected USB devices on Linux and Windows.

## Supported Platforms

| Platform | Current Devices | Historical Devices |
|---|---|---|
| Windows | `Get-PnpDevice` | Registry `HKLM\SYSTEM\CurrentControlSet\Enum\USB` |
| Linux | `/sys/bus/usb/devices/` sysfs | Kernel logs (`journalctl -k`) |

## Output

```
 [+] Scanning USB devices...
 [+] Found 3 devices (2 connected, 1 historical):

 [CONNECTED] USB-C Digital AV Multiport Adapter
   ├─ VID:    1A40
   ├─ PID:    0801
   └─ Serial: Unknown

 [CONNECTED] USB Keyboard
   ├─ VID:    046A
   ├─ PID:    0011
   └─ Serial: 12345678

 [DISCONNECTED] USB Flash Drive
   ├─ VID:    0781
   ├─ PID:    5583
   └─ Serial: 4C53000101022610
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
