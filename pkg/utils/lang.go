package utils

func T[T any](cond bool, onTrue, onFalse T) T {
	if cond {
		return onTrue
	}

	return onFalse
}

func GetOptionalParam[T any](param ...T) *T {
	if len(param) > 0 {
		return &param[0]
	}

	return nil
}
