// This file is subject to the terms and conditions defined in
// file 'LICENSE.txt', which is part of this source code package.
package printer

import (
	"fmt"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"math"
	"strconv"
	"strings"
	"time"
)

const dateTimeLayout = "2006-Jan-02"

var PrettyStrconv Strconv = &struct {
	boolEmoji
	durationHuman
	timestampDate
}{}

var SimpleStrconv Strconv = &struct {
	simpleTimestamp
	simpleDuration
	simpleBool
}{}

type simpleTimestamp struct{}

func (p *simpleTimestamp) FormatTimestamp(dt *timestamppb.Timestamp) string {
	if dt == nil || (dt.Nanos == 0 && dt.Seconds == 0) {
		return ""
	}
	return dt.AsTime().Format(dateTimeLayout)
}

type simpleBool struct{}

func (p *simpleBool) FormatBool(b bool) string {
	return strconv.FormatBool(b)
}

type simpleDuration struct{}

func (p *simpleDuration) FormatDuration(dt *durationpb.Duration) string {
	if dt == nil || (dt.Nanos == 0 && dt.Seconds == 0) {
		return ""
	}
	return dt.AsDuration().String()
}

type boolEmoji struct{}

func (p *boolEmoji) FormatBool(b bool) string {
	if b {
		return ":heavy_check_mark:"
	}
	return "-"
}

type timestampDate struct{}

func (p *timestampDate) FormatTimestamp(dt *timestamppb.Timestamp) string {
	if dt == nil || (dt.Nanos == 0 && dt.Seconds == 0) {
		return "-"
	}
	return dt.AsTime().Format(dateTimeLayout)
}

type durationHuman struct{}

func (p *durationHuman) FormatDuration(duration *durationpb.Duration) string {
	if duration == nil || (duration.Nanos == 0 && duration.Seconds == 0) {
		return "-"
	}
	return humanizeDuration(duration.AsDuration())
}

func humanizeDuration(duration time.Duration) string {
	days := int64(duration.Hours() / 24)
	hours := int64(math.Mod(duration.Hours(), 24))
	minutes := int64(math.Mod(duration.Minutes(), 60))
	seconds := int64(math.Mod(duration.Seconds(), 60))

	chunks := []struct {
		singularName string
		amount       int64
	}{
		{"day", days},
		{"hour", hours},
		{"minute", minutes},
		{"second", seconds},
	}

	var parts []string

	for _, chunk := range chunks {
		switch chunk.amount {
		case 0:
			continue
		case 1:
			parts = append(parts, fmt.Sprintf("%d %s", chunk.amount, chunk.singularName))
		default:
			parts = append(parts, fmt.Sprintf("%d %ss", chunk.amount, chunk.singularName))
		}
	}

	return strings.Join(parts, " ")
}
