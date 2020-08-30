package hook

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/qiniu/qmgo/field"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type User struct {
	field.DefaultField `bson:",inline"`

	Name         string    `bson:"name"`
	Age          int       `bson:"age"`
	CreateTimeAt time.Time `bson:"createTimeAt"`
	UpdateTimeAt int64     `bson:"updateTimeAt"`
	// for test
	AfterCalled bool `bson:"afterCalled"`
}

func (u *User) BeforeInsert() error {
	if u.Name == "jz" || u.Name == "xm" {
		u.Age = 17
	}
	return nil
}

func (u *User) AfterInsert() error {
	u.AfterCalled = true
	return nil
}

func (u *User) CustomFields() field.CustomFieldsBuilder {
	return field.NewCustom().SetCreateAt("CreateTimeAt").SetUpdateAt("UpdateTimeAt")
}

func TestInsertHook(t *testing.T) {
	ast := require.New(t)

	u := &User{Name: "jz", Age: 7}
	err := Do(u, BeforeInsert)
	ast.NoError(err)
	ast.Equal(17, u.Age)
	// default fields
	ast.NotEqual(time.Time{}, u.CreateAt)
	ast.NotEqual(time.Time{}, u.UpdateAt)
	// custom fields
	ast.NotEqual(time.Time{}, u.CreateTimeAt)
	ast.NotEqual(time.Time{}, u.UpdateTimeAt)

	err = Do(u, AfterInsert)
	ast.NoError(err)
	ast.True(u.AfterCalled)

	u1, u2 := &User{Name: "jz", Age: 7}, &User{Name: "xm", Age: 8}
	us := []interface{}{u1, u2}
	err = Do(us, BeforeInsert)
	ast.NoError(err)
	err = Do(us, AfterInsert)
	ast.NoError(err)
	for _, v := range us {
		vv := v.(*User)
		if vv.Name == "jz" {
			ast.Equal(17, vv.Age)
		}
		if vv.Name == "xm" {
			ast.Equal(17, vv.Age)
		}
		ast.True(vv.AfterCalled)
	}

}

func (u *User) BeforeUpdate() error {
	if u.Name == "jz" || u.Name == "xm" {
		u.Age = 17
	}
	return nil
}

func (u *User) AfterUpdate() error {
	u.AfterCalled = true
	return nil
}
func TestUpdateHook(t *testing.T) {
	ast := require.New(t)

	u := &User{Name: "jz", Age: 7}
	err := Do(u, BeforeUpdate)
	ast.NoError(err)
	ast.Equal(17, u.Age)
	// default field
	ast.NotEqual(time.Time{}, u.UpdateAt)

	// custom fields
	ast.NotEqual(time.Time{}, u.UpdateTimeAt)

	err = Do(u, AfterUpdate)
	ast.NoError(err)
	ast.True(u.AfterCalled)

	u1, u2 := &User{Name: "jz", Age: 7}, &User{Name: "xm", Age: 8}
	us := []interface{}{u1, u2}
	err = Do(us, BeforeUpdate)
	ast.NoError(err)
	err = Do(us, AfterUpdate)
	ast.NoError(err)
	for _, v := range us {
		vv := v.(*User)
		if vv.Name == "jz" {
			ast.Equal(17, vv.Age)
		}
		if vv.Name == "xm" {
			ast.Equal(17, vv.Age)
		}
		ast.True(vv.AfterCalled)
	}

}

func (u *User) BeforeQuery() error {
	if u.Name == "jz" || u.Name == "xm" {
		u.Age = 17
	}
	fmt.Println("into before query")
	sliceBeforeQueryCount++
	return nil
}

func (u *User) AfterQuery() error {
	u.AfterCalled = true
	return nil
}
func TestQueryHook(t *testing.T) {
	ast := require.New(t)

	u := &User{Name: "jz", Age: 7}
	err := Do(u, BeforeQuery)
	ast.NoError(err)
	ast.Equal(17, u.Age)

	err = Do(u, AfterQuery)
	ast.NoError(err)
	ast.True(u.AfterCalled)

	u1, u2 := &User{Name: "jz", Age: 7}, &User{Name: "xm", Age: 8}
	us := []interface{}{u1, u2}
	err = Do(us, BeforeQuery)
	ast.NoError(err)
	err = Do(us, AfterQuery)
	ast.NoError(err)
	for _, v := range us {
		vv := v.(*User)
		if vv.Name == "jz" {
			ast.Equal(17, vv.Age)
		}
		if vv.Name == "xm" {
			ast.Equal(17, vv.Age)
		}
		ast.True(vv.AfterCalled)
	}

	uss := []*User{&User{Name: "jz"}, &User{Name: "xm"}}
	Do(&uss, BeforeQuery)

}
func (u *User) BeforeRemove() error {
	if u.Name == "jz" || u.Name == "xm" {
		u.Age = 17
	}
	return nil
}

func (u *User) AfterRemove() error {
	u.AfterCalled = true
	return nil
}

func TestRemoveHook(t *testing.T) {
	ast := require.New(t)

	u := &User{Name: "jz", Age: 7}
	err := Do(u, BeforeRemove)
	ast.NoError(err)
	ast.Equal(17, u.Age)

	err = Do(u, AfterRemove)
	ast.NoError(err)
	ast.True(u.AfterCalled)

	u1, u2 := &User{Name: "jz", Age: 7}, &User{Name: "xm", Age: 8}
	us := []interface{}{u1, u2}
	err = Do(us, BeforeRemove)
	ast.NoError(err)
	err = Do(us, AfterRemove)
	ast.NoError(err)
	for _, v := range us {
		vv := v.(*User)
		if vv.Name == "jz" {
			ast.Equal(17, vv.Age)
		}
		if vv.Name == "xm" {
			ast.Equal(17, vv.Age)
		}
		ast.True(vv.AfterCalled)
	}

}

type UserError struct {
	Name string
	Age  int

	// for test
	mock.Mock `bson:"-"`
}

func (u *UserError) BeforeInsert() error {
	return nil
}

func (u *UserError) AfterInsert() error {
	args := u.Called()
	return args.Error(0)
}
func TestSliceError(t *testing.T) {
	ast := require.New(t)

	u1, u2 := &UserError{Name: "jz", Age: 7}, &UserError{Name: "xm", Age: 8}
	us := []interface{}{u1, u2}

	u1.On("AfterInsert").Return(nil)
	u2.On("AfterInsert").Return(errors.New("called"))
	err := Do(us, AfterInsert)
	ast.Equal("called", err.Error())

}

type UserNoHook struct {
	Name string
	Age  int

	// for test
	AfterCalled bool
}

func TestUserNoHook(t *testing.T) {
	ast := require.New(t)

	u := &UserNoHook{Name: "jz", Age: 7}
	err := Do(u, BeforeInsert)
	ast.NoError(err)
	ast.Equal(7, u.Age)

	err = Do(u, AfterInsert)
	ast.NoError(err)

	u1, u2 := &UserNoHook{Name: "jz", Age: 7}, &UserNoHook{Name: "xm", Age: 8}
	us := []interface{}{u1, u2}
	err = Do(us, BeforeInsert)
	ast.NoError(err)
	err = Do(us, AfterInsert)
	ast.NoError(err)
	for _, v := range us {
		vv := v.(*UserNoHook)
		if vv.Name == "jz" {
			ast.Equal(7, vv.Age)
		}
		if vv.Name == "xm" {
			ast.Equal(8, vv.Age)
		}
		ast.False(vv.AfterCalled)
	}

	err = Do(u, BeforeUpdate)
	ast.NoError(err)
	err = Do(u, AfterUpdate)
	ast.NoError(err)
	err = Do(us, BeforeUpdate)
	ast.NoError(err)
	err = Do(us, AfterUpdate)
	ast.NoError(err)

	err = Do(u, BeforeQuery)
	ast.NoError(err)
	err = Do(u, AfterQuery)
	ast.NoError(err)
	err = Do(us, BeforeQuery)
	ast.NoError(err)
	err = Do(us, AfterQuery)
	ast.NoError(err)

	err = Do(u, BeforeRemove)
	ast.NoError(err)
	err = Do(u, AfterRemove)
	ast.NoError(err)
	err = Do(us, BeforeRemove)
	ast.NoError(err)
	err = Do(us, AfterRemove)
	ast.NoError(err)
}

var sliceBeforeQueryCount = 0

func TestSliceHook(t *testing.T) {
	sliceBeforeQueryCount = 0
	ast := require.New(t)

	u := &User{Name: "jz"}
	Do(u, BeforeQuery)

	uss := []*User{&User{Name: "jz"}, &User{Name: "xm"}}
	Do(uss, BeforeQuery)

	Do(&uss, BeforeQuery)

	ast.Equal(5, sliceBeforeQueryCount)

}
