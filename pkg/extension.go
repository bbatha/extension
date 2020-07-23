// Package extension provides some building blocks for the extension interface
// pattern.
//
// To .Net developers this package essentially provides IServiceProvider.
package extension

import (
	"reflect"
)

// Extension should be embedded in another interface that you wish to extend later. Interfaces
// that embed the base interface (the first one to implement Extension) can transparently
// delegate to specialized implemenations of interfaces without special casting
type Extension interface {
	As(target interface{}) bool
}

// CastTo returns true if self can be cast to the interface pointed to by target, and mutates target
// to set it to self.
func CastTo(self interface{}, target interface{}) bool {
	targetType := reflect.TypeOf(target)
	if targetType.Kind() != reflect.Ptr {
		return false
	}
	targetInterface := targetType.Elem()
	if targetInterface.Kind() != reflect.Interface {
		return false
	}
	if !reflect.TypeOf(self).Implements(targetInterface) {
		return false
	}
	reflect.ValueOf(target).Elem().Set(reflect.ValueOf(self))
	return true
}

// CanCast if self implements the interface target. If target is not an a *pointer to* an
// interface false is returned.
func CanCast(self interface{}, target interface{}) bool {
	targetType := reflect.TypeOf(target)
	if targetType.Kind() != reflect.Ptr {
		return false
	}
	targetInterface := targetType.Elem()
	if targetInterface.Kind() != reflect.Interface {
		return false
	}
	if !reflect.TypeOf(self).Implements(targetInterface) {
		return false
	}
	return true
}

// As recursively delegates through self until a type that implements the interface *pointed to* by
// target. If the current self implements the interface it is returned otherwise it delegates to the
// self's As implementation which may choose to recurse through potential delegates.
func As(self Extension, target interface{}) bool {
	if CastTo(self, target) {
		return true
	}
	return self.As(target)
}
