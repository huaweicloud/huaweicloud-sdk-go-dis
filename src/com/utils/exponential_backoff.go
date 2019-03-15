package utils

import (
	"time"
)

var STOP int64 = -1

//The default initial interval.
var DEFAULT_INITIAL_INTERVAL int64 = 50

//The default multiplier (increases the interval by 50%).
var DEFAULT_MULTIPLIER float64 = 1.5

//The default maximum back off time.
var DEFAULT_MAX_INTERVAL int64 = 30000

//The default maximum elapsed time
var DEFAULT_MAX_ELAPSED_TIME int64 = 0x7fffffffffffffff

type ExponentialBackOff struct {
	InitialInterval    int64
	Multiplier         float64
	MaxInterval        int64
	MaxElapsedTime     int64
	CurrentElapsedTime int64
	CurrentInterval    int64
}

func NewExponentialBackOff(initialInterval int64, multiplier float64, maxInterval, maxElapsedTime int64) *ExponentialBackOff {
	if multiplier < 1 {
		return nil
	}

	return &ExponentialBackOff{InitialInterval: initialInterval, Multiplier: multiplier, MaxInterval: maxInterval, MaxElapsedTime: maxElapsedTime, CurrentInterval: -1, CurrentElapsedTime: 0}
}

func (e *ExponentialBackOff) BackOff(backOff int64) int64 {
	if backOff != STOP {
		time.Sleep(time.Duration(backOff) * time.Millisecond)
	}
	return backOff
}

func (e *ExponentialBackOff) GetNextBackOff() int64 {
	if e.CurrentElapsedTime >= e.MaxElapsedTime {
		return STOP
	}

	nextInterval := e.computeNextInterval()
	e.CurrentElapsedTime += nextInterval
	return nextInterval
}

func (e *ExponentialBackOff) computeNextInterval() int64 {
	maxInterval := e.MaxInterval

	if e.CurrentInterval >= maxInterval {
		return maxInterval
	} else if e.CurrentInterval < 0 {
		initialInterval := e.InitialInterval
		if initialInterval < maxInterval {
			e.CurrentInterval = initialInterval
		} else {
			e.CurrentInterval = maxInterval
		}
	} else {
		e.CurrentInterval = e.multiplyInterval(maxInterval)
	}
	return e.CurrentInterval
}

func (e *ExponentialBackOff) multiplyInterval(maxInterval int64) int64 {
	i := int64(float64(e.CurrentInterval) * e.Multiplier)
	if i > maxInterval {
		return maxInterval
	}
	return i
}

func (e *ExponentialBackOff) ResetCurrentInterval() {
	e.CurrentInterval = -1
}
