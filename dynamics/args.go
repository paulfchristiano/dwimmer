package dynamics

import "github.com/paulfchristiano/dwimmer/term"

func Args0(f func(Dwimmer, *term.SettingT) term.T) func(Dwimmer, *term.SettingT, ...term.T) term.T {
	return func(d Dwimmer, s *term.SettingT, ts ...term.T) term.T {
		if len(ts) != 0 {
			panic("Wrong number of arguments")
		}
		return f(d, s)
	}
}

func Args1(f func(Dwimmer, *term.SettingT, term.T) term.T) func(Dwimmer, *term.SettingT, ...term.T) term.T {
	return func(d Dwimmer, s *term.SettingT, ts ...term.T) term.T {
		if len(ts) != 1 {
			panic("Wrong number of arguments")
		}
		return f(d, s, ts[0])
	}
}

func Args2(f func(Dwimmer, *term.SettingT, term.T, term.T) term.T) func(Dwimmer, *term.SettingT, ...term.T) term.T {
	return func(d Dwimmer, s *term.SettingT, ts ...term.T) term.T {
		if len(ts) != 2 {
			panic("Wrong number of arguments")
		}
		return f(d, s, ts[0], ts[1])
	}
}
func Args3(f func(Dwimmer, *term.SettingT, term.T, term.T, term.T) term.T) func(Dwimmer, *term.SettingT, ...term.T) term.T {
	return func(d Dwimmer, s *term.SettingT, ts ...term.T) term.T {
		if len(ts) != 3 {
			panic("Wrong number of arguments")
		}
		return f(d, s, ts[0], ts[1], ts[2])
	}
}
func Args4(f func(Dwimmer, *term.SettingT, term.T, term.T, term.T, term.T) term.T) func(Dwimmer, *term.SettingT, ...term.T) term.T {
	return func(d Dwimmer, s *term.SettingT, ts ...term.T) term.T {
		if len(ts) != 4 {
			panic("Wrong number of arguments")
		}
		return f(d, s, ts[0], ts[1], ts[2], ts[3])
	}
}
func Args5(f func(Dwimmer, *term.SettingT, term.T, term.T, term.T, term.T, term.T) term.T) func(Dwimmer, *term.SettingT, ...term.T) term.T {
	return func(d Dwimmer, s *term.SettingT, ts ...term.T) term.T {
		if len(ts) != 5 {
			panic("Wrong number of arguments")
		}
		return f(d, s, ts[0], ts[1], ts[2], ts[3], ts[4])
	}
}
