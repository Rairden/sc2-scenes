package scene

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/Rairden/sc2-scenes/api/game"
	"github.com/Rairden/sc2-scenes/api/ui"
	"github.com/google/go-cmp/cmp"
)

type Menu int

type Scene struct {
	Menu Menu
	Prev Menu
	Game game.Game
	UI   ui.ActiveScreens
}

var (
	UIEndpoint   = "http://localhost:6119/ui"
	GameEndpoint = "http://localhost:6119/game"
	DisplayTime  = "http://localhost:6119/game/displayTime"
	firstLoad    = game.Game{Players: []any{}}
)

const (
	NIL Menu = iota
	HOME
	CAMPAIGN
	COOP
	VERSUS
	CUSTOM
	COLLECTION
	REPLAYS

	INGAME
	INREPLAY
	LOADING
	SCORESCREEN
)

func (m Menu) String() string {
	return [...]string{"closed", "home", "campaign", "coop", "versus", "custom",
		"collection", "replays", "in game", "in replay", "loading", "score screen"}[m]
}

// SetMenu sets the current state of the menu. http.Get() will return an error
// if you 1) exit the game or 2) load an old replay version.
func (s *Scene) SetMenu() error {
	ui1, err := GetUI(UIEndpoint)
	game1, err2 := GetGame(GameEndpoint)
	s.UI = ui1
	s.Game = game1
	s.Prev = s.Menu

	switch {
	case err != nil || err2 != nil || s.UI.ActiveScreens == nil:
		s.Menu = NIL
		return err
	case s.CheckMenu("ScreenLoading"):
		s.Menu = LOADING
	case s.CheckMenu("ScreenScore"):
		s.Menu = SCORESCREEN
	case s.IsFirstLoad():
		if len(s.UI.ActiveScreens) > 0 {
			s.Menu = HOME
		} else {
			s.Menu = s.Prev
		}
	case s.IsInGame():
		s.Menu = INGAME
	case s.IsReplay():
		s.Menu = INREPLAY
	default:
		s.Menu = HOME
	}

	return nil
}

func (s *Scene) IsInGame() bool {
	if s.Prev == LOADING || s.Prev == HOME {
		return len(s.UI.ActiveScreens) == 0 && len(s.Game.Players) > 0 &&
			!s.Game.IsReplay && s.Game.DisplayTime == 0
	}
	return len(s.UI.ActiveScreens) == 0 && !s.Game.IsReplay
}

func (s *Scene) IsReplay() bool {
	if s.Prev == LOADING || s.Prev == HOME {
		return len(s.UI.ActiveScreens) == 0 && len(s.Game.Players) > 0 && s.Game.DisplayTime == 0
	}
	return len(s.UI.ActiveScreens) == 0 && s.Game.IsReplay && len(s.Game.Players) > 0
}

// IsFirstLoad the /game response persists state from the last match or replay. So, the only time
// it is all zero values is when SC2 is first started.
func (s *Scene) IsFirstLoad() bool {
	return cmp.Equal(s.Game, firstLoad)
}

func (s *Scene) CheckMenu(menu string) bool {
	for _, screen := range s.UI.ActiveScreens {
		screen1 := screen.(string)
		scr := bytes.Index([]byte(screen1), []byte(menu))
		if scr == 0 {
			return true
		}
	}
	return false
}

func GetUI(url string) (ui.ActiveScreens, error) {
	resp, err := http.Get(url)
	if err != nil {
		return ui.ActiveScreens{}, err
	}

	body, err := io.ReadAll(resp.Body)

	var screens ui.ActiveScreens
	json.Unmarshal(body, &screens)

	return screens, nil
}

func GetGame(url string) (game.Game, error) {
	resp, err := http.Get(url)
	if err != nil {
		return game.Game{}, err
	}

	body, err := io.ReadAll(resp.Body)

	var game1 game.Game
	json.Unmarshal(body, &game1)

	return game1, nil
}

func PrettyPrint(resp any) string {
	body, _ := json.Marshal(resp)
	var prettyJSON bytes.Buffer
	json.Indent(&prettyJSON, body, "", "  ")
	return prettyJSON.String()
}
