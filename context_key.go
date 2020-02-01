package camforchat

// ContextKey constructs new key for values saved in context
type ContextKey string

func (c ContextKey) String() string {
	return "ctx" + string(c)
}
