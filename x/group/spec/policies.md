# Decision Policies

```go
type ThresholdDecisionPolicy struct {
	// Specifies the number of votes that must be accumulated in order for a decision to be made by the group.
	// A member gets as many votes as is indicated by their Weight field.
	// A big integer is used here to avoid any potential vulnerabilities from overflow errors
	// where large weight and threshold values are used.
	DecisionThreshold sdk.Int `json:"decision_threshold"`
}

type PercentageDecisionPolicy struct {
	Percent sdk.Dec `json:"percent"`
}
```