package cluster

// PromptAxisView is the current prompt and axis labeling shown to players.
type PromptAxisView struct {
	PromptText string
	XMinLabel  string
	XMaxLabel  string
	YMinLabel  string
	YMaxLabel  string
}

// DotView is one plotted player coordinate.
type DotView struct {
	Nickname        string
	X               float64
	Y               float64
	Points          int
	AnimationDelay  int
	IsOutlier       bool
	IsCurrentPlayer bool
}

// StandingView is one row in standings.
type StandingView struct {
	Nickname           string
	DistanceFromCenter float64
	HasCenterDistance  bool
	DistanceFromYou    float64
	HasDistanceFromYou bool
	IsCurrentPlayer    bool
	TotalPoints        int
	AvgPointsPerRound  float64
}
