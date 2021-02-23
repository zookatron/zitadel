package query

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore"
	usr_repo "github.com/caos/zitadel/internal/repository/user"

	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/crypto"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/id"
	iam_repo "github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

type QuerySide struct {
	iamID        string
	eventstore   *eventstore.Eventstore
	idGenerator  id.Generator
	secretCrypto crypto.Crypto
}

type Config struct {
	Eventstore     *eventstore.Eventstore
	SystemDefaults sd.SystemDefaults
}

func StartQuerySide(config *Config) (repo *QuerySide, err error) {
	repo = &QuerySide{
		iamID:       config.SystemDefaults.IamID,
		eventstore:  config.Eventstore,
		idGenerator: id.SonyFlakeGenerator,
	}
	iam_repo.RegisterEventMappers(repo.eventstore)
	usr_repo.RegisterEventMappers(repo.eventstore)

	repo.secretCrypto, err = crypto.NewAESCrypto(config.SystemDefaults.IDPConfigVerificationKey)
	if err != nil {
		return nil, err
	}
	return repo, nil
}

func (r *QuerySide) IAMByID(ctx context.Context, id string) (_ *iam_model.IAM, err error) {
	readModel, err := r.iamByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return readModelToIAM(readModel), nil
}

func (r *QuerySide) iamByID(ctx context.Context, id string) (_ *ReadModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	readModel := NewReadModel(id)
	err = r.eventstore.FilterToQueryReducer(ctx, readModel)
	if err != nil {
		return nil, err
	}

	return readModel, nil
}