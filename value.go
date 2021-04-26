package gcache

import "time"

type cacheValue interface {
	GetKey() interface{}
	GetValue() interface{}
	SetKey(key interface{})
	SetValue(val interface{})
	IsDeleted() bool
	SetDeadline(deadline time.Time)
}

func newValue(key, val interface{}, expiry time.Duration) cacheValue {
	if expiry > 0 {
		return &deadlineValue{
			baseValue: baseValue{
				key:   key,
				value: val,
			},
			deadline: time.Now().Add(expiry),
		}
	}
	return &baseValue{
		key:   key,
		value: val,
	}
}

type baseValue struct {
	key   interface{}
	value interface{}
}

func (v *baseValue) GetKey() interface{} {
	return v.key
}
func (v *baseValue) GetValue() interface{} {
	return v.value
}
func (v *baseValue) SetKey(key interface{}) {
	v.key = key
}
func (v *baseValue) SetValue(val interface{}) {
	v.value = val
}
func (v *baseValue) IsDeleted() bool {
	return false
}
func (v *baseValue) SetDeadline(deadline time.Time) {
	panic(`baseValue not support SetDeadline`)
}

type deadlineValue struct {
	baseValue
	deadline time.Time
}

func (v *deadlineValue) IsDeleted() bool {
	return !v.deadline.After(time.Now())
}
func (v *deadlineValue) SetDeadline(deadline time.Time) {
	v.deadline = deadline
}
