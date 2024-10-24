#

## REPO

- [ ] Cleanup golangci
- [ ] Add CI for tests
- [ ] Add CI for build
- [ ] Cleanup TODOs from code

## v0.2

- [ ] STATS: FIX: unused series not discarded
- [ ] STATS: DEBUG: active series count Graph
- [ ] APP: Config file
  - [ ] Colors
- [ ] PLOT: Thresholds styling minimal implementation
- [ ] SHORTCUT: Hide to Tray?
- [ ] APP: Switch to passthrough after timeout
- [ ] APP: Categories for Graphs
- [ ] APP: Switch for whole Graph Categories
- [ ] APP: Switch for Graph Set profiles(multiple sets)
- [ ] SYSTRAY: Reset position(fix window outside of boundaries)
- [ ] APP: Unify logging implementation

## v0.3

- [ ] APP: Compact mode: show only graphs with triggered thresholds
- [ ] APP: Compact mode: show only small icons that thresholds currently triggered
- [ ] LAYOUT: Horizontal layout
- [ ] LAYOUT: Columns

## STATS

- [ ] STATS: STORAGE I/O stats
- [ ] STATS: NET stats
- [ ] STATS: GPU stats
- [ ] STATS: MEM stats
- [ ] STATS: SENSORS
  - [ ] Sensors require elevation and gopsutil doesn't show anything interesting on Win
  - [ ] Try prometheus windows_exporter stats
    - [ ] Check what sensors are available
    - [ ] Check if we can use exporter inside our code without running separate process
- [ ] STATS: Ping to address(i.e. RETN)
- [ ] STATS: Track cpu/mem for focused app by shortcut

## PLOT

- [ ] PLOT: Grid
- [ ] PLOT: Ticks visualization
- [ ] PLOT: Thresholds styling
- [ ] PLOT: Simple Line without filling
- [ ] PLOT: line + different color for filling(gradient?)
