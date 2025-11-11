package external

import "encoding/json"

type Event struct {
	ID              string        `json:"id"`
	Object          string        `json:"object"` // "event"
	APIVersion      string        `json:"api_version"`
	Created         int64         `json:"created"`
	Data            EventData     `json:"data"`
	Livemode        bool          `json:"livemode"`
	PendingWebhooks int           `json:"pending_webhooks"`
	Request         *EventRequest `json:"request,omitempty"`
	Type            string        `json:"type"` // e.g., "setup_intent.created"
}

type EventData struct {
	Object             SetupIntent    `json:"object"`
	PreviousAttributes map[string]any `json:"previous_attributes,omitempty"`
}

type EventRequest struct {
	ID             *string `json:"id"`
	IdempotencyKey *string `json:"idempotency_key"`
}

// ---- SetupIntent and nested options ----

type SetupIntent struct {
	ID                      string                `json:"id"`
	Object                  string                `json:"object"` // "setup_intent"
	Application             *string               `json:"application"`
	AutomaticPaymentMethods any                   `json:"automatic_payment_methods"` // keep flexible
	CancellationReason      *string               `json:"cancellation_reason"`
	ClientSecret            *string               `json:"client_secret"`
	Created                 int64                 `json:"created"`
	Customer                *string               `json:"customer"`
	Description             *string               `json:"description"`
	FlowDirections          *[]string             `json:"flow_directions"`
	LastSetupError          json.RawMessage       `json:"last_setup_error"` // nullable, unknown shape
	LatestAttempt           *string               `json:"latest_attempt"`
	Livemode                bool                  `json:"livemode"`
	Mandate                 *string               `json:"mandate"`
	Metadata                map[string]string     `json:"metadata"`
	NextAction              json.RawMessage       `json:"next_action"` // nullable, variant shape
	OnBehalfOf              *string               `json:"on_behalf_of"`
	PaymentMethod           *string               `json:"payment_method"`
	PaymentMethodOptions    *PaymentMethodOptions `json:"payment_method_options,omitempty"`
	PaymentMethodTypes      []string              `json:"payment_method_types"`
	SingleUseMandate        *string               `json:"single_use_mandate"`
	Status                  string                `json:"status"`
	Usage                   string                `json:"usage"`
}

type PaymentMethodOptions struct {
	ACSSDebit *ACSSDebitOptions `json:"acss_debit,omitempty"`
}

type ACSSDebitOptions struct {
	Currency           string              `json:"currency"`
	MandateOptions     *ACSSMandateOptions `json:"mandate_options,omitempty"`
	VerificationMethod string              `json:"verification_method"`
}

type ACSSMandateOptions struct {
	IntervalDescription string `json:"interval_description"`
	PaymentSchedule     string `json:"payment_schedule"`
	TransactionType     string `json:"transaction_type"`
}
