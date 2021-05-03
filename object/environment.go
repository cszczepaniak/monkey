package object

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnvironment() *Environment {
	return &Environment{store: make(map[string]Object)}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func (e *Environment) Get(name string) (Object, bool) {
	v, ok := e.store[name]
	if !ok && e.outer != nil {
		v, ok = e.outer.Get(name)
	}
	return v, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
