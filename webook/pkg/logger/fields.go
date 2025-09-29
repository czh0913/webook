package logger

func String(Key, Val string) Field {
	return Field{
		Key:   Val,
		Value: Val,
	}
}

func Error(err error) Field {
	return Field{
		Key:   "error",
		Value: err,
	}
}
