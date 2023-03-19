package service

import (
	"github.com/richmondgoh8/go-casbin/graph/model"
	mock_ports "github.com/richmondgoh8/go-casbin/internal/mocks/core/ports"
	"github.com/richmondgoh8/go-casbin/pkg/utils"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestGetHeartbeat(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	ctx := utils.GetTestGinContext()

	tests := []struct {
		name         string
		doMockRepo   func(repository *mock_ports.MockHealthRepository)
		expectedResp *model.HealthResponse
		err          error
	}{
		{
			name: "Test Case Positive",
			doMockRepo: func(repository *mock_ports.MockHealthRepository) {
				repository.EXPECT().GetHealth().Return(&model.HealthResponse{
					Message: "OK",
				})
			},
			expectedResp: &model.HealthResponse{
				Message: "OK",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockHealthRepo := mock_ports.NewMockHealthRepository(mockCtrl)
			tc.doMockRepo(mockHealthRepo)

			healthSvc := NewHealthSvc(mockHealthRepo)
			data := healthSvc.GetHeartbeat(ctx)
			assert.Equal(t, tc.expectedResp, data)
		})
	}
}
