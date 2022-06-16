package tag

import (
	"context"
	"errors"
	mockservice "github.com/dmalykh/tagservice/mocks"
	mockrepository "github.com/dmalykh/tagservice/mocks/repository"
	mocks "github.com/dmalykh/tagservice/mocks/repository/transaction"
	"github.com/dmalykh/tagservice/tagservice"
	"github.com/dmalykh/tagservice/tagservice/model"
	"github.com/dmalykh/tagservice/tagservice/repository"
	"github.com/dmalykh/tagservice/tagservice/repository/transaction"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"testing"
)

func TestTagService_Delete(t *testing.T) {
	var errunknown = errors.New(`unknown`)
	tests := []struct {
		name              string
		TagGetByIdReturns func() (model.Tag, error)
		TxBeginTxReturns  func() (transaction.Transaction, error)
		err               error
	}{
		{
			name: `not found tag error`,
			TagGetByIdReturns: func() (model.Tag, error) {
				return model.Tag{}, repository.ErrFindTag
			},
			TxBeginTxReturns: func() (transaction.Transaction, error) { return nil, nil },
			err:              tagservice.ErrTagNotFound,
		},
		{
			name: `unknown err when getting tag`,
			TagGetByIdReturns: func() (model.Tag, error) {
				return model.Tag{}, errunknown
			},
			TxBeginTxReturns: func() (transaction.Transaction, error) { return nil, nil },
			err:              errunknown,
		},
		{
			name: `error to start Transaction`,
			TagGetByIdReturns: func() (model.Tag, error) {
				return model.Tag{}, nil
			},
			TxBeginTxReturns: func() (transaction.Transaction, error) {
				return nil, errunknown
			},
			err: errunknown,
		},
		{
			name: `error remove relation`,
			TagGetByIdReturns: func() (model.Tag, error) {
				return model.Tag{}, nil
			},
			TxBeginTxReturns: func() (transaction.Transaction, error) {
				var tx = mocks.NewTransaction(t)

				var rel = mockrepository.NewRelation(t)
				rel.On(`Delete`, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errunknown)
				tx.On(`Relation`).Return(rel)
				tx.On(`Rollback`, mock.Anything).Return(nil)

				return tx, nil
			},
			err: errunknown,
		},
		{
			name: `error rollback remove relation`,
			TagGetByIdReturns: func() (model.Tag, error) {
				return model.Tag{}, nil
			},
			TxBeginTxReturns: func() (transaction.Transaction, error) {
				var tx = mocks.NewTransaction(t)

				var rel = mockrepository.NewRelation(t)
				rel.On(`Delete`, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New(``))
				tx.On(`Relation`).Return(rel)

				tx.On(`Rollback`, mock.Anything).Return(errunknown)
				return tx, nil
			},
			err: errunknown,
		},
		{
			name: `error remove tag`,
			TagGetByIdReturns: func() (model.Tag, error) {
				return model.Tag{}, nil
			},
			TxBeginTxReturns: func() (transaction.Transaction, error) {
				var tx = mocks.NewTransaction(t)

				var rel = mockrepository.NewRelation(t)
				rel.On(`Delete`, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				tx.On(`Relation`).Return(rel)

				var tag = mockrepository.NewTag(t)
				tag.On(`DeleteById`, mock.Anything, mock.Anything).Return(errunknown)
				tx.On(`Tag`).Return(tag)

				tx.On(`Rollback`, mock.Anything).Return(nil)
				return tx, nil
			},
			err: errunknown,
		},
		{
			name: `error rollback remove tag`,
			TagGetByIdReturns: func() (model.Tag, error) {
				return model.Tag{}, nil
			},
			TxBeginTxReturns: func() (transaction.Transaction, error) {
				var tx = mocks.NewTransaction(t)

				var rel = mockrepository.NewRelation(t)
				rel.On(`Delete`, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				tx.On(`Relation`).Return(rel)

				var tag = mockrepository.NewTag(t)
				tag.On(`DeleteById`, mock.Anything, mock.Anything).Return(errors.New(``))
				tx.On(`Tag`).Return(tag)

				tx.On(`Rollback`, mock.Anything).Return(errunknown)
				return tx, nil
			},
			err: errunknown,
		},
		{
			name: `error commit`,
			TagGetByIdReturns: func() (model.Tag, error) {
				return model.Tag{}, nil
			},
			TxBeginTxReturns: func() (transaction.Transaction, error) {
				var tx = mocks.NewTransaction(t)

				var rel = mockrepository.NewRelation(t)
				rel.On(`Delete`, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				tx.On(`Relation`).Return(rel)

				var tag = mockrepository.NewTag(t)
				tag.On(`DeleteById`, mock.Anything, mock.Anything).Return(nil)
				tx.On(`Tag`).Return(tag)

				tx.On(`Commit`, mock.Anything).Return(errunknown)
				return tx, nil
			},
			err: errunknown,
		},
		{
			name: `no error`,
			TagGetByIdReturns: func() (model.Tag, error) {
				return model.Tag{}, nil
			},
			TxBeginTxReturns: func() (transaction.Transaction, error) {
				var tx = mocks.NewTransaction(t)

				var rel = mockrepository.NewRelation(t)
				rel.On(`Delete`, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				tx.On(`Relation`).Return(rel)

				var tag = mockrepository.NewTag(t)
				tag.On(`DeleteById`, mock.Anything, mock.Anything).Return(nil)
				tx.On(`Tag`).Return(tag)

				tx.On(`Commit`, mock.Anything).Return(nil)
				return tx, nil
			},
			err: nil,
		},
	}

	logger, _ := zap.NewDevelopment()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize service
			var s = func() tagservice.Tag {
				var tx = mocks.NewTransactioner(t)
				tx.On(`BeginTx`, mock.Anything, mock.Anything).Return(tt.TxBeginTxReturns()).Maybe()

				var tagrepo = mockrepository.NewTag(t)
				tagrepo.On(`GetById`, mock.Anything, mock.Anything).
					Return(tt.TagGetByIdReturns()).Maybe()

				return &TagService{
					tagRepository:      tagrepo,
					relationRepository: mockrepository.NewRelation(t),
					transaction:        tx,
					log:                logger,
				}
			}()

			//goland:noinspection GoContextTodo
			var err = s.Delete(context.TODO(), 88)
			if tt.err == nil {
				assert.NoError(t, err)
				return
			}
			assert.ErrorIs(t, err, tt.err)
		})
	}
}

func TestTagService_GetRelationEntities(t1 *testing.T) {

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
				return model.Namespace{Id: 11}, nil
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
				return model.Namespace{Id: 11}, errors.New(``)
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
	// Should receive all laptops that has: "RAM" (512 or 1024) and "CPU" (2.8 or 3.2) and "Display size" (between 13 and 15)

	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			var relationrepo = mockrepository.NewRelation(t1)
			relationrepo.On(`Get`, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tt.RelationGetReturns()).Maybe()

			var namespaceservice = mockservice.NewNamespace(t1)
			namespaceservice.On(`GetByName`, mock.Anything, mock.Anything).Return(tt.NamespaceGetByNameReturns())

			t := &TagService{
				log:                zap.NewNop(),
				namespaceService:   namespaceservice,
				relationRepository: relationrepo,
			}

			relations, err := t.GetRelations(context.TODO(), &model.EntityFilter{
				TagId:     tt.tagGroups,
				Namespace: []string{`any`},
			})
			tt.wantErr(t1, err)
			tt.want(t1, relations)

		})
	}
}
