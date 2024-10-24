# gOVERMON

System status overlay

## Motivation

Ever run out of disk space unexpectedly?
Want to keep an eye on memory usage while working with memory-leak-prone
software?
(Hi, Unreal. How it's going, Blender?)

## How to Use

### Supported Systems

Tested on: Windows

All underlying libs and most of the code should work on other platforms.
However, a small part of the code is not handling other platfroms at the
moment (e.g. hotkeys, extended memory graph)

### Shortcuts

- **Win + Shift + O**: Toggle passthrough mode
- **Win + Shift + I**: Exit

### Available Stats

- System memory usage/commit
- Free disk space

### Planned Stats

- Ping to a specified address(es)
- CPU usage
- Network I/O
- Disk I/O
- Temperature Sensors
- GPU Usage and Sensors

### That Won't Be Implemented

- FPS counter: There are enough convenient ways to track framerate, and too much
  risks and hustle implementing it.

## Development

### Why Ebitengine as GUI backend?

I tested most of the available Go GUI frameworks.
They all have similar (or worse) overhead compared to [Ebitengine](https://github.com/hajimehoshi/ebiten),
and Ebitengine is the only one with an up-to-date backend with cross-platform
mouse passthrough, transparency and active support.

### How to contribute

1. Fork
2. Commit
3. Open pull request
