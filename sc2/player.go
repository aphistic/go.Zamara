package sc2

const (
	PlayerUnknown = iota
	PlayerHuman
	PlayerComputer
)

const (
	RaceUnknown = iota
	RaceRandom
	RaceTerran
	RaceProtoss
	RaceZerg
)

const (
	DifficultyUnknown = iota
	DifficultyVeryEasy
	DifficultyEasy
	DifficultyMedium
	DifficultyHard
	DifficultyVeryHard
	DifficultyInsane
)

const (
	ColorUnknown = iota
	ColorRed
	ColorBlue
	ColorTeal
	ColorPurple
	ColorYellow
	ColorOrange
	ColorGreen
	ColorLightPink
	ColorViolet
	ColorLightGrey
	ColorDarkGreen
	ColorBrown
	ColorLightGreen
	ColorDarkGrey
	ColorPink
)

type Player struct {
	Name string
	Id   int64

	Type int

	Team       int
	Color      Color
	NamedColor int

	ChosenRace int
	ActualRace int
	Difficulty int
	Handicap   int

	Outcome int
}

func newPlayer(value *serializedValue) (player *Player, err error) {
	player = new(Player)
	player.load(value)

	return
}

func (player *Player) load(value *serializedValue) (err error) {
	raceMap := map[string]int{
		"Protoss": RaceProtoss,
		"Terran":  RaceTerran,
		"Zerg":    RaceZerg,
	}

	player.Name = value.i(0).asString()
	player.Id = value.i(1).i(3).asInt64()
	player.Color.A = int(value.i(3).i(0).asInt64())
	player.Color.R = int(value.i(3).i(1).asInt64())
	player.Color.G = int(value.i(3).i(2).asInt64())
	player.Color.B = int(value.i(3).i(3).asInt64())
	player.Team = int(value.i(5).asInt64())
	player.Handicap = int(value.i(6).asInt64())
	player.Outcome = int(value.i(8).asInt64())

	val, exists := raceMap[value.i(2).asString()]
	if !exists {
		val = RaceUnknown
	}
	player.ActualRace = val

	return
}
