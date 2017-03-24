package object

func NewChildEnv(parent *Env) *Env {
	env := NewEnv()
	env.parent = parent
	return env
}

func NewEnv() *Env {
	s := make(map[string]Object)
	return &Env{state: s, parent: nil}
}

type Env struct {
	state  map[string]Object
	parent *Env
}

func (e *Env) Get(name string) (Object, bool) {
	obj, ok := e.state[name]
	if !ok && e.parent != nil {
		obj, ok = e.parent.Get(name)
	}
	return obj, ok
}

func (e *Env) Set(name string, obj Object) Object {
	e.state[name] = obj
	return obj
}
