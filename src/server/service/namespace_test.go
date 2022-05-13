package service

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	mockrepository "tagservice/mocks/server/repository"
	mocks "tagservice/mocks/server/repository/transaction"
	"tagservice/server"
	"tagservice/server/model"
	"tagservice/server/repository"
	"tagservice/server/repository/transaction"
	"testing"
)

func TestNamespaceService_Create(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	var s = &NamespaceService{log: logger}
	t.Run(`create with error`, func(t *testing.T) {
		var repo = mockrepository.NewNamespace(t)
		s.namespaceRepository = repo
		repo.On(`Create`, mock.Anything, `ambar`).Return(model.Namespace{}, errors.New(``))
		//goland:noinspection GoContextTodo
		_, err := s.Create(context.TODO(), `ambar`)
		assert.Error(t, err)
		repo.AssertNumberOfCalls(t, `Create`, 1)

	})
	t.Run(`create without error`, func(t *testing.T) {
		var repo = mockrepository.NewNamespace(t)
		s.namespaceRepository = repo
		repo.On(`Create`, mock.Anything, `zamok`).Return(model.Namespace{}, nil)
		_, err := s.Create(context.TODO(), `zamok`)
		assert.NoError(t, err)
		repo.AssertNumberOfCalls(t, `Create`, 1)
	})
}

func TestNamespaceService_Delete(t *testing.T) {
	var errunknown = errors.New(`unknown`)
	tests := []struct {
		name                       string
		NamespaceGetByIdReturns    func() (model.Namespace, error)
		NamespaceDeleteByIdReturns error
		TxBeginTxReturns           func(tx transaction.Transaction) (transaction.Transaction, error)
		TxRollbackReturns          error
		TxCommitReturns            error
		RelationDeleteReturns      error
		err                        error
	}{
		{
			name: `not found namespace error`,
			NamespaceGetByIdReturns: func() (model.Namespace, error) {
				return model.Namespace{}, repository.ErrNotFound
			},
			TxBeginTxReturns: func(tx transaction.Transaction) (transaction.Transaction, error) { return nil, nil },
			err:              ErrNamespaceNotFound,
		},
		{
			name: `unknown err when getting namespace`,
			NamespaceGetByIdReturns: func() (model.Namespace, error) {
				return model.Namespace{}, errunknown
			},
			TxBeginTxReturns: func(tx transaction.Transaction) (transaction.Transaction, error) { return nil, nil },
			err:              errunknown,
		},
		{
			name: `error to start transaction`,
			NamespaceGetByIdReturns: func() (model.Namespace, error) {
				return model.Namespace{}, nil
			},
			TxBeginTxReturns: func(tx transaction.Transaction) (transaction.Transaction, error) {
				return tx, errunknown
			},
			err: errunknown,
		},
		{
			name: `error remove relation`,
			NamespaceGetByIdReturns: func() (model.Namespace, error) {
				return model.Namespace{}, nil
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
			NamespaceGetByIdReturns: func() (model.Namespace, error) {
				return model.Namespace{}, nil
			},
			TxBeginTxReturns: func(tx transaction.Transaction) (transaction.Transaction, error) {
				return tx, nil
			},
			RelationDeleteReturns: errors.New(`anything`),
			TxRollbackReturns:     errunknown,
			err:                   errunknown,
		},
		{
			name: `error remove namespace`,
			NamespaceGetByIdReturns: func() (model.Namespace, error) {
				return model.Namespace{}, nil
			},
			TxBeginTxReturns: func(tx transaction.Transaction) (transaction.Transaction, error) {
				return tx, nil
			},
			RelationDeleteReturns:      nil,
			NamespaceDeleteByIdReturns: errunknown,
			TxRollbackReturns:          nil,
			err:                        errunknown,
		},
		{
			name: `error rollback remove namespace`,
			NamespaceGetByIdReturns: func() (model.Namespace, error) {
				return model.Namespace{}, nil
			},
			TxBeginTxReturns: func(tx transaction.Transaction) (transaction.Transaction, error) {
				return tx, nil
			},
			RelationDeleteReturns:      nil,
			NamespaceDeleteByIdReturns: errors.New(`anything`),
			TxRollbackReturns:          errunknown,
			err:                        errunknown,
		},
		{
			name: `error commit`,
			NamespaceGetByIdReturns: func() (model.Namespace, error) {
				return model.Namespace{}, nil
			},
			TxBeginTxReturns: func(tx transaction.Transaction) (transaction.Transaction, error) {
				return tx, nil
			},
			RelationDeleteReturns:      nil,
			NamespaceDeleteByIdReturns: nil,
			TxCommitReturns:            errunknown,
			err:                        errunknown,
		},
		{
			name: `no error`,
			NamespaceGetByIdReturns: func() (model.Namespace, error) {
				return model.Namespace{}, nil
			},
			TxBeginTxReturns: func(tx transaction.Transaction) (transaction.Transaction, error) {
				return tx, nil
			},
			RelationDeleteReturns:      nil,
			NamespaceDeleteByIdReturns: nil,
			TxCommitReturns:            nil,
			err:                        nil,
		},
	}

	logger, _ := zap.NewDevelopment()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize service
			var s = func() server.Namespace {
				var tx = mocks.NewTransaction(t)
				tx.On(`BeginTx`, mock.Anything, mock.Anything).Return(tt.TxBeginTxReturns(tx)).Maybe()
				tx.On(`Rollback`, mock.Anything).Return(tt.TxRollbackReturns).Maybe()
				tx.On(`Commit`, mock.Anything).Return(tt.TxCommitReturns).Maybe()

				var namespacerepo = mockrepository.NewNamespace(t)
				namespacerepo.On(`GetById`, mock.Anything, mock.Anything).
					Return(tt.NamespaceGetByIdReturns()).Maybe()
				namespacerepo.On(`DeleteById`, mock.Anything, mock.Anything).
					Return(tt.NamespaceDeleteByIdReturns).Maybe()

				var relationRepo = mockrepository.NewRelation(t)
				relationRepo.On(`Delete`, mock.Anything, mock.Anything).
					Return(tt.RelationDeleteReturns).Maybe()

				return &NamespaceService{
					namespaceRepository: namespacerepo,
					relationRepository:  relationRepo,
					transaction:         tx,
					log:                 logger,
				}
			}()

			var err = s.Delete(context.TODO(), 88)
			if tt.err == nil {
				assert.NoError(t, err)
				return
			}
			assert.ErrorIs(t, err, tt.err)
		})
	}
}
