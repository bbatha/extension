package extension

import (
	"testing"
)

type Base interface {
	Extension
	Base() int32
}

type Extended interface {
	Base
	Extended() int32
}

type OnlyBase struct{}

var _ Base = (*OnlyBase)(nil)

func (ob OnlyBase) Base() int32 {
	return 1
}

func (ob OnlyBase) As(interface{}) bool {
	return false
}

type ExtendedBase struct{}

var _ Base = (*ExtendedBase)(nil)

func (eb ExtendedBase) Base() int32 {
	return 1
}

func (eb ExtendedBase) Extended() int32 {
	return 1
}

func (eb ExtendedBase) As(interface{}) bool {
	return false
}

type Wrapper struct {
	base Base
}

var _ Base = (*Wrapper)(nil)

func (w Wrapper) As(target interface{}) bool {
	if CastTo(w, target) {
		return true
	}
	return As(w.base, target)
}

func (w Wrapper) Base() int32 {
	b := w.base.Base()

	var e Extended
	if As(w, &e) {
		return b + e.Extended()
	}
	return b
}

type Override struct {
	base Base
}

var _ Base = (*Override)(nil)

func (o Override) As(target interface{}) bool {
	if CastTo(o, target) {
		return true
	}
	return As(o.base, target)
}

func (o Override) Base() int32 {
	return o.base.Base()
}

func (o Override) Extended() int32 {
	return 2
}

func TestExtension(t *testing.T) {
	ob := OnlyBase{}
	var e Extended
	if As(ob, &e) {
		t.Error("OB should not implement extended")
	}

	eb := ExtendedBase{}
	if !As(eb, &e) {
		t.Error("EB should implement extended")
	}
	if e.Extended() != 1 {
		t.Error(("EB.Extended should return 1"))
	}

	wob := Wrapper{ob}
	if As(wob, &e) {
		t.Error("Wrapped OnlyBase should not implement ExtendedBase")
	}
	if wob.Base() != 1 {
		t.Error("Wrapped OnlyBase should return 1")
	}

	web := Wrapper{eb}
	if !As(web, &e) {
		t.Error("Wrapped ExtendedBase should implement ExtendedBase")
	}
	if web.Base() != 2 {
		t.Error("Wrapped ExtendedBase should return 2")
	}
	if e.Extended() != 1 {
		t.Error("web.Extended should return 1")
	}

	oob := Override{ob}
	if !As(oob, &e) {
		t.Error("Override OnlyBase should implement ExtendedBase")
	}
	if oob.Base() != 1 {
		t.Error("Override OnlyBase should return 1")
	}
	if e.Extended() != 2 {
		t.Error("oeb.Extended should return 2")
	}

	oeb := Override{eb}
	if !As(oeb, &e) {
		t.Error("Override ExtendedBase should implement ExtendedBase")
	}
	if oeb.Base() != 1 {
		t.Error("Override ExtendedBase should return 1")
	}
	if e.Extended() != 2 {
		t.Error("oeb.Extended should return 2")
	}
}
