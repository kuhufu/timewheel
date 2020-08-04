package timewheel

type Option func(wheel *Wheel)

func WithOneGoRoutine() Option {
	return func(wheel *Wheel) {
		wheel.onlyOneGoRoutine = true
	}
}
