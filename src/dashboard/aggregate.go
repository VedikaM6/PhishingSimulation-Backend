package dashboard

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserCountsGauge
func AggregateUserCountsData(logsColl *mongo.Collection) (GaugeData, error) {
	gData := GaugeData{
		Type: UserCountsGauge,
	}

	// create the aggregation pipeline to count the attack results by user
	aggPipeline := bson.A{
		bson.D{
			{"$unwind",
				bson.D{
					{"path", "$TargetRecipients"},
					{"preserveNullAndEmptyArrays", false},
				},
			},
		},
		bson.D{
			{"$group",
				bson.D{
					{"_id", "$TargetRecipients.Name"},
					{"TotalAttacks", bson.D{{"$sum", 1}}},
					{"NumAttacksPassed",
						bson.D{
							{"$sum",
								bson.D{
									{"$cond",
										bson.A{
											"$TargetRecipients.IsClicked",
											0,
											1,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	// create a local struct to represent the aggregation results
	type aggResult struct {
		UserName         string `json:"userName" bson:"_id"`
		TotalAttacks     int    `json:"totalAttacks" bson:"TotalAttacks"`
		NumAttacksPassed int    `json:"numAttacksPassed" bson:"NumAttacksPassed"`
	}

	// create a slice to store all the results
	allResults := make([]aggResult, 0)

	// execute the aggregation
	ctx := context.TODO()
	cur, err := logsColl.Aggregate(ctx, aggPipeline)
	if err != nil {
		fmt.Printf("[AggregateUserCountsData] Failed to aggregate data: %+v", err)
		return gData, err
	}

	// make a deferred call to close the Mongo cursor
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		// decode the current document
		var currDoc aggResult
		err := cur.Decode(&currDoc)
		if err != nil {
			fmt.Printf("[AggregateUserCountsData] Failed to decode document: %+v", err)
			continue
		}

		// add the current result to allResults
		allResults = append(allResults, currDoc)
	}

	if cur.Err() != nil {
		fmt.Printf("[AggregateUserCountsData] A Mongo cursor error occurred!: %+v", cur.Err())
	}

	// add 'allResults' to 'gData'
	gData.Data = allResults

	return gData, nil
}

// TeamTotalPerformance
func AggregateTeamTotalPerformanceData(logsColl *mongo.Collection) (GaugeData, error) {
	gData := GaugeData{
		Type: TeamTotalPerformance,
	}

	// create the aggregation pipeline to count the all-time attack results
	aggPipeline := bson.A{
		bson.D{
			{"$unwind",
				bson.D{
					{"path", "$TargetRecipients"},
					{"preserveNullAndEmptyArrays", false},
				},
			},
		},
		bson.D{
			{"$group",
				bson.D{
					{"_id", "RES"},
					{"TotalAttacks", bson.D{{"$sum", 1}}},
					{"NumAttacksPassed",
						bson.D{
							{"$sum",
								bson.D{
									{"$cond",
										bson.A{
											"$TargetRecipients.IsClicked",
											0,
											1,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	// create a local struct to represent the aggregation results
	type aggResult struct {
		Id               string `json:"id" bson:"_id"`
		TotalAttacks     int    `json:"totalAttacks" bson:"TotalAttacks"`
		NumAttacksPassed int    `json:"numAttacksPassed" bson:"NumAttacksPassed"`
	}

	// execute the aggregation
	ctx := context.TODO()
	cur, err := logsColl.Aggregate(ctx, aggPipeline)
	if err != nil {
		fmt.Printf("[AggregateTeamTotalPerformanceData] Failed to aggregate data: %+v", err)
		return gData, err
	}

	// make a deferred call to close the Mongo cursor
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		// decode the current document
		var currDoc aggResult
		err := cur.Decode(&currDoc)
		if err != nil {
			fmt.Printf("[AggregateTeamTotalPerformanceData] Failed to decode document: %+v", err)
			continue
		}

		// set the current result as the gauge data
		gData.Data = currDoc
	}

	if cur.Err() != nil {
		fmt.Printf("[AggregateTeamTotalPerformanceData] A Mongo cursor error occurred!: %+v", cur.Err())
	}

	return gData, nil
}

// EmailCountsGauge
func AggregateEmailCountsData(logsColl *mongo.Collection) (GaugeData, error) {
	gData := GaugeData{
		Type: EmailCountsGauge,
	}

	// create the aggregation pipeline to count the attack results by email
	aggPipeline := bson.A{
		bson.D{
			{"$unwind",
				bson.D{
					{"path", "$TargetRecipients"},
					{"preserveNullAndEmptyArrays", false},
				},
			},
		},
		bson.D{
			{"$group",
				bson.D{
					{"_id", "$UsedEmail.Name"},
					{"TotalAttacks", bson.D{{"$sum", 1}}},
					{"NumAttacksPassed",
						bson.D{
							{"$sum",
								bson.D{
									{"$cond",
										bson.A{
											"$TargetRecipients.IsClicked",
											0,
											1,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	// create a local struct to represent the aggregation results
	type aggResult struct {
		EmailName        string `json:"emailName" bson:"_id"`
		TotalAttacks     int    `json:"totalAttacks" bson:"TotalAttacks"`
		NumAttacksPassed int    `json:"numAttacksPassed" bson:"NumAttacksPassed"`
	}

	// create a slice to store all the results
	allResults := make([]aggResult, 0)

	// execute the aggregation
	ctx := context.TODO()
	cur, err := logsColl.Aggregate(ctx, aggPipeline)
	if err != nil {
		fmt.Printf("[AggregateEmailCountsData] Failed to aggregate data: %+v", err)
		return gData, err
	}

	// make a deferred call to close the Mongo cursor
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		// decode the current document
		var currDoc aggResult
		err := cur.Decode(&currDoc)
		if err != nil {
			fmt.Printf("[AggregateEmailCountsData] Failed to decode document: %+v", err)
			continue
		}

		// add the current result to allResults
		allResults = append(allResults, currDoc)
	}

	if cur.Err() != nil {
		fmt.Printf("[AggregateEmailCountsData] A Mongo cursor error occurred!: %+v", cur.Err())
	}

	// add 'allResults' to 'gData'
	gData.Data = allResults

	return gData, nil
}
