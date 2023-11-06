package internal

type Kill struct {
	Tick int
}

type Player struct {
	UserId int
	Kills  []Kill
}
