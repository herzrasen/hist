package handler

import (
	"errors"
	"github.com/herzrasen/hist/args"
	"github.com/herzrasen/hist/config"
	"github.com/herzrasen/hist/handler/mocks"
	"github.com/herzrasen/hist/record"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHandler_Handle(t *testing.T) {
	t.Run("handle list", func(t *testing.T) {
		m := mocks.NewHistClient(t)
		m.On("List", mock.Anything).Return([]record.Record{}, nil)
		h := Handler{Client: m}
		err := h.Handle(args.Args{
			List: &args.ListCmd{},
		})
		require.NoError(t, err)
	})

	t.Run("list error", func(t *testing.T) {
		m := mocks.NewHistClient(t)
		m.On("List", mock.Anything).Return([]record.Record{}, errors.New("some error"))
		h := Handler{Client: m}
		err := h.Handle(args.Args{
			List: &args.ListCmd{},
		})
		require.Error(t, err)
	})

	t.Run("handle record", func(t *testing.T) {
		m := mocks.NewHistClient(t)
		m.On("Update", mock.Anything).Return(nil)
		h := Handler{
			Client: m,
			Config: &config.Config{},
		}
		err := h.Handle(args.Args{
			Record: &args.RecordCmd{Command: "ls -alF"},
		})
		require.NoError(t, err)
	})

	t.Run("record error", func(t *testing.T) {
		m := mocks.NewHistClient(t)
		m.On("Update", mock.Anything).Return(errors.New("some error"))
		h := Handler{
			Client: m,
			Config: &config.Config{},
		}
		err := h.Handle(args.Args{
			Record: &args.RecordCmd{Command: "ls -alF"},
		})
		require.Error(t, err)
	})

	t.Run("handle delete", func(t *testing.T) {
		m := mocks.NewHistClient(t)
		m.On("Delete", mock.Anything).Return(nil)
		h := Handler{
			Client: m,
			Config: &config.Config{},
		}
		err := h.Handle(args.Args{
			Delete: &args.DeleteCmd{Pattern: "foo"},
		})
		require.NoError(t, err)
	})

	t.Run("delete error", func(t *testing.T) {
		m := mocks.NewHistClient(t)
		m.On("Delete", mock.Anything).Return(errors.New("some error"))
		h := Handler{
			Client: m,
			Config: &config.Config{},
		}
		err := h.Handle(args.Args{
			Delete: &args.DeleteCmd{Pattern: "foo"},
		})
		require.Error(t, err)
	})

	t.Run("handle get", func(t *testing.T) {
		m := mocks.NewHistClient(t)
		m.On("Get", mock.Anything).Return("some-command", nil)
		h := Handler{
			Client: m,
			Config: &config.Config{},
		}
		err := h.Handle(args.Args{
			Get: &args.GetCmd{Index: 101},
		})
		require.NoError(t, err)
	})

	t.Run("get error", func(t *testing.T) {
		m := mocks.NewHistClient(t)
		m.On("Get", mock.Anything).Return("", errors.New("some error"))
		h := Handler{
			Client: m,
			Config: &config.Config{},
		}
		err := h.Handle(args.Args{
			Get: &args.GetCmd{Index: 101},
		})
		require.Error(t, err)
	})

}
