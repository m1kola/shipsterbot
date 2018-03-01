package gomockhelpers

// MatcherFunc allows to define a closure function
// that matches an argument with captured values
type MatcherFunc func(x interface{}) bool

func (MatcherFunc) String() string {
	return "matches values captured in a func"
}

// Matches calls a closure function
func (f MatcherFunc) Matches(x interface{}) bool {
	return f(x)
}
