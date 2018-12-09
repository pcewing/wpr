# Wallpaper Rotator

A simple wallpaper rotator that uses `feh`, written in go.

## Configuration

Create a config file at `$HOME/.config/wpr/wprrc.json` that looks like:
```json
{
  "WallpaperDir":"/home/username/Pictures/Wallpapers",
  "DisplayCount":2,
  "Interval":120
}
```

**Options**
* _WallpaperDir_: The directory containing wallpapers; the program will look through this recursively and currently does no filtering, so make sure all files in ll child directories are valid wallpapers.
* _DisplayCount_: The number of physical (Or virtual) displays that need to have wallpapers set.
* _Interval_: The number of seconds to wait between rotating the wallpapers

## Building

```bash
mkdir -p $GOPATH/src/github.com/pcewing/
git clone https://github.com/pcewing/wpr $GOPATH/src/github.com/pcewing/wpr
go install $GOPATH/src/github.com/pcewing/wpr
```

## Running
```bash
wpr
```
