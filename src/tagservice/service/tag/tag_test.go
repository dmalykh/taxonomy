package tag_test

import (
	"context"
	"errors"
	"testing"

	mockservice "github.com/dmalykh/tagservice/mocks"
	mockrepository "github.com/dmalykh/tagservice/mocks/repository"
	mocks "github.com/dmalykh/tagservice/mocks/repository/transaction"
	"github.com/dmalykh/tagservice/tagservice"
	"github.com/dmalykh/tagservice/tagservice/model"
	"github.com/dmalykh/tagservice/tagservice/repository"
	"github.com/dmalykh/tagservice/tagservice/repository/transaction"
	"github.com/dmalykh/tagservice/tagservice/service/tag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestTagService_Delete(t *testing.T) {
	t.Parallel()

	errunknown := errors.New(`unknown`)

	tests := []struct {
		name              string
		TagGetByIDReturns func() (model.Tag, error)
		TxBeginTxReturns  func() (transaction.Transaction, error)
		err               error
	}{
		{
			name: `not found tag error`,
			TagGetByIDReturns: func() (model.Tag, error) {
				return model.Tag{}, repository.ErrFindTag
			},
			TxBeginTxReturns: func() (transaction.Transaction, error) { return nil, nil },
			err:              tagservice.ErrTagNotFound,
		},
		{
			name: `unknown err when getting tag`,
			TagGetByIDReturns: func() (model.Tag, error) {
				return model.Tag{}, errunknown
			},
			TxBeginTxReturns: func() (transaction.Transaction, error) { return nil, nil },
			err:              errunknown,
		},
		{
			name: `error to start Transaction`,
			TagGetByIDReturns: func() (model.Tag, error) {
				return model.Tag{}, nil
			},
			TxBeginTxReturns: func() (transaction.Transaction, error) {
				return nil, errunknown
			},
			err: errunknown,
		},
		{
			name: `error remove relation`,
			TagGetByIDReturns: func() (model.Tag, error) {
				return model.Tag{}, nil
			},
			TxBeginTxReturns: func() (transaction.Transaction, error) {
				tx := mocks.NewTransaction(t)

				rel := mockrepository.NewRelation(t)
				rel.On(`Delete`, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errunknown)
				tx.On(`Relation`).Return(rel)
				tx.On(`Rollback`, mock.Anything).Return(nil)

				return tx, nil
			},
			err: errunknown,
		},
		{
			name: `error rollback remove relation`,
			TagGetByIDReturns: func() (model.Tag, error) {
				return model.Tag{}, nil
			},
			TxBeginTxReturns: func() (transaction.Transaction, error) {
				tx := mocks.NewTransaction(t)

				rel := mockrepository.NewRelation(t)
				rel.On(`Delete`, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(errors.New(``))
				tx.On(`Relation`).Return(rel)

				tx.On(`Rollback`, mock.Anything).Return(errunknown)

				return tx, nil
			},
			err: errunknown,
		},
		{
			name: `error remove tag`,
			TagGetByIDReturns: func() (model.Tag, error) {
				return model.Tag{}, nil
			},
			TxBeginTxReturns: func() (transaction.Transaction, error) {
				tx := mocks.NewTransaction(t)

				rel := mockrepository.NewRelation(t)
				rel.On(`Delete`, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				tx.On(`Relation`).Return(rel)

				tag := mockrepository.NewTag(t)
				tag.On(`DeleteByID`, mock.Anything, mock.Anything).Return(errunknown)
				tx.On(`Tag`).Return(tag)

				tx.On(`Rollback`, mock.Anything).Return(nil)

				return tx, nil
			},
			err: errunknown,
		},
		{
			name: `error rollback remove tag`,
			TagGetByIDReturns: func() (model.Tag, error) {
				return model.Tag{}, nil
			},
			TxBeginTxReturns: func() (transaction.Transaction, error) {
				tx := mocks.NewTransaction(t)

				rel := mockrepository.NewRelation(t)
				rel.On(`Delete`, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				tx.On(`Relation`).Return(rel)

				tag := mockrepository.NewTag(t)
				tag.On(`DeleteByID`, mock.Anything, mock.Anything).Return(errors.New(``))
				tx.On(`Tag`).Return(tag)

				tx.On(`Rollback`, mock.Anything).Return(errunknown)

				return tx, nil
			},
			err: errunknown,
		},
		{
			name: `error commit`,
			TagGetByIDReturns: func() (model.Tag, error) {
				return model.Tag{}, nil
			},
			TxBeginTxReturns: func() (transaction.Transaction, error) {
				tx := mocks.NewTransaction(t)

				rel := mockrepository.NewRelation(t)
				rel.On(`Delete`, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				tx.On(`Relation`).Return(rel)

				tag := mockrepository.NewTag(t)
				tag.On(`DeleteByID`, mock.Anything, mock.Anything).Return(nil)
				tx.On(`Tag`).Return(tag)

				tx.On(`Commit`, mock.Anything).Return(errunknown)

				return tx, nil
			},
			err: errunknown,
		},
		{
			name: `no error`,
			TagGetByIDReturns: func() (model.Tag, error) {
				return model.Tag{}, nil
			},
			TxBeginTxReturns: func() (transaction.Transaction, error) {
				tx := mocks.NewTransaction(t)

				rel := mockrepository.NewRelation(t)
				rel.On(`Delete`, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				tx.On(`Relation`).Return(rel)

				tag := mockrepository.NewTag(t)
				tag.On(`DeleteByID`, mock.Anything, mock.Anything).Return(nil)
				tx.On(`Tag`).Return(tag)

				tx.On(`Commit`, mock.Anything).Return(nil)

				return tx, nil
			},
			err: nil,
		},
	}

	logger, _ := zap.NewDevelopment()

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Initialize service
			s := func() tagservice.Tag {
				tx := mocks.NewTransactioner(t)
				tx.On(`BeginTx`, mock.Anything, mock.Anything).Return(tt.TxBeginTxReturns()).Maybe()

				tagrepo := mockrepository.NewTag(t)
				tagrepo.On(`GetByID`, mock.Anything, mock.Anything).
					Return(tt.TagGetByIDReturns()).Maybe()

				return tag.New(&tag.Config{
					Transaction:        tx,
					TagRepository:      tagrepo,
					RelationRepository: mockrepository.NewRelation(t),
					Logger:             logger,
				})
			}()

			//goland:noinspection GoContextTodo
			err := s.Delete(context.TODO(), 88)
			if tt.err == nil {
				assert.NoError(t, err)

				return
			}
			assert.ErrorIs(t, err, tt.err)
		})
	}
}

//goland:noinspection GoContextTodo
func TestTagService_GetRelations(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                      string
		tagGroups                 [][]uint
		RelationGetReturns        func() ([]model.Relation, error)
		NamespaceGetByNameReturns func() (model.Namespace, error)
		want                      func(t assert.TestingT, relations []model.Relation)
		wantErr                   assert.ErrorAssertionFunc
	}{
		{
			name:      `no relations`,
			tagGroups: [][]uint{{34}, {92, 96}},
			NamespaceGetByNameReturns: func() (model.Namespace, error) {
				return model.Namespace{ID: 11}, nil
			},
			RelationGetReturns: func() ([]model.Relation, error) {
				return nil, nil
			},
			want: func(t assert.TestingT, relations []model.Relation) {
				assert.Len(t, relations, 0)
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
			},
		},
		{
			name:      `namespace not found`,
			tagGroups: [][]uint{{34}, {92, 96}},
			NamespaceGetByNameReturns: func() (model.Namespace, error) {
				return model.Namespace{ID: 11}, errors.New(``)
			},
			RelationGetReturns: func() ([]model.Relation, error) {
				return nil, nil
			},
			want: func(t assert.TestingT, relations []model.Relation) {
				assert.Len(t, relations, 0)
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, tagservice.ErrTagNamespaceNotFound)
			},
		},
	}

	// Create categories "RAM", "CPU", "Display size"
	// Should receive all laptops that has: "RAM" (512 or 1024) and "CPU" (2.8 or 3.2) and "Display size"
	// (between 13 and 15)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			relationrepo := mockrepository.NewRelation(t)
			relationrepo.On(`Get`, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
				Return(tt.RelationGetReturns()).Maybe()

			namespaceservice := mockservice.NewNamespace(t)
			namespaceservice.On(`GetByName`, mock.Anything, mock.Anything).Return(tt.NamespaceGetByNameReturns())

			tagService := tag.New(&tag.Config{
				NamespaceService:   namespaceservice,
				RelationRepository: relationrepo,
				Logger:             zap.NewNop(),
			})

			relations, err := tagService.GetRelations(context.TODO(), &model.EntityFilter{
				TagID:     tt.tagGroups,
				Namespace: []string{`any`},
			})
			tt.wantErr(t, err)
			tt.want(t, relations)
		})
	}
}
