package context

import (
	"github.com/stretchr/testify/assert"
	"k8s.io/api/admission/v1beta1"
	"testing"
)

func TestHasChanged(t *testing.T) {
	ctx := createTestContext(`{"a": {"b": 1, "c": 2}, "d": 3}`, `{"a": {"b": 2, "c": 2}, "d": 4}`)

	val, err := ctx.HasChanged("a.b")
	assert.NoError(t, err)
	assert.True(t, val)

	val, err = ctx.HasChanged("a.c")
	assert.NoError(t, err)
	assert.False(t, val)

	val, err = ctx.HasChanged("d")
	assert.NoError(t, err)
	assert.True(t, val)

	val, err = ctx.HasChanged("a.x.y")
	assert.Error(t, err)
}

func TestRequestNotInitialize(t *testing.T) {
	request := &v1beta1.AdmissionRequest{}
	ctx := NewContext()
	ctx.AddRequest(request)

	_, err := ctx.HasChanged("x.y.z")
	assert.Error(t, err)
}

func TestMissingOldObject(t *testing.T) {
	request := &v1beta1.AdmissionRequest{}
	ctx := NewContext()
	ctx.AddRequest(request)
	request.Object.Raw = []byte(`{"a": {"b": 1, "c": 2}, "d": 3}`)

	_, err := ctx.HasChanged("a.b")
	assert.Error(t, err)
}

func TestMissingObject(t *testing.T) {
	request := &v1beta1.AdmissionRequest{}
	ctx := NewContext()
	ctx.AddRequest(request)
	request.OldObject.Raw = []byte(`{"a": {"b": 1, "c": 2}, "d": 3}`)

	_, err := ctx.HasChanged("a.b")
	assert.Error(t, err)
}

func createTestContext(obj, oldObj string) *Context {
	request := &v1beta1.AdmissionRequest{}
	request.Operation = "UPDATE"
	request.Object.Raw = []byte(obj)
	request.OldObject.Raw = []byte(oldObj)

	ctx := NewContext()
	ctx.AddRequest(request)
	return ctx
}
