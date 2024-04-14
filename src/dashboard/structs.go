package dashboard

type GaugeType string

const (
	UserCountsGauge          = "UserCounts"
	TeamTotalPerformance     = "TotalTeamPerformance"
	EmailCountsGauge         = "EmailCounts"
	TeamPerformanceLastWeek  = "TeamPerformanceLastWeek"
	ScheduledAttacksNextWeek = "ScheduledAttacksNextWeek"
	ScheduledAttacksForUsers = "ScheduledAttacksForUsers"
)

// Represents data for all gauge types
type GaugeData struct {
	Type GaugeType   `json:"type" bson:"Type"`
	Data interface{} `json:"data" bson:"Data"`
}
