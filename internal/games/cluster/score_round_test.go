package cluster

import (
	"math"
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

	dots, _, _, distances, winners, outliers := scoreRound(submissions, 0.5, 0.5, uuidByte(1).String())
	if len(dots) != 3 {
		t.Fatalf("expected 3 dots, got %d", len(dots))
	}

	if got := distances[uuidByte(1).String()]; got != 0 {
		t.Fatalf("expected current player distance 0, got %v", got)
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

func TestScoreRoundTwoPlayersExposeViewerDistancesWithoutOutliers(t *testing.T) {
	hostID := uuidByte(1)
	guestID := uuidByte(2)
	submissions := []submissionRecord{
		{
			PlayerID: hostID,
			Nickname: "Host",
			X:        0.10,
			Y:        0.20,
		},
		{
			PlayerID: guestID,
			Nickname: "Guest",
			X:        0.85,
			Y:        0.65,
		},
	}

	_, roundPoints, _, distances, winners, outliers := scoreRound(submissions, 0.475, 0.425, hostID.String())

	if roundPoints[hostID.String()] != roundPoints[guestID.String()] {
		t.Fatalf("expected two-player centroid scoring tie, got %d vs %d", roundPoints[hostID.String()], roundPoints[guestID.String()])
	}

	if got := distances[guestID.String()]; math.Abs(got-0.87) > 0.01 {
		t.Fatalf("expected guest distance about 0.87, got %.4f", got)
	}

	if len(winners) != 2 {
		t.Fatalf("expected tied winners for two-player round, got %v", winners)
	}

	if len(outliers) != 0 {
		t.Fatalf("expected no outlier summary when everyone ties, got %v", outliers)
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
