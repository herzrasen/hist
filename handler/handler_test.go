package handler

import (
	"errors"
	"github.com/herzrasen/hist/args"
	"github.com/herzrasen/hist/handler/mocks"
	"github.com/herzrasen/hist/record"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"os"
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
		m.On("Record", mock.Anything).Return(nil)
		h := Handler{
			Client: m,
		}
		err := h.Handle(args.Args{
			Record: &args.RecordCmd{Command: "ls -alF"},
		})
		require.NoError(t, err)
	})

	t.Run("record error", func(t *testing.T) {
		m := mocks.NewHistClient(t)
		m.On("Record", mock.Anything).Return(errors.New("some error"))
		h := Handler{
			Client: m,
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
		}
		err := h.Handle(args.Args{
			Get: &args.GetCmd{Index: 101},
		})
		require.Error(t, err)
	})

	t.Run("handle tidy", func(t *testing.T) {
		m := mocks.NewHistClient(t)
		m.On("Tidy").Return(nil)
		h := Handler{
			Client: m,
		}
		err := h.Handle(args.Args{
			Tidy: &args.TidyCmd{},
		})
		require.NoError(t, err)
	})

	t.Run("tidy error", func(t *testing.T) {
		m := mocks.NewHistClient(t)
		m.On("Tidy", mock.Anything).Return(errors.New("some error"))
		h := Handler{
			Client: m,
		}
		err := h.Handle(args.Args{
			Tidy: &args.TidyCmd{},
		})
		require.Error(t, err)
	})

	t.Run("handle search", func(t *testing.T) {
		s := mocks.NewSearchClient(t)
		s.On("Show", mock.Anything).
			Return(nil)
		h := Handler{
			Searcher: s,
		}
		err := h.Handle(args.Args{
			Search: &args.SearchCmd{},
		})
		require.NoError(t, err)
	})

	t.Run("search error", func(t *testing.T) {
		s := mocks.NewSearchClient(t)
		s.On("Show", mock.Anything).
			Return(errors.New("some error"))
		h := Handler{
			Searcher: s,
		}
		err := h.Handle(args.Args{
			Search: &args.SearchCmd{},
		})
		require.Error(t, err)
	})

	t.Run("handle import", func(t *testing.T) {
		f, err := os.CreateTemp("", "hist-test-*")
		require.NoError(t, err)
		defer os.Remove(f.Name())
		m := mocks.NewHistClient(t)
		m.On("Import", mock.Anything).Return(nil)
		h := Handler{
			Client: m,
		}
		err = h.Handle(args.Args{
			Import: &args.ImportCmd{
				Path: f.Name(),
			},
		})
		require.NoError(t, err)
	})

	t.Run("import error", func(t *testing.T) {
		f, err := os.CreateTemp("", "hist-test-*")
		require.NoError(t, err)
		defer os.Remove(f.Name())
		m := mocks.NewHistClient(t)
		m.On("Import", mock.Anything).Return(errors.New("some error"))
		h := Handler{
			Client: m,
		}
		err = h.Handle(args.Args{
			Import: &args.ImportCmd{
				Path: f.Name(),
			},
		})
		require.Error(t, err)
	})

	t.Run("import error (file not found)", func(t *testing.T) {
		f, err := os.CreateTemp("", "hist-test-*")
		require.NoError(t, err)
		// remove the file to ensure that it is a valid file that does not exist
		err = os.Remove(f.Name())
		require.NoError(t, err)
		m := mocks.NewHistClient(t)
		h := Handler{
			Client: m,
		}
		err = h.Handle(args.Args{
			Import: &args.ImportCmd{
				Path: f.Name(),
			},
		})
		require.Error(t, err)
	})
}
