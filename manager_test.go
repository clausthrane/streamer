package streamer

import (
	"context"
	"testing"

	mock_streamer "github.com/clausthrane/streamer/.mocks"
	"github.com/golang/mock/gomock"

	"github.com/clausthrane/streamer/test/utils"
)

type managerFixtures struct {
	connector *mock_streamer.MockConnector
}

func setupManager(t *testing.T, f func(context.Context, *StreamManager, *managerFixtures)) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	connector := mock_streamer.NewMockConnector(ctrl)
	utils.WithContext(func(ctx context.Context) {
		f(ctx, NewStreamManager(ctx), &managerFixtures{connector: connector})
	})
}

func TestGivenAManagerHandlerFuncCanBeRegistered(t *testing.T) {
	setupManager(t, func(ctx context.Context, manager *StreamManager, fix *managerFixtures) {
		fix.connector.EXPECT().ID().Return("id")
		manager.RegisterFunc(fix.connector, func(ctx context.Context, in interface{}, out interface{}, e error) {
		})
	})
}
