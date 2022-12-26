# Quick start

Open SC2 and run:

```sh
go run ./cmd/scenes.go
```

# Usage

```go
import "sc2-scenes/scene"

func main() {
	scn := scene.Scene{}
	scn.SetMenu()
	fmt.Println(scene.PrettyPrint(scn.UI.ActiveScreens))
	fmt.Println(scene.PrettyPrint(scn.Game))
}
```

## Output

```json
[
  "ScreenBackgroundSC2/ScreenBackgroundSC2",
  "ScreenReplay/ScreenReplay",
  "ScreenNavigationSC2/ScreenNavigationSC2",
  "ScreenForegroundSC2/ScreenForegroundSC2"
]
{
  "isReplay": false,
  "displayTime": 0,
  "players": []
}
```

# Build

```sh
go get -u github.com/google/go-cmp/cmp
go build ./cmd/scenes.go
```
