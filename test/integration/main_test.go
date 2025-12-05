//go:build integration

package integration

import (
	"context"
	"os"
	"testing"

	"github.com/LackOfMorals/unofficialMcp/test/integration/dbservice"
)

var dbs = dbservice.NewDBService()

func TestMain(m *testing.M) {
	ctx := context.Background()

	dbs.Start(ctx)

	code := m.Run()

	dbs.Stop(ctx)

	os.Exit(code)
}
