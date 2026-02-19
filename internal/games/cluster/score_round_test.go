package cluster

import (
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
)

func TestScoreRoundReturnsWinnersAndOutliers(t *testing.T) {
	submissions := []submissionRecord{
		{
			PlayerID: uuidByte(1),
			Nickname: "Near",
			X:        0.52,
			Y:        0.48,
		},
		{
			PlayerID: uuidByte(2),
			Nickname: "Mid",
			X:        0.70,
			Y:        0.50,
		},
		{
			PlayerID: uuidByte(3),
			Nickname: "Far",
			X:        0.98,
			Y:        0.98,
		},
	}

	dots, _, winners, outliers := scoreRound(submissions, 0.5, 0.5, uuidByte(1).String())
	if len(dots) != 3 {
		t.Fatalf("expected 3 dots, got %d", len(dots))
	}

	if len(winners) != 1 || winners[0] != "Near" {
		t.Fatalf("expected Near as winner, got %v", winners)
	}

	if len(outliers) != 1 || outliers[0] != "Far" {
		t.Fatalf("expected Far as outlier, got %v", outliers)
	}

	if dots[0].AnimationDelay != 0 || dots[1].AnimationDelay != 80 || dots[2].AnimationDelay != 160 {
		t.Fatalf("unexpected animation delays: %d, %d, %d", dots[0].AnimationDelay, dots[1].AnimationDelay, dots[2].AnimationDelay)
	}
}

func uuidByte(n byte) pgtype.UUID {
	var b [16]byte
	b[15] = n
	return pgtype.UUID{
		Bytes: b,
		Valid: true,
	}
}
