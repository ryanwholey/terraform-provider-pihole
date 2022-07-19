package pihole

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateCNAMERecord(t *testing.T) {
	acceptance(t)

	t.Parallel()

	ctx := context.Background()

	clients := map[string]*Client{
		"password": newAccTestClient(t, ctx, Config{
			Password: os.Getenv("PIHOLE_PASSWORD"),
			URL:      os.Getenv("PIHOLE_URL"),
		}),
		"apiToken": newAccTestClient(t, ctx, Config{
			URL:      os.Getenv("PIHOLE_URL"),
			APIToken: doubleHash256(os.Getenv("PIHOLE_PASSWORD")),
		}),
	}

	for name, c := range clients {
		c := c

		t.Run(fmt.Sprintf("%s client", name), func(t *testing.T) {
			t.Parallel()

			record := &CNAMERecord{
				Domain: fmt.Sprintf("test-%s.com", randomSuffix()),
				Target: fmt.Sprintf("test-%s.com", randomSuffix()),
			}

			t.Cleanup(func() {
				require.NoError(t, c.DeleteCNAMERecord(ctx, record.Domain))
			})

			_, err := c.CreateCNAMERecord(ctx, record)
			require.NoError(t, err)

			cname, err := c.GetCNAMERecord(ctx, record.Domain)
			require.NoError(t, err)

			assert.Equal(t, strings.ToLower(record.Domain), cname.Domain)
			assert.Equal(t, strings.ToLower(record.Target), cname.Target)
		})
	}

	// c := newAccTestClient(t, ctx, Config{
	// 	Password: os.Getenv("PIHOLE_PASSWORD"),
	// 	URL:      os.Getenv("PIHOLE_URL"),
	// })

}
