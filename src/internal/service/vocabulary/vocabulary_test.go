package vocabulary

import (
	"context"
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
		data             *model.VocabularyData
		assert           func(t *testing.T, voc *model.Vocabulary, err error)
		repositoryReturn func() (*model.Vocabulary, error)
	}{
		{
			name: "Create Vocabulary Successfully",
			data: &model.VocabularyData{
				Name: `XXX`,
			},
			assert: func(t *testing.T, voc *model.Vocabulary, err error) {
				assert.NoError(t, err)
				assert.Equal(t, uint64(121), voc.ID)
				assert.Equal(t, `XXX`, voc.Data.Name)
			},
			repositoryReturn: func() (*model.Vocabulary, error) {
				return &model.Vocabulary{ID: 121, Data: model.VocabularyData{Name: `XXX`}}, nil
			},
		},
		// Add more test cases @TODO
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.SetUp(t)
			var ctx = context.Background()
			vocabularyRepository := mock.Mock[repository.Vocabulary]()
			service := VocabularyService{
				log:                  zap.NewNop(),
				vocabularyRepository: vocabularyRepository,
			}

			mock.When(vocabularyRepository.Create(mock.Exact[context.Context](ctx), mock.Any[*model.VocabularyData]())).
				ThenAnswer(func(args []any) []any {
					voc, err := tt.repositoryReturn()
					return []any{voc, err}
				})

			result, err := service.Create(ctx, tt.data)
			tt.assert(t, result, err)
		})
	}
}
