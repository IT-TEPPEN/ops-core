package errors

type ErrorLayer string
type ErrorFeature string

const (
	LayerDomain         ErrorLayer = "D"
	LayerApplication    ErrorLayer = "A"
	LayerInfrastructure ErrorLayer = "I"
	LayerPresentation   ErrorLayer = "P"
)

const (
	FeatureAuthentication ErrorFeature = "AT"
)
