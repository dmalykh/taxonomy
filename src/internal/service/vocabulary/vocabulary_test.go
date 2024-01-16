package vocabulary

import (
	"context"
	"errors"
	"github.com/dmalykh/taxonomy/taxonomy/model"
	"github.com/dmalykh/taxonomy/taxonomy/repository"
	"github.com/ovechkin-dm/mockio/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestVocabularyService_Create(t *testing.T) {
	tests := []struct {
		name             string
		mockData         *model.VocabularyData
		expectedResult   *model.Vocabulary
		expectedError    error
		repositoryReturn func() (*model.Vocabulary, error)
	}{
		{
			name:     "Create Vocabulary Successfully",
			mockData: &model.VocabularyData{
				// provide necessary data
			},
			expectedResult: &model.Vocabulary{
				// provide necessary data
			},
			repositoryReturn: func() (*model.Vocabulary, error) {

			},
		},
		{
			name:     "Error Creating Vocabulary",
			mockData: &model.VocabularyData{
				// provide necessary data
			},
			expectedResult: nil,
			expectedError:  errors.New("some error"),
			repositoryReturn: func() (*model.Vocabulary, error) {

			},
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mock.SetUp(t)
			var ctx = context.Background()
			vocabularyRepository := mock.Mock[repository.Vocabulary]()
			service := VocabularyService{
				log:                  zap.NewNop(),
				vocabularyRepository: vocabularyRepository,
			}

			mock.When(vocabularyRepository.Create(mock.Exact[context.Context](ctx), tt.mockData)).
				ThenAnswer(func(args []any) []any {
					voc, err := tt.repositoryReturn()
					return []any{voc, err}
				})

			// Act
			result, err := service.Create(context.Background(), tt.mockData)

			// Assert
			assert.Equal(t, tt.expectedResult, result)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}
