package main

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Kentik struct {
	ID             primitive.ObjectID `bson:"_id" json:"id"`
	V              int                //`bson:"__v" json:"__v"`
	AgentCity      string             //`json:"agent_city"`
	AgentCountry   string             //`bson:"agentCountry" json:"agent_country"`
	AgentLatitude  float32            //`bson:"agentLatitude" json:"agent_latitude"`
	AgentLongitude float32            //`bson:"agentLongitude" json:"agent_longitude"`
	AgentRegion    string             //`bson:"agentRegion" json:"agent_region"`
	DateStamp      string             //`bson:"date_stamp" json:"date_stamp"`
	LatencyMs      float32            //`bson:"latencyMs" json:"latency_ms"`
	Multiaddr      string             //`bson:"multiaddr" json:"multiaddr"`
	Provider       string             //`bson:"provider" json:"provider"`
	TestId         string             //`bson:"testId" json:"test_id"`
}
