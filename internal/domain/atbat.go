package domain

type AtBat struct {
	GameID   string
	LeagueID string
	TeamID   string
	PlayerID string

	Seq     int
	Inning  int
	Half    string
	Strikes int
	Balls   int
	Fouls   int
	Result  string
	RBI     int
	Pitches int
}