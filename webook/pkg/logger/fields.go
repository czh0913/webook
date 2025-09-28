package logger

func String(Key, Val string) Field {
	return Field{
		Key:   Val,
		Value: Val,
	}
}
