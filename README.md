# vdi_poke

Sends a keypress (F15 by default, or Alt+Tab) on a timer to keep a Windows session from going idle. Single static `.exe`, no dependencies.

## Build

Requires Docker (no local Go install needed):

```sh
make        # -> bin/vdi_poke.exe
make vet
```

Or download `vdi_poke.exe` from the [latest release](https://github.com/PAllisonVSO/vdi_poke/releases/latest).

## Releasing

Push a `v*` tag and GitHub Actions builds `vdi_poke.exe` (with the tag as its version) plus a `.sha256` checksum and publishes them as a release:

```sh
git tag v0.0.2
git push origin v0.0.2
```

## Run

Copy `vdi_poke.exe` to the Windows machine and run it from a console:

```
vdi_poke.exe                       # F15 every 10 minutes
vdi_poke.exe -i 5m                 # custom interval (long form: -interval)
vdi_poke.exe -u 17:30              # stop automatically at 5:30 PM (long form: -until)
vdi_poke.exe -k alttab             # Alt+Tab instead of F15 (long form: -key)
vdi_poke.exe -v                    # print version (long form: -version)
vdi_poke.exe ?                     # help (also /?, -?, -h)
```

Press Ctrl+C to stop.

## Notes

- F15 is a no-op key on almost every system, so it doesn't disturb whatever you're doing.
- `-k alttab` really changes the focused window; keystrokes typed at that moment land in the other window.
- The interval must be shorter than the session's idle-lock policy.
- Hard session caps (e.g. Horizon "forcibly disconnect after N hours") are not affected.
- Environments enforcing AppLocker/WDAC may block an unsigned executable.
