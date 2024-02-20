package main

import (
	"context"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/data-preservation-programs/spade-tenant/config"
	"github.com/data-preservation-programs/spade-tenant/db"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

type SP struct {
	Version        int       `bson:"__v"`
	Provider       string    `bson:"_id"`
	AgentCity      string    `bson:"agentCity"`
	AgentCountry   string    `bson:"agentCountry"`
	AgentLatitude  float64   `bson:"agentLatitude"`
	AgentLongitude float64   `bson:"agentLongitude"`
	AgentRegion    string    `bson:"agentRegion"`
	Date           time.Time `bson:"date"`
	LatencyMs      float64   `bson:"latencyMs"`
	Multiaddr      string    `bson:"multiaddr"`
	TestId         string    `bson:"testId"`
}

func main() {
	godotenv.Load()
	ctx := context.TODO()
	sps := fetchSpData(ctx)

	config := config.InitConfig()
	database, err := db.OpenDatabase(config.DB_URL, config.DEBUG, config.DRY_RUN)

	if err != nil {
		panic(err)
	}

	var clauses []db.TenantSPEligibilityClauses
	database.Model(&db.TenantSPEligibilityClauses{}).Find(&clauses)
	tenantClauseMap := createTenantClauseMap(clauses)

	for tenantId, clauses := range tenantClauseMap {
		for _, sp := range sps {
			isEligible := verifyEligibility(clauses, sp)
			updateDatabase(isEligible, database, sp, tenantId)
		}
	}
}

func fetchSpData(ctx context.Context) []SP {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(os.Getenv("MONGODB_URI")).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	collection := client.Database(os.Getenv("MONGODB_DATABASE")).Collection(os.Getenv("MONGODB_SP_COLLECTION"))

	cur, _ := collection.Find(context.TODO(), bson.D{{}})
	var sps []SP = []SP{}
	for cur.Next(context.TODO()) {
		var sp SP
		err := cur.Decode(&sp)
		if err != nil {
			log.Fatal(err)
		}

		sps = append(sps, sp)
	}

	return sps
}

func verifyEligibility(clauses []db.TenantSPEligibilityClauses, sp SP) bool {
	for _, clause := range clauses {
		if !compare(clause.ClauseValue, getField(&sp, clause.ClauseAttribute), clause.ClauseOperator) {
			return false
		}
	}
	return true
}

func createTenantClauseMap(clauses []db.TenantSPEligibilityClauses) map[db.ID][]db.TenantSPEligibilityClauses {
	tenantClauseMap := make(map[db.ID][]db.TenantSPEligibilityClauses)
	for _, clause := range clauses {
		if _, ok := tenantClauseMap[clause.TenantID]; !ok {
			tenantClauseMap[clause.TenantID] = make([]db.TenantSPEligibilityClauses, 0)
		}

		tenantClauseMap[clause.TenantID] = append(tenantClauseMap[clause.TenantID], clause)
	}

	return tenantClauseMap
}

func updateDatabase(isValid bool, database *gorm.DB, sp SP, tenantId db.ID) {
	// Find the number part of the provider. ex: f0123 -> 123
	spID, err := strconv.Atoi(sp.Provider[2:])
	if err != nil {
		log.Println("Could not parse ", sp.Provider, " to an integer. Skipping this sp.")
		return
	}

	var tsp db.TenantsSPs
	tsp.TenantID = tenantId
	tsp.SPID = db.ID(spID)

	rowsAffected := database.Find(&tsp).RowsAffected
	if isValid {
		if rowsAffected == 0 {
			database.Create(&db.SP{SPID: db.ID(spID)})
			database.Create(&db.TenantsSPs{TenantID: tenantId, SPID: db.ID(spID), TenantSpState: db.TenantSpStateEligible})
		} else {
			tsp.TenantSpState = db.TenantSpStateEligible
			database.Updates(&tsp)
		}
	} else {
		if rowsAffected != 0 {
			tsp.TenantSpState = db.TenantSpStateDisabled
			database.Updates(&tsp)
		}
	}
}

func getField(sp *SP, field string) string {
	r := reflect.ValueOf(sp)
	f := reflect.Indirect(r).FieldByName(field)
	return f.String()
}

func compare(tenantSetting []string, spSetting string, comparison db.ComparisonOperator) bool {
	switch comparison {
	case db.IncludedIn:
		return contains(tenantSetting, spSetting)
	case db.ExcludedFrom:
		return !contains(tenantSetting, spSetting)
	default:
		return false
	}
}

func contains(tenantSetting []string, spSetting string) bool {
	for _, setting := range tenantSetting {
		if strings.EqualFold(strings.ToLower(setting), strings.ToLower(spSetting)) {
			return true
		}
	}

	return false
}
