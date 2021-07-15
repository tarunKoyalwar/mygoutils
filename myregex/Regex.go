package myregex

import "regexp"

type RegObject struct {
	Object *regexp.Regexp
}

func NewRegObject(exp string) *RegObject {
	return &RegObject{Object: regexp.MustCompile(exp)}
}

func (R RegObject) RegexReplaceBytes(old *[]byte, repl []byte) *[]byte {
	newbytes := R.Object.ReplaceAllLiteral(*old, repl)
	return &newbytes
}

func (R RegObject) RegexMatchBytes(data *[]byte) bool {
	return R.Object.Match(*data)
}
