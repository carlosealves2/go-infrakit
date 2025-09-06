package attribute

type KeyValue struct{}

func String(k, v string) KeyValue { return KeyValue{} }
func Int(k string, v int) KeyValue { return KeyValue{} }
func Bool(k string, v bool) KeyValue { return KeyValue{} }
