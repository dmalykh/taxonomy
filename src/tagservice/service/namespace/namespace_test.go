package namespace_test

import (
	"context"
	"errors"
	"testing"

	mockrepository "github.com/dmalykh/tagservice/mocks/repository"
	mocks "github.com/dmalykh/tagservice/mocks/repository/transaction"
	"github.com/dmalykh/tagservice/tagservice"
	"github.com/dmalykh/tagservice/tagservice/model"
	"github.com/dmalykh/tagservice/tagservice/repository"
	"github.com/dmalykh/tagservice/tagservice/repository/transaction"
	"github.com/dmalykh/tagservice/tagservice/service/namespace"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

//goland:noinspection GoContextTodo
func TestNamespaceService_Create(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run(`create with error`, func(t *testing.T) {
		repo := mockrepository.NewNamespace(t)
		s := namespace.New(&namespace.Config{
			Logger:              logger,
			NamespaceRepository: repo,
		})
		repo.On(`Create`, mock.Anything, `ambar`).Return(model.Namespace{}, errors.New(``))
		//goland:noinspection GoContextTodo
		_, err := s.Create(context.TODO(), `ambar`)
		assert.Error(t, err)
		repo.AssertNumberOfCalls(t, `Create`, 1)
	})

	t.Run(`create without error`, func(t *testing.T) {
		t.Parallel()
		repo := mockrepository.NewNamespace(t)
		s := namespace.New(&namespace.Config{
			Logger:              logger,
			NamespaceRepository: repo,
		})
		repo.On(`Create`, mock.Anything, `zamok`).Return(model.Namespace{}, nil)
		_, err := s.Create(context.TODO(), `zamok`)
		assert.NoError(t, err)
		repo.AssertNumberOfCalls(t, `Create`, 1)
	})
}

//goland:noinspection GoContextTodo
func TestNamespaceService_Delete(t *testing.T) {
	t.Parallel()

	errunknown := errors.New(`unknown`)

	tests := []struct {
		name                       string
		NamespaceGetByIDReturns    func() (model.Namespace, error)
		NamespaceDeleteByIDReturns error
		TxBeginTxReturns           func() (transaction.Transaction, error)
		RelationDeleteReturns      error
		err                        error
	}{
		{
			name: `not found namespace error`,
			NamespaceGetByIDReturns: func() (model.Namespace, error) {
				return model.Namespace{}, repository.ErrFindNamespace
			},
			TxBeginTxReturns: func() (transaction.Transaction, error) { return nil, nil },
			err:              tagservice.ErrNamespaceNotFound,
		},
		{
			name: `unknown err when getting namespace`,
			NamespaceGetByIDReturns: func() (model.Namespace, error) {
				return model.Namespace{}, errunknown
			},
			TxBeginTxReturns: func() (transaction.Transaction, error) { return nil, nil },
			err:              errunknown,
		},
		{
			name: `error to start transaction`,
			NamespaceGetByIDReturns: func() (model.Namespace, error) {
				return model.Namespace{}, nil
			},
			TxBeginTxReturns: func() (transaction.Transaction, error) {
				return mocks.NewTransaction(t), errunknown
			},
			err: errunknown,
		},
		{
			name: `error remove relation`,
			NamespaceGetByIDReturns: func() (model.Namespace, error) {
				return model.Namespace{}, nil
			},
			TxBeginTxReturns: func() (transaction.Transaction, error) {
				tx := mocks.NewTransaction(t)

				rel := mockrepository.NewRelation(t)
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
			NamespaceGetByIDReturns: func() (model.Namespace, error) {
				return model.Namespace{}, nil
			},
			TxBeginTxReturns: func() (transaction.Transaction, error) {
				tx := mocks.NewTransaction(t)

				rel := mockrepository.NewRelation(t)
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
			NamespaceGetByIDReturns: func() (model.Namespace, error) {
				return model.Namespace{}, nil
			},
			TxBeginTxReturns: func() (transaction.Transaction, error) {
				tx := mocks.NewTransaction(t)

				rel := mockrepository.NewRelation(t)
				rel.On(`Delete`, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				tx.On(`Relation`).Return(rel)

				ns := mockrepository.NewNamespace(t)
				ns.On(`DeleteByID`, mock.Anything, mock.Anything).Return(errunknown)
				tx.On(`Namespace`).Return(ns)

				tx.On(`Rollback`, mock.Anything).Return(nil)

				return tx, nil
			},
			RelationDeleteReturns:      nil,
			NamespaceDeleteByIDReturns: errunknown,
			err:                        errunknown,
		},
		{
			name: `error rollback remove namespace`,
			NamespaceGetByIDReturns: func() (model.Namespace, error) {
				return model.Namespace{}, nil
			},
			TxBeginTxReturns: func() (transaction.Transaction, error) {
				tx := mocks.NewTransaction(t)

				rel := mockrepository.NewRelation(t)
				rel.On(`Delete`, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				tx.On(`Relation`).Return(rel)

				ns := mockrepository.NewNamespace(t)
				ns.On(`DeleteByID`, mock.Anything, mock.Anything).Return(errunknown)
				tx.On(`Namespace`).Return(ns)

				tx.On(`Rollback`, mock.Anything).Return(errunknown)

				return tx, nil
			},
			RelationDeleteReturns:      nil,
			NamespaceDeleteByIDReturns: errors.New(`anything`),
			err:                        errunknown,
		},
		{
			name: `error commit`,
			NamespaceGetByIDReturns: func() (model.Namespace, error) {
				return model.Namespace{}, nil
			},
			TxBeginTxReturns: func() (transaction.Transaction, error) {
				tx := mocks.NewTransaction(t)

				rel := mockrepository.NewRelation(t)
				rel.On(`Delete`, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				tx.On(`Relation`).Return(rel)

				ns := mockrepository.NewNamespace(t)
				ns.On(`DeleteByID`, mock.Anything, mock.Anything).Return(nil)
				tx.On(`Namespace`).Return(ns)

				tx.On(`Commit`, mock.Anything).Return(errunknown)

				return tx, nil
			},
			RelationDeleteReturns:      nil,
			NamespaceDeleteByIDReturns: nil,
			err:                        errunknown,
		},
		{
			name: `no error`,
			NamespaceGetByIDReturns: func() (model.Namespace, error) {
				return model.Namespace{}, nil
			},
			TxBeginTxReturns: func() (transaction.Transaction, error) {
				tx := mocks.NewTransaction(t)

				rel := mockrepository.NewRelation(t)
				rel.On(`Delete`, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				tx.On(`Relation`).Return(rel)

				ns := mockrepository.NewNamespace(t)
				ns.On(`DeleteByID`, mock.Anything, mock.Anything).Return(nil)
				tx.On(`Namespace`).Return(ns)

				tx.On(`Commit`, mock.Anything).Return(nil)

				return tx, nil
			},
			RelationDeleteReturns:      nil,
			NamespaceDeleteByIDReturns: nil,
			err:                        nil,
		},
	}

	logger, _ := zap.NewDevelopment()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Initialize service
			s := func() tagservice.Namespace {
				tx := mocks.NewTransactioner(t)
				tx.On(`BeginTx`, mock.Anything, mock.Anything).Return(tt.TxBeginTxReturns()).Maybe()

				namespacerepo := mockrepository.NewNamespace(t)
				namespacerepo.On(`GetByID`, mock.Anything, mock.Anything).
					Return(tt.NamespaceGetByIDReturns()).Maybe()

				relationRepo := mockrepository.NewRelation(t)
				relationRepo.On(`Delete`, mock.Anything, mock.Anything).
					Return(tt.RelationDeleteReturns).Maybe()

				return namespace.New(&namespace.Config{
					Logger:              logger,
					NamespaceRepository: namespacerepo,
					RelationRepository:  relationRepo,
					Transaction:         tx,
				})
			}()

			err := s.Delete(context.TODO(), 88)
			if tt.err == nil {
				assert.NoError(t, err)

				return
			}
			assert.ErrorIs(t, err, tt.err)
		})
	}
}
