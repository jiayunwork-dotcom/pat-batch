package spc

// ZoneClassification categorizes a point's position relative to control chart
// zones. Zone A: between 2σ and 3σ; Zone B: between 1σ and 2σ; Zone C: within 1σ.
type ZoneClassification struct {
	Index int
	Zone  string // "A+", "A-", "B+", "B-", "C+", "C-", "OOC" (out of control)
}

// ClassifyZones assigns each point to its control chart zone relative to CL and sigma.
func ClassifyZones(values []float64, cl, sigma float64) []ZoneClassification {
	if sigma <= 0 {
		return nil
	}
	result := make([]ZoneClassification, len(values))
	for i, v := range values {
		diff := v - cl
		absDiff := diff
		if absDiff < 0 {
			absDiff = -absDiff
		}
		sign := "+"
		if diff < 0 {
			sign = "-"
		}

		var zone string
		switch {
		case absDiff > 3*sigma:
			zone = "OOC"
		case absDiff > 2*sigma:
			zone = "A" + sign
		case absDiff > sigma:
			zone = "B" + sign
		default:
			zone = "C" + sign
		}
		result[i] = ZoneClassification{Index: i, Zone: zone}
	}
	return result
}

// ZoneRuleCheck applies the Western Electric zone rules:
// 1. One point in Zone A or beyond (= Nelson rule 1)
// 2. Two of three consecutive in Zone A on same side (= Nelson rule 5)
// 3. Four of five consecutive in Zone B or beyond on same side (= Nelson rule 6)
// 4. Eight consecutive points on one side of CL (similar to Nelson rule 2 with 8)
// Returns indices of violations per rule.
func ZoneRuleCheck(zones []ZoneClassification) map[int][]int {
	violations := map[int][]int{
		1: nil,
		2: nil,
		3: nil,
		4: nil,
	}

	// Rule 1: any OOC point.
	for _, z := range zones {
		if z.Zone == "OOC" {
			violations[1] = append(violations[1], z.Index)
		}
	}

	// Rule 2: two of three in Zone A on same side.
	for i := 0; i < len(zones)-2; i++ {
		aPlus := 0
		aMinus := 0
		for j := 0; j < 3; j++ {
			z := zones[i+j].Zone
			if z == "A+" || z == "OOC" {
				aPlus++
			}
			if z == "A-" {
				aMinus++
			}
		}
		if aPlus >= 2 || aMinus >= 2 {
			violations[2] = append(violations[2], i)
		}
	}

	// Rule 3: four of five in Zone B or beyond on same side.
	for i := 0; i < len(zones)-4; i++ {
		bPlus := 0
		bMinus := 0
		for j := 0; j < 5; j++ {
			z := zones[i+j].Zone
			if z == "A+" || z == "B+" || z == "OOC" {
				bPlus++
			}
			if z == "A-" || z == "B-" {
				bMinus++
			}
		}
		if bPlus >= 4 || bMinus >= 4 {
			violations[3] = append(violations[3], i)
		}
	}

	// Rule 4: eight consecutive on same side.
	for i := 0; i < len(zones)-7; i++ {
		positive := 0
		negative := 0
		for j := 0; j < 8; j++ {
			z := zones[i+j].Zone
			if z == "C+" || z == "B+" || z == "A+" || z == "OOC" {
				positive++
			} else {
				negative++
			}
		}
		if positive == 8 || negative == 8 {
			violations[4] = append(violations[4], i)
		}
	}

	return violations
}

// ProcessStabilityScore computes a 0-100 stability score based on zone rule
// violations. 100 means no violations; each violation type reduces the score.
func ProcessStabilityScore(values []float64, cl, sigma float64) int {
	zones := ClassifyZones(values, cl, sigma)
	violations := ZoneRuleCheck(zones)

	score := 100
	// Each rule violation type reduces score.
	if len(violations[1]) > 0 {
		score -= 30 // OOC is severe
	}
	if len(violations[2]) > 0 {
		score -= 20
	}
	if len(violations[3]) > 0 {
		score -= 15
	}
	if len(violations[4]) > 0 {
		score -= 15
	}
	if score < 0 {
		score = 0
	}
	return score
}
