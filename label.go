package basecfg

type Label string

func (l *Label) Key(k string) string {
	if l == nil || *l == "" {
		return k
	}

	return string(*l) + "." + k
}

func (l *Label) String() string {
	if l == nil {
		return ""
	}

	return string(*l)
}
