package Space

import "time"

func (space *Space) GetTime() time.Time {
	return space.time
}

func (space *Space) GetDeltaTime() time.Duration {
	return space.deltaTime
}

func (space *Space) setTime(t time.Time) {
	space.time = t
}

func (space *Space) setDeltaTime(t time.Duration) {
	space.deltaTime = t
}
