package resolver_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rafaelsq/boiler/pkg/entity"
	gentity "github.com/rafaelsq/boiler/pkg/graphql/internal/entity"
	"github.com/rafaelsq/boiler/pkg/graphql/internal/resolver"
	"github.com/rafaelsq/boiler/pkg/iface"
	"github.com/rafaelsq/boiler/pkg/mock"
	"github.com/stretchr/testify/assert"
)

func TestUserUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// succeed
	{
		user := &entity.User{ID: 4, Name: "John Doe"}

		m := mock.NewMockService(ctrl)
		r := resolver.NewUser(m)

		m.EXPECT().
			GetUserByID(gomock.Any(), user.ID).
			Return(user, nil)

		u, err := r.User(context.TODO(), strconv.Itoa(user.ID))
		assert.Nil(t, err)
		assert.NotNil(t, u)
		assert.Equal(t, strconv.Itoa(user.ID), u.ID)
		assert.Equal(t, user.Name, u.Name)
	}

	// fail if invalid ID
	{
		m := mock.NewMockService(ctrl)
		r := resolver.NewUser(m)

		u, err := r.User(context.TODO(), "fail")
		assert.Nil(t, u)
		assert.Equal(t, err, iface.ErrInvalidID)
	}

	// fails if service fails
	{
		user := &entity.User{ID: 4, Name: "John Doe"}

		m := mock.NewMockService(ctrl)
		r := resolver.NewUser(m)

		m.EXPECT().
			GetUserByID(gomock.Any(), user.ID).
			Return(nil, fmt.Errorf("opz"))

		u, err := r.User(context.TODO(), strconv.Itoa(user.ID))
		assert.Nil(t, u)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "opz")
	}
}

func TestUserUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// succeed
	{
		user := &entity.User{ID: 4, Name: "John Doe"}

		m := mock.NewMockService(ctrl)
		r := resolver.NewUser(m)

		m.EXPECT().
			FilterUsers(gomock.Any(), iface.FilterUsers{Limit: 2}).
			Return([]*entity.User{user}, nil)

		users, err := r.Users(context.TODO(), 2)
		assert.Nil(t, err)
		assert.NotNil(t, users)
		assert.Equal(t, len(users), 1)
	}

	// fails if service fails
	{
		m := mock.NewMockService(ctrl)
		r := resolver.NewUser(m)

		m.EXPECT().
			FilterUsers(gomock.Any(), iface.FilterUsers{Limit: 4}).
			Return(nil, fmt.Errorf("opz"))

		users, err := r.Users(context.TODO(), 4)
		assert.Nil(t, users)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "opz")
	}
}

func TestUserEmails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// succeed
	{
		user := &entity.Email{ID: 4, Address: "a@b.c"}

		m := mock.NewMockService(ctrl)
		r := resolver.NewUser(m)

		m.EXPECT().
			FilterEmails(gomock.Any(), iface.FilterEmails{UserID: user.ID}).
			Return([]*entity.Email{user}, nil)

		emails, err := r.Emails(context.TODO(), &gentity.User{ID: strconv.Itoa(user.ID)})
		assert.Nil(t, err)
		assert.NotNil(t, emails)
		assert.Equal(t, len(emails), 1)
	}

	// fail if invalid ID
	{
		m := mock.NewMockService(ctrl)
		r := resolver.NewUser(m)

		emails, err := r.Emails(context.TODO(), &gentity.User{ID: "0"})
		assert.Nil(t, emails)
		assert.Equal(t, err, iface.ErrInvalidID)
	}

	// fails if service fails
	{
		m := mock.NewMockService(ctrl)
		r := resolver.NewUser(m)

		m.EXPECT().
			FilterEmails(gomock.Any(), iface.FilterEmails{UserID: 2}).
			Return(nil, fmt.Errorf("opz"))

		users, err := r.Emails(context.TODO(), &gentity.User{ID: strconv.Itoa(2)})
		assert.Nil(t, users)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "opz")
	}
}
