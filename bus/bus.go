// Package bus publish sự kiện lên EventBridge bus của khu vực.
package bus

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
)

func Publish(ctx context.Context, c *eventbridge.Client, busName, source, detailType string, detail any) error {
	b, err := json.Marshal(detail)
	if err != nil {
		return err
	}
	out, err := c.PutEvents(ctx, &eventbridge.PutEventsInput{
		Entries: []types.PutEventsRequestEntry{{
			EventBusName: aws.String(busName),
			Source:       aws.String(source),
			DetailType:   aws.String(detailType),
			Detail:       aws.String(string(b)),
		}},
	})
	if err != nil {
		return fmt.Errorf("put %s: %w", detailType, err)
	}
	if out.FailedEntryCount > 0 {
		return fmt.Errorf("put %s: %s", detailType, aws.ToString(out.Entries[0].ErrorMessage))
	}
	return nil
}
