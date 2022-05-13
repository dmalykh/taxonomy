package service

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	mockservice "tagservice/mocks/server"
	mockrepository "tagservice/mocks/server/repository"
	mocks "tagservice/mocks/server/repository/transaction"
	"tagservice/server"
	"tagservice/server/model"
	"tagservice/server/repository"
	"tagservice/server/repository/transaction"
	"testing"
)

func TestTagService_Delete(t *testing.T) {
	var errunknown = errors.New(`unknown`)
	tests := []struct {
		name                  string
		TagGetByIdReturns     func() (model.Tag, error)
		TagDeleteByIdReturns  error
		TxBeginTxReturns      func(tx transaction.Transaction) (transaction.Transaction, error)
		TxRollbackReturns     error
		TxCommitReturns       error
		RelationDeleteReturns error
		err                   error
	}{
		{
			name: `not found tag error`,
			TagGetByIdReturns: func() (model.Tag, error) {
				return model.Tag{}, repository.ErrNotFound
			},
			TxBeginTxReturns: func(tx transaction.Transaction) (transaction.Transaction, error) { return nil, nil },
			err:              ErrTagNotFound,
		},
		{
			name: `unknown err when getting tag`,
			TagGetByIdReturns: func() (model.Tag, error) {
				return model.Tag{}, errunknown
			},
			TxBeginTxReturns: func(tx transaction.Transaction) (transaction.Transaction, error) { return nil, nil },
			err:              errunknown,
		},
		{
			name: `error to start transaction`,
			TagGetByIdReturns: func() (model.Tag, error) {
				return model.Tag{}, nil
			},
			TxBeginTxReturns: func(tx transaction.Transaction) (transaction.Transaction, error) {
				return tx, errunknown
			},
			err: errunknown,
		},
		{
			name: `error remove relation`,
			TagGetByIdReturns: func() (model.Tag, error) {
				return model.Tag{}, nil
			},
			TxBeginTxReturns: func(tx transaction.Transaction) (transaction.Transaction, error) {
				return tx, nil
			},
			RelationDeleteReturns: errunknown,
			TxRollbackReturns:     nil,
			err:                   errunknown,
		},
		{
			name: `error rollback remove relation`,
			TagGetByIdReturns: func() (model.Tag, error) {
				return model.Tag{}, nil
			},
			TxBeginTxReturns: func(tx transaction.Transaction) (transaction.Transaction, error) {
				return tx, nil
			},
			RelationDeleteReturns: errors.New(`anything`),
			TxRollbackReturns:     errunknown,
			err:                   errunknown,
		},
		{
			name: `error remove tag`,
			TagGetByIdReturns: func() (model.Tag, error) {
				return model.Tag{}, nil
			},
			TxBeginTxReturns: func(tx transaction.Transaction) (transaction.Transaction, error) {
				return tx, nil
			},
			RelationDeleteReturns: nil,
			TagDeleteByIdReturns:  errunknown,
			TxRollbackReturns:     nil,
			err:                   errunknown,
		},
		{
			name: `error rollback remove tag`,
			TagGetByIdReturns: func() (model.Tag, error) {
				return model.Tag{}, nil
			},
			TxBeginTxReturns: func(tx transaction.Transaction) (transaction.Transaction, error) {
				return tx, nil
			},
			RelationDeleteReturns: nil,
			TagDeleteByIdReturns:  errors.New(`anything`),
			TxRollbackReturns:     errunknown,
			err:                   errunknown,
		},
		{
			name: `error commit`,
			TagGetByIdReturns: func() (model.Tag, error) {
				return model.Tag{}, nil
			},
			TxBeginTxReturns: func(tx transaction.Transaction) (transaction.Transaction, error) {
				return tx, nil
			},
			RelationDeleteReturns: nil,
			TagDeleteByIdReturns:  nil,
			TxCommitReturns:       errunknown,
			err:                   errunknown,
		},
		{
			name: `no error`,
			TagGetByIdReturns: func() (model.Tag, error) {
				return model.Tag{}, nil
			},
			TxBeginTxReturns: func(tx transaction.Transaction) (transaction.Transaction, error) {
				return tx, nil
			},
			RelationDeleteReturns: nil,
			TagDeleteByIdReturns:  nil,
			TxCommitReturns:       nil,
			err:                   nil,
		},
	}

	logger, _ := zap.NewDevelopment()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize service
			var s = func() server.Tag {
				var tx = mocks.NewTransaction(t)
				tx.On(`BeginTx`, mock.Anything, mock.Anything).Return(tt.TxBeginTxReturns(tx)).Maybe()
				tx.On(`Rollback`, mock.Anything).Return(tt.TxRollbackReturns).Maybe()
				tx.On(`Commit`, mock.Anything).Return(tt.TxCommitReturns).Maybe()

				var tagrepo = mockrepository.NewTag(t)
				tagrepo.On(`GetById`, mock.Anything, mock.Anything).
					Return(tt.TagGetByIdReturns()).Maybe()
				tagrepo.On(`DeleteById`, mock.Anything, mock.Anything).
					Return(tt.TagDeleteByIdReturns).Maybe()

				var relationRepo = mockrepository.NewRelation(t)
				relationRepo.On(`Delete`, mock.Anything, mock.Anything).
					Return(tt.RelationDeleteReturns).Maybe()

				return &TagService{
					tagRepository:      tagrepo,
					relationRepository: relationRepo,
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
		tagGroups                 [][]uint64
		RelationGetReturns        func() ([]model.Relation, error)
		NamespaceGetByNameReturns func() (model.Namespace, error)
		want                      func(t assert.TestingT, relations []model.Relation)
		wantErr                   assert.ErrorAssertionFunc
	}{
		{
			name:      `no relations`,
			tagGroups: [][]uint64{{34}, {92, 96}},
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
			tagGroups: [][]uint64{{34}, {92, 96}},
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
				return assert.ErrorIs(t, err, ErrTagNamespaceNotFound)
			},
		},
		{
			name:      `duplicates removed`,
			tagGroups: [][]uint64{{34}, {92, 96}, {33}},
			NamespaceGetByNameReturns: func() (model.Namespace, error) {
				return model.Namespace{Id: 11}, nil
			},
			RelationGetReturns: func() ([]model.Relation, error) {
				return []model.Relation{
					{EntityId: 22},
					{EntityId: 22},
					{EntityId: 53},
				}, nil
			},
			want: func(t assert.TestingT, relations []model.Relation) {
				assert.Len(t, relations, 2)
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
			},
		},
		{
			name: `no groups`,
			NamespaceGetByNameReturns: func() (model.Namespace, error) {
				return model.Namespace{Id: 11}, nil
			},
			tagGroups: [][]uint64{{}, {}, {}},
			RelationGetReturns: func() ([]model.Relation, error) {
				return []model.Relation{
					{EntityId: 93},
					{EntityId: 22},
					{EntityId: 22},
					{EntityId: 53},
					{EntityId: 53},
				}, nil
			},
			want: func(t assert.TestingT, relations []model.Relation) {
				assert.Len(t, relations, 3)
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
			},
		},
		{
			name: `empty groups`,
			NamespaceGetByNameReturns: func() (model.Namespace, error) {
				return model.Namespace{Id: 11}, nil
			},
			RelationGetReturns: func() ([]model.Relation, error) {
				return []model.Relation{
					{EntityId: 93},
					{EntityId: 22},
					{EntityId: 22},
					{EntityId: 53},
				}, nil
			},
			want: func(t assert.TestingT, relations []model.Relation) {
				assert.Len(t, relations, 0)
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
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

			relations, err := t.GetRelationEntities(context.TODO(), `any`, tt.tagGroups)
			tt.wantErr(t1, err)
			tt.want(t1, relations)

		})
	}
}
