package projection

import (
	"context"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/repository/features"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/repository/org"
)

type FeatureProjection struct {
	crdb.StatementHandler
}

const (
	FeatureTable = "zitadel.projections.features"
)

func NewFeatureProjection(ctx context.Context, config crdb.StatementHandlerConfig) *FeatureProjection {
	p := &FeatureProjection{}
	config.ProjectionName = FeatureTable
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *FeatureProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: org.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  org.FeaturesSetEventType,
					Reduce: p.reduceFeatureSet,
				},
				{
					Event:  org.FeaturesRemovedEventType,
					Reduce: p.reduceFeatureRemoved,
				},
			},
		},
		{
			Aggregate: iam.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  iam.FeaturesSetEventType,
					Reduce: p.reduceFeatureSet,
				},
			},
		},
	}
}

const (
	FeatureAggregateIDCol              = "aggregate_id"
	FeatureCreationDateCol             = "creation_date"
	FeatureChangeDateCol               = "change_date"
	FeatureSequenceCol                 = "sequence"
	FeatureIsDefaultCol                = "is_default"
	FeatureTierNameCol                 = "tier_name"
	FeatureTierDescriptionCol          = "tier_description"
	FeatureStateCol                    = "state"
	FeatureStateDescriptionCol         = "state_description"
	FeatureAuditLogRetentionCol        = "audit_log_retention"
	FeatureLoginPolicyFactorsCol       = "login_policy_factors"
	FeatureLoginPolicyIDPCol           = "login_policy_idp"
	FeatureLoginPolicyPasswordlessCol  = "login_policy_passwordless"
	FeatureLoginPolicyRegistrationCol  = "login_policy_registration"
	FeatureLoginPolicyUsernameLoginCol = "login_policy_username_login"
	FeatureLoginPolicyPasswordResetCol = "login_policy_password_reset"
	FeaturePasswordComplexityPolicyCol = "password_complexity_policy"
	FeatureLabelPolicyPrivateLabelCol  = "label_policy_private_label"
	FeatureLabelPolicyWatermarkCol     = "label_policy_watermark"
	FeatureCustomDomainCol             = "custom_domain"
	FeaturePrivacyPolicyCol            = "privacy_policy"
	FeatureMetadataUserCol             = "metadata_user"
	FeatureCustomTextMessageCol        = "custom_text_message"
	FeatureCustomTextLoginCol          = "custom_text_login"
	FeatureLockoutPolicyCol            = "lockout_policy"
	FeatureActionsCol                  = "actions"
)

func (p *FeatureProjection) reduceFeatureSet(event eventstore.EventReader) (*handler.Statement, error) {
	var featureEvent features.FeaturesSetEvent
	var isDefault bool
	switch e := event.(type) {
	case *iam.FeaturesSetEvent:
		featureEvent = e.FeaturesSetEvent
		isDefault = true
	case *org.FeaturesSetEvent:
		featureEvent = e.FeaturesSetEvent
		isDefault = false
	default:
		logging.LogWithFields("HANDL-M9ets", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.FeaturesSetEventType, iam.FeaturesSetEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-K0erf", "reduce.wrong.event.type")
	}

	cols := []handler.Column{
		handler.NewCol(FeatureAggregateIDCol, featureEvent.Aggregate().ID),
		handler.NewCol(FeatureCreationDateCol, featureEvent.CreationDate()),
		handler.NewCol(FeatureChangeDateCol, featureEvent.CreationDate()),
		handler.NewCol(FeatureSequenceCol, featureEvent.Sequence()),
		handler.NewCol(FeatureIsDefaultCol, isDefault),
	}
	if featureEvent.TierName != nil {
		cols = append(cols, handler.NewCol(FeatureTierNameCol, *featureEvent.TierName))
	}
	if featureEvent.TierDescription != nil {
		cols = append(cols, handler.NewCol(FeatureTierDescriptionCol, *featureEvent.TierDescription))
	}
	if featureEvent.State != nil {
		cols = append(cols, handler.NewCol(FeatureStateCol, *featureEvent.State))
	}
	if featureEvent.StateDescription != nil {
		cols = append(cols, handler.NewCol(FeatureStateDescriptionCol, *featureEvent.StateDescription))
	}
	if featureEvent.AuditLogRetention != nil {
		cols = append(cols, handler.NewCol(FeatureAuditLogRetentionCol, *featureEvent.AuditLogRetention))
	}
	if featureEvent.LoginPolicyFactors != nil {
		cols = append(cols, handler.NewCol(FeatureLoginPolicyFactorsCol, *featureEvent.LoginPolicyFactors))
	}
	if featureEvent.LoginPolicyIDP != nil {
		cols = append(cols, handler.NewCol(FeatureLoginPolicyIDPCol, *featureEvent.LoginPolicyIDP))
	}
	if featureEvent.LoginPolicyPasswordless != nil {
		cols = append(cols, handler.NewCol(FeatureLoginPolicyPasswordlessCol, *featureEvent.LoginPolicyPasswordless))
	}
	if featureEvent.LoginPolicyRegistration != nil {
		cols = append(cols, handler.NewCol(FeatureLoginPolicyRegistrationCol, *featureEvent.LoginPolicyRegistration))
	}
	if featureEvent.LoginPolicyUsernameLogin != nil {
		cols = append(cols, handler.NewCol(FeatureLoginPolicyUsernameLoginCol, *featureEvent.LoginPolicyUsernameLogin))
	}
	if featureEvent.LoginPolicyPasswordReset != nil {
		cols = append(cols, handler.NewCol(FeatureLoginPolicyPasswordResetCol, *featureEvent.LoginPolicyPasswordReset))
	}
	if featureEvent.PasswordComplexityPolicy != nil {
		cols = append(cols, handler.NewCol(FeaturePasswordComplexityPolicyCol, *featureEvent.PasswordComplexityPolicy))
	}
	if featureEvent.LabelPolicyPrivateLabel != nil {
		cols = append(cols, handler.NewCol(FeatureLabelPolicyPrivateLabelCol, *featureEvent.LabelPolicyPrivateLabel))
	}
	if featureEvent.LabelPolicyWatermark != nil {
		cols = append(cols, handler.NewCol(FeatureLabelPolicyWatermarkCol, *featureEvent.LabelPolicyWatermark))
	}
	if featureEvent.CustomDomain != nil {
		cols = append(cols, handler.NewCol(FeatureCustomDomainCol, *featureEvent.CustomDomain))
	}
	if featureEvent.PrivacyPolicy != nil {
		cols = append(cols, handler.NewCol(FeaturePrivacyPolicyCol, *featureEvent.PrivacyPolicy))
	}
	if featureEvent.MetadataUser != nil {
		cols = append(cols, handler.NewCol(FeatureMetadataUserCol, *featureEvent.MetadataUser))
	}
	if featureEvent.CustomTextMessage != nil {
		cols = append(cols, handler.NewCol(FeatureCustomTextMessageCol, *featureEvent.CustomTextMessage))
	}
	if featureEvent.CustomTextLogin != nil {
		cols = append(cols, handler.NewCol(FeatureCustomTextLoginCol, *featureEvent.CustomTextLogin))
	}
	if featureEvent.LockoutPolicy != nil {
		cols = append(cols, handler.NewCol(FeatureLockoutPolicyCol, *featureEvent.LockoutPolicy))
	}
	if featureEvent.Actions != nil {
		cols = append(cols, handler.NewCol(FeatureActionsCol, *featureEvent.Actions))
	}
	return crdb.NewUpsertStatement(
		&featureEvent,
		cols), nil
}

func (p *FeatureProjection) reduceFeatureRemoved(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*org.FeaturesRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-fN903", "seq", event.Sequence(), "expectedType", org.FeaturesRemovedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-0p4rf", "reduce.wrong.event.type")
	}
	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(FeatureAggregateIDCol, e.Aggregate().ID),
		},
	), nil
}