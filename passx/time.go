package passx

type Time struct {
	// M, T, W, R, or F
	day rune

	// Military-time-like integer (see comment in Times).
	time int
}
