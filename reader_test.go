package streamer

import (
	"context"
	"errors"
	"sync"
	"testing"

	mock_streamer "github.com/clausthrane/streamer/.mocks"
	mock_source "github.com/clausthrane/streamer/.mocks/source"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/clausthrane/streamer/test/utils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

type fixtures struct {
	connector *mock_streamer.MockConnector
	stream    *mock_source.MockStream
}

func setup(t *testing.T, f func(context.Context, *fixtures)) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockConnector := mock_streamer.NewMockConnector(ctrl)
	mockStream := mock_source.NewMockStream(ctrl)
	utils.WithContext(func(ctx context.Context) {
		f(ctx, &fixtures{connector: mockConnector, stream: mockStream})
	})

}

func TestWhenCallingReceiveOnceThenTheHandlerIsInvokedWithTheResultOfTheValueReturnedFromTheStream(t *testing.T) {
	setup(t, func(ctx context.Context, fix *fixtures) {
		expected := "expected"
		fix.connector.EXPECT().Connect(gomock.Any()).Return(fix.stream, nil)
		fix.stream.EXPECT().Receive(gomock.Any()).Return(expected, nil)

		var v interface{}
		s := NewReader(ctx, fix.connector, func(ctx context.Context, in interface{}, out interface{}, e error) {
			v = in
		})
		assert.NoError(t, s.connect())

		s.receiveOnce()
		assert.Equal(t, expected, v)
	})
}

func TestWhenStartingAReaderThenTheStreamIsReadAsync(t *testing.T) {
	setup(t, func(ctx context.Context, fix *fixtures) {
		expected := "expected"
		fix.connector.EXPECT().Connect(gomock.Any()).Return(fix.stream, nil)
		fix.stream.EXPECT().Receive(gomock.Any()).Return(expected, nil).AnyTimes()

		w := sync.WaitGroup{}
		w.Add(1)

		var v interface{}
		s := NewReader(ctx, fix.connector, func(ctx context.Context, in interface{}, out interface{}, e error) {
			if v == nil {
				v = in
				w.Done()
			}
		})

		require.NoError(t, s.Start())
		w.Wait()
		s.Stop()
		assert.Equal(t, expected, v)
		goleak.VerifyNone(t)
	})
}

func TestWhenAConnectionFailsToOpenTheSteamIsNeverRead(t *testing.T) {
	setup(t, func(ctx context.Context, fix *fixtures) {
		fix.connector.EXPECT().Connect(gomock.Any()).Return(nil, errors.New("no"))
		s := NewReader(ctx, fix.connector, nil)
		assert.Error(t, s.Start())
		goleak.VerifyNone(t)
	})
}
