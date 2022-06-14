package namespace

import (
	"context"
	"errors"
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
		TxBeginTxReturns           func() (transaction.Transaction, error)
		RelationDeleteReturns      error
		err                        error
	}{
		{
			name: `not found namespace error`,
			NamespaceGetByIdReturns: func() (model.Namespace, error) {
				return model.Namespace{}, repository.ErrFindNamespace
			},
			TxBeginTxReturns: func() (transaction.Transaction, error) { return nil, nil },
			err:              ErrNamespaceNotFound,
		},
		{
			name: `unknown err when getting namespace`,
			NamespaceGetByIdReturns: func() (model.Namespace, error) {
				return model.Namespace{}, errunknown
			},
			TxBeginTxReturns: func() (transaction.Transaction, error) { return nil, nil },
			err:              errunknown,
		},
		{
			name: `error to start transaction`,
			NamespaceGetByIdReturns: func() (model.Namespace, error) {
				return model.Namespace{}, nil
			},
			TxBeginTxReturns: func() (transaction.Transaction, error) {
				return mocks.NewTransaction(t), errunknown
			},
			err: errunknown,
		},
		{
			name: `error remove relation`,
			NamespaceGetByIdReturns: func() (model.Namespace, error) {
				return model.Namespace{}, nil
			},
			TxBeginTxReturns: func() (transaction.Transaction, error) {
				var tx = mocks.NewTransaction(t)

				var rel = mockrepository.NewRelation(t)
				rel.On(`Delete`, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errunknown)
				tx.On(`Relation`).Return(rel)
				tx.On(`Rollback`, mock.Anything).Return(nil)

				return tx, nil
			},
			RelationDeleteReturns: errunknown,
			err:                   errunknown,
		},
		{
			name: `error rollback remove relation`,
			NamespaceGetByIdReturns: func() (model.Namespace, error) {
				return model.Namespace{}, nil
			},
			TxBeginTxReturns: func() (transaction.Transaction, error) {
				var tx = mocks.NewTransaction(t)

				var rel = mockrepository.NewRelation(t)
				rel.On(`Delete`, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New(``))
				tx.On(`Relation`).Return(rel)

				tx.On(`Rollback`, mock.Anything).Return(errunknown)
				return tx, nil
			},
			RelationDeleteReturns: errors.New(`anything`),
			err:                   errunknown,
		},
		{
			name: `error remove namespace`,
			NamespaceGetByIdReturns: func() (model.Namespace, error) {
				return model.Namespace{}, nil
			},
			TxBeginTxReturns: func() (transaction.Transaction, error) {
				var tx = mocks.NewTransaction(t)

				var rel = mockrepository.NewRelation(t)
				rel.On(`Delete`, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				tx.On(`Relation`).Return(rel)

				var ns = mockrepository.NewNamespace(t)
				ns.On(`DeleteById`, mock.Anything, mock.Anything).Return(errunknown)
				tx.On(`Namespace`).Return(ns)

				tx.On(`Rollback`, mock.Anything).Return(nil)
				return tx, nil
			},
			RelationDeleteReturns:      nil,
			NamespaceDeleteByIdReturns: errunknown,
			err:                        errunknown,
		},
		{
			name: `error rollback remove namespace`,
			NamespaceGetByIdReturns: func() (model.Namespace, error) {
				return model.Namespace{}, nil
			},
			TxBeginTxReturns: func() (transaction.Transaction, error) {
				var tx = mocks.NewTransaction(t)

				var rel = mockrepository.NewRelation(t)
				rel.On(`Delete`, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				tx.On(`Relation`).Return(rel)

				var ns = mockrepository.NewNamespace(t)
				ns.On(`DeleteById`, mock.Anything, mock.Anything).Return(errunknown)
				tx.On(`Namespace`).Return(ns)

				tx.On(`Rollback`, mock.Anything).Return(errunknown)
				return tx, nil
			},
			RelationDeleteReturns:      nil,
			NamespaceDeleteByIdReturns: errors.New(`anything`),
			err:                        errunknown,
		},
		{
			name: `error commit`,
			NamespaceGetByIdReturns: func() (model.Namespace, error) {
				return model.Namespace{}, nil
			},
			TxBeginTxReturns: func() (transaction.Transaction, error) {
				var tx = mocks.NewTransaction(t)

				var rel = mockrepository.NewRelation(t)
				rel.On(`Delete`, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				tx.On(`Relation`).Return(rel)

				var ns = mockrepository.NewNamespace(t)
				ns.On(`DeleteById`, mock.Anything, mock.Anything).Return(nil)
				tx.On(`Namespace`).Return(ns)

				tx.On(`Commit`, mock.Anything).Return(errunknown)
				return tx, nil
			},
			RelationDeleteReturns:      nil,
			NamespaceDeleteByIdReturns: nil,
			err:                        errunknown,
		},
		{
			name: `no error`,
			NamespaceGetByIdReturns: func() (model.Namespace, error) {
				return model.Namespace{}, nil
			},
			TxBeginTxReturns: func() (transaction.Transaction, error) {
				var tx = mocks.NewTransaction(t)

				var rel = mockrepository.NewRelation(t)
				rel.On(`Delete`, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				tx.On(`Relation`).Return(rel)

				var ns = mockrepository.NewNamespace(t)
				ns.On(`DeleteById`, mock.Anything, mock.Anything).Return(nil)
				tx.On(`Namespace`).Return(ns)

				tx.On(`Commit`, mock.Anything).Return(nil)
				return tx, nil
			},
			RelationDeleteReturns:      nil,
			NamespaceDeleteByIdReturns: nil,
			err:                        nil,
		},
	}

	logger, _ := zap.NewDevelopment()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize service
			var s = func() tagservice.Namespace {
				var tx = mocks.NewTransactioner(t)
				tx.On(`BeginTx`, mock.Anything, mock.Anything).Return(tt.TxBeginTxReturns()).Maybe()

				var namespacerepo = mockrepository.NewNamespace(t)
				namespacerepo.On(`GetById`, mock.Anything, mock.Anything).
					Return(tt.NamespaceGetByIdReturns()).Maybe()

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
