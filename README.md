# Killaport

Tired of playing detective every time some random port is being a little b\*tch and refusing to die?

You know the drill:

```bash
netstat -ano | findstr :3000
tasklist | findstr "12345"
taskkill /PID 12345 /F
```

...then realize you typed the wrong PID and start over. Classic Monday vibes.

**Killaport** fixes that nonsense.

One binary. One port number. Zero brain cells required.

# How to use (literally 3 seconds)

1. Download the source code and build with `go build -o killaport.exe main.go`

2. Double-click built exe file (or run from terminal if you're feeling fancy)

3. Type the port you hate (e.g. `3000`, `8080`, `5432`) and press Enter

4. Watch it hunt and destroy. Done.

Example run:

```
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║                                                                                                  ║
║     ___  __    ___  ___       ___       ________  ________  ________  ________  _________        ║
║    |\  \|\  \ |\  \|\  \     |\  \     |\   __  \|\   __  \|\   __  \|\   __  \|\___   ___\      ║
║    \ \  \/  /|\ \  \ \  \    \ \  \    \ \  \|\  \ \  \|\  \ \  \|\  \ \  \|\  \|___ \  \_|      ║
║     \ \   ___  \ \  \ \  \    \ \  \    \ \   __  \ \   ____\ \  \\\  \ \   _  _\   \ \  \       ║
║      \ \  \\ \  \ \  \ \  \____\ \  \____\ \  \ \  \ \  \___|\ \  \\\  \ \  \\  \|   \ \  \      ║
║       \ \__\\ \__\ \__\ \_______\ \_______\ \__\ \__\ \__\    \ \_______\ \__\\ _\    \ \__\     ║
║        \|__| \|__|\|__|\|_______|\|_______|\|__|\|__|\|__|     \|_______|\|__|\|__|    \|__|     ║
║                                                                                                  ║
║                                                                                                  ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝

Enter port to kill: 3000
Found PID(s) on port 3000: [8764 8764]
Killed PID 8764
Press Enter to exit...
```

# Features because why not list them

- Kills all processes on the given port (TCP only for now)

- Shows you what it's about to murder before doing it

- Works on Windows (the .exe life)

- No admin rights needed in most cases

- Zero dependencies, single file, ~4-5MB

# Note for the Unix people staring angrily

This is Windows-first because `netstat` + `taskkill` pain is real on Windows.

If enough macOS/Linux users cry in the issues, maybe there will be a multi-platform version later. No promises.

Now go free those ports and stop swearing at your terminal.

---

### Happy killing,

**Veinz 🌿**
