package printer

import (
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Strconv interface {
	FormatBool(b bool) string

	FormatTimestamp(dt *timestamppb.Timestamp) string

	FormatDuration(dt *durationpb.Duration) string
}
