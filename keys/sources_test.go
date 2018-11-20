package keys_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/LUSHDigital/microservice-core-golang/keys"
)

func TestSources(t *testing.T) {
	// Only run these tests to manually check that the correct keys are fetched
	// The results are printed to STDOUT
	if os.Getenv("TEST_KEY_SOURCES") != "" {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
		defer cancel()

		func(ctx context.Context) {
			fmt.Print(">>>>> BEGIN STAGING <<<<<\n\n")
			defer fmt.Print(">>>>> END STAGING <<<<<\n\n")
			bts, err := keys.StagingTokenPublicKeySources.GetKey(ctx)
			if err != nil {
				t.Fatal(err)
			}
			if len(bts) < 1 {
				t.Fatal(bts)
			}

			fmt.Println(string(bts))
		}(ctx)

		func(ctx context.Context) {
			fmt.Print(">>>>> BEGIN PRODUCTION <<<<<\n\n")
			defer fmt.Print(">>>>> END PRODUCTION <<<<<\n\n")
			bts, err := keys.ProductionTokenPublicKeySources.GetKey(ctx)
			if err != nil {
				t.Fatal(err)
			}
			if len(bts) < 1 {
				t.Fatal(bts)
			}

			fmt.Println(string(bts))
		}(ctx)
	}
}
