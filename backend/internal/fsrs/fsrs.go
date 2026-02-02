package fsrs

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"
)

// State represents the learning state of a card.
type State int

const (
	Learning   State = 1
	Review     State = 2
	Relearning State = 3
)

// Rating represents a user's review rating.
type Rating int

const (
	Again Rating = 1
	Hard  Rating = 2
	Good  Rating = 3
	Easy  Rating = 4
)

const (
	StabilityMin      = 0.001
	MinDifficulty     = 1.0
	MaxDifficulty     = 10.0
	InitialStabMax    = 100.0
	FSRSDefaultDecay  = 0.1542
	NumParams         = 21
	DefaultMaxIvl     = 36500
	DefaultRetention  = 0.9
)

// DefaultParameters are the FSRS v6 default model weights.
var DefaultParameters = [NumParams]float64{
	0.212, 1.2931, 2.3065, 8.2956,
	6.4133, 0.8334, 3.0194, 0.001,
	1.8722, 0.1666, 0.796, 1.4835,
	0.0614, 0.2629, 1.6483, 0.6014,
	1.8729, 0.5425, 0.0912, 0.0658,
	FSRSDefaultDecay,
}

var lowerBounds = [NumParams]float64{
	StabilityMin, StabilityMin, StabilityMin, StabilityMin,
	1.0, 0.001, 0.001, 0.001,
	0.0, 0.0, 0.001, 0.001,
	0.001, 0.001, 0.0, 0.0,
	1.0, 0.0, 0.0, 0.0,
	0.1,
}

var upperBounds = [NumParams]float64{
	InitialStabMax, InitialStabMax, InitialStabMax, InitialStabMax,
	10.0, 4.0, 4.0, 0.75,
	4.5, 0.8, 3.5, 5.0,
	0.25, 0.9, 4.0, 1.0,
	6.0, 2.0, 2.0, 0.8,
	0.8,
}

type fuzzRange struct {
	start, end, factor float64
}

var fuzzRanges = []fuzzRange{
	{2.5, 7.0, 0.15},
	{7.0, 20.0, 0.1},
	{20.0, math.Inf(1), 0.05},
}

// Card is the in-memory scheduling state. Not a DB model.
type Card struct {
	State      State
	Step       *int
	Stability  *float64
	Difficulty *float64
	Due        time.Time
	LastReview *time.Time
}

// ReviewOutput contains the updated card and metadata.
type ReviewOutput struct {
	Card           Card
	ScheduledDays  int
	ElapsedDays    int
	Retrievability float64
}

// ReviewEntry is a minimal review record for RescheduleCard.
type ReviewEntry struct {
	Rating         Rating
	ReviewDatetime time.Time
}

// Scheduler holds FSRS parameters and configuration.
type Scheduler struct {
	Params           [NumParams]float64
	DesiredRetention float64
	LearningSteps    []time.Duration
	RelearningSteps  []time.Duration
	MaximumInterval  int
	EnableFuzzing    bool

	rng    *rand.Rand
	decay  float64
	factor float64
}

// Option configures a Scheduler.
type Option func(*Scheduler)

func WithParams(p [NumParams]float64) Option {
	return func(s *Scheduler) { s.Params = p }
}

func WithDesiredRetention(r float64) Option {
	return func(s *Scheduler) { s.DesiredRetention = r }
}

func WithLearningSteps(steps []time.Duration) Option {
	return func(s *Scheduler) { s.LearningSteps = steps }
}

func WithRelearningSteps(steps []time.Duration) Option {
	return func(s *Scheduler) { s.RelearningSteps = steps }
}

func WithMaximumInterval(days int) Option {
	return func(s *Scheduler) { s.MaximumInterval = days }
}

func WithEnableFuzzing(enable bool) Option {
	return func(s *Scheduler) { s.EnableFuzzing = enable }
}

func WithRNG(rng *rand.Rand) Option {
	return func(s *Scheduler) { s.rng = rng }
}

// NewScheduler creates a Scheduler with defaults matching py-fsrs.
func NewScheduler(opts ...Option) (*Scheduler, error) {
	s := &Scheduler{
		Params:           DefaultParameters,
		DesiredRetention: DefaultRetention,
		LearningSteps:    []time.Duration{1 * time.Minute, 10 * time.Minute},
		RelearningSteps:  []time.Duration{10 * time.Minute},
		MaximumInterval:  DefaultMaxIvl,
		EnableFuzzing:    true,
	}
	for _, opt := range opts {
		opt(s)
	}
	if err := validateParams(s.Params); err != nil {
		return nil, err
	}
	s.decay = -s.Params[20]
	s.factor = math.Pow(0.9, 1.0/s.decay) - 1
	if s.rng == nil {
		s.rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	return s, nil
}

func validateParams(p [NumParams]float64) error {
	var errs []string
	for i := 0; i < NumParams; i++ {
		if p[i] < lowerBounds[i] || p[i] > upperBounds[i] {
			errs = append(errs, fmt.Sprintf("parameters[%d] = %g is out of bounds: (%g, %g)", i, p[i], lowerBounds[i], upperBounds[i]))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("one or more parameters are out of bounds:\n%s", strings.Join(errs, "\n"))
	}
	return nil
}

// Retrievability returns the current probability of recall.
func (s *Scheduler) Retrievability(card Card, now time.Time) float64 {
	if card.LastReview == nil || card.Stability == nil {
		return 0
	}
	elapsedDays := int(now.Sub(*card.LastReview).Hours() / 24)
	if elapsedDays < 0 {
		elapsedDays = 0
	}
	return math.Pow(1+s.factor*float64(elapsedDays)/(*card.Stability), s.decay)
}

// ReviewCard is the core scheduling function.
func (s *Scheduler) ReviewCard(card Card, rating Rating, now time.Time) ReviewOutput {
	var daysSinceLastReview *int
	if card.LastReview != nil {
		d := int(now.Sub(*card.LastReview).Hours() / 24)
		daysSinceLastReview = &d
	}

	var nextInterval time.Duration

	switch card.State {
	case Learning:
		// update stability and difficulty
		if card.Stability == nil || card.Difficulty == nil {
			stab := s.initialStability(rating)
			diff := s.initialDifficulty(rating, true)
			card.Stability = &stab
			card.Difficulty = &diff
		} else if daysSinceLastReview != nil && *daysSinceLastReview < 1 {
			stab := s.shortTermStability(*card.Stability, rating)
			diff := s.nextDifficulty(*card.Difficulty, rating)
			card.Stability = &stab
			card.Difficulty = &diff
		} else {
			r := s.Retrievability(card, now)
			stab := s.nextStability(*card.Difficulty, *card.Stability, r, rating)
			diff := s.nextDifficulty(*card.Difficulty, rating)
			card.Stability = &stab
			card.Difficulty = &diff
		}

		// calculate next interval
		if len(s.LearningSteps) == 0 || (*card.Step >= len(s.LearningSteps) && (rating == Hard || rating == Good || rating == Easy)) {
			card.State = Review
			card.Step = nil
			days := s.nextInterval(*card.Stability)
			nextInterval = time.Duration(days) * 24 * time.Hour
		} else {
			switch rating {
			case Again:
				zero := 0
				card.Step = &zero
				nextInterval = s.LearningSteps[*card.Step]
			case Hard:
				step := *card.Step
				if step == 0 && len(s.LearningSteps) == 1 {
					nextInterval = s.LearningSteps[0] * 3 / 2
				} else if step == 0 && len(s.LearningSteps) >= 2 {
					nextInterval = (s.LearningSteps[0] + s.LearningSteps[1]) / 2
				} else {
					nextInterval = s.LearningSteps[step]
				}
			case Good:
				step := *card.Step
				if step+1 == len(s.LearningSteps) {
					card.State = Review
					card.Step = nil
					days := s.nextInterval(*card.Stability)
					nextInterval = time.Duration(days) * 24 * time.Hour
				} else {
					newStep := step + 1
					card.Step = &newStep
					nextInterval = s.LearningSteps[newStep]
				}
			case Easy:
				card.State = Review
				card.Step = nil
				days := s.nextInterval(*card.Stability)
				nextInterval = time.Duration(days) * 24 * time.Hour
			}
		}

	case Review:
		// update stability and difficulty
		if daysSinceLastReview != nil && *daysSinceLastReview < 1 {
			stab := s.shortTermStability(*card.Stability, rating)
			card.Stability = &stab
		} else {
			r := s.Retrievability(card, now)
			stab := s.nextStability(*card.Difficulty, *card.Stability, r, rating)
			card.Stability = &stab
		}
		diff := s.nextDifficulty(*card.Difficulty, rating)
		card.Difficulty = &diff

		// calculate next interval
		switch rating {
		case Again:
			if len(s.RelearningSteps) == 0 {
				days := s.nextInterval(*card.Stability)
				nextInterval = time.Duration(days) * 24 * time.Hour
			} else {
				card.State = Relearning
				zero := 0
				card.Step = &zero
				nextInterval = s.RelearningSteps[0]
			}
		default: // Hard, Good, Easy
			days := s.nextInterval(*card.Stability)
			nextInterval = time.Duration(days) * 24 * time.Hour
		}

	case Relearning:
		// update stability and difficulty
		if daysSinceLastReview != nil && *daysSinceLastReview < 1 {
			stab := s.shortTermStability(*card.Stability, rating)
			diff := s.nextDifficulty(*card.Difficulty, rating)
			card.Stability = &stab
			card.Difficulty = &diff
		} else {
			r := s.Retrievability(card, now)
			stab := s.nextStability(*card.Difficulty, *card.Stability, r, rating)
			diff := s.nextDifficulty(*card.Difficulty, rating)
			card.Stability = &stab
			card.Difficulty = &diff
		}

		// calculate next interval
		if len(s.RelearningSteps) == 0 || (*card.Step >= len(s.RelearningSteps) && (rating == Hard || rating == Good || rating == Easy)) {
			card.State = Review
			card.Step = nil
			days := s.nextInterval(*card.Stability)
			nextInterval = time.Duration(days) * 24 * time.Hour
		} else {
			switch rating {
			case Again:
				zero := 0
				card.Step = &zero
				nextInterval = s.RelearningSteps[0]
			case Hard:
				step := *card.Step
				if step == 0 && len(s.RelearningSteps) == 1 {
					nextInterval = s.RelearningSteps[0] * 3 / 2
				} else if step == 0 && len(s.RelearningSteps) >= 2 {
					nextInterval = (s.RelearningSteps[0] + s.RelearningSteps[1]) / 2
				} else {
					nextInterval = s.RelearningSteps[step]
				}
			case Good:
				step := *card.Step
				if step+1 == len(s.RelearningSteps) {
					card.State = Review
					card.Step = nil
					days := s.nextInterval(*card.Stability)
					nextInterval = time.Duration(days) * 24 * time.Hour
				} else {
					newStep := step + 1
					card.Step = &newStep
					nextInterval = s.RelearningSteps[newStep]
				}
			case Easy:
				card.State = Review
				card.Step = nil
				days := s.nextInterval(*card.Stability)
				nextInterval = time.Duration(days) * 24 * time.Hour
			}
		}
	}

	if s.EnableFuzzing && card.State == Review {
		nextInterval = s.getFuzzedInterval(nextInterval)
	}

	due := now.Add(nextInterval)
	card.Due = due
	card.LastReview = &now

	scheduledDays := int(nextInterval.Hours() / 24)
	elapsedDays := 0
	if daysSinceLastReview != nil {
		elapsedDays = *daysSinceLastReview
	}

	return ReviewOutput{
		Card:           card,
		ScheduledDays:  scheduledDays,
		ElapsedDays:    elapsedDays,
		Retrievability: s.Retrievability(card, now),
	}
}

// RescheduleCard replays review history with the current scheduler parameters.
func (s *Scheduler) RescheduleCard(card Card, reviews []ReviewEntry) Card {
	rescheduled := Card{
		State: Learning,
		Step:  intPtr(0),
		Due:   card.Due,
	}
	for _, rev := range reviews {
		out := s.ReviewCard(rescheduled, rev.Rating, rev.ReviewDatetime)
		rescheduled = out.Card
	}
	return rescheduled
}

// --- internal helpers ---

func (s *Scheduler) initialStability(rating Rating) float64 {
	stab := s.Params[int(rating)-1]
	return clampStability(stab)
}

func (s *Scheduler) initialDifficulty(rating Rating, clamp bool) float64 {
	diff := s.Params[4] - math.Exp(s.Params[5]*float64(rating-1)) + 1
	if clamp {
		diff = clampDifficulty(diff)
	}
	return diff
}

func (s *Scheduler) shortTermStability(stability float64, rating Rating) float64 {
	increase := math.Exp(s.Params[17]*(float64(rating)-3+s.Params[18])) * math.Pow(stability, -s.Params[19])
	if rating == Good || rating == Easy {
		if increase < 1.0 {
			increase = 1.0
		}
	}
	stab := stability * increase
	return clampStability(stab)
}

func (s *Scheduler) nextDifficulty(difficulty float64, rating Rating) float64 {
	arg1 := s.initialDifficulty(Easy, false)
	deltaDiff := -(s.Params[6] * (float64(rating) - 3))
	linearDamped := (10.0 - difficulty) * deltaDiff / 9.0
	arg2 := difficulty + linearDamped
	next := s.Params[7]*arg1 + (1-s.Params[7])*arg2
	return clampDifficulty(next)
}

func (s *Scheduler) nextStability(difficulty, stability, retrievability float64, rating Rating) float64 {
	var next float64
	if rating == Again {
		next = s.nextForgetStability(difficulty, stability, retrievability)
	} else {
		next = s.nextRecallStability(difficulty, stability, retrievability, rating)
	}
	return clampStability(next)
}

func (s *Scheduler) nextForgetStability(difficulty, stability, retrievability float64) float64 {
	longTerm := s.Params[11] *
		math.Pow(difficulty, -s.Params[12]) *
		(math.Pow(stability+1, s.Params[13]) - 1) *
		math.Exp((1-retrievability)*s.Params[14])

	shortTerm := stability / math.Exp(s.Params[17]*s.Params[18])

	return math.Min(longTerm, shortTerm)
}

func (s *Scheduler) nextRecallStability(difficulty, stability, retrievability float64, rating Rating) float64 {
	hardPenalty := 1.0
	if rating == Hard {
		hardPenalty = s.Params[15]
	}
	easyBonus := 1.0
	if rating == Easy {
		easyBonus = s.Params[16]
	}
	return stability * (1 +
		math.Exp(s.Params[8])*
			(11-difficulty)*
			math.Pow(stability, -s.Params[9])*
			(math.Exp((1-retrievability)*s.Params[10])-1)*
			hardPenalty*
			easyBonus)
}

func (s *Scheduler) nextInterval(stability float64) int {
	ivl := (stability / s.factor) * (math.Pow(s.DesiredRetention, 1.0/s.decay) - 1)
	ivl = math.Round(ivl)
	if ivl < 1 {
		ivl = 1
	}
	if ivl > float64(s.MaximumInterval) {
		ivl = float64(s.MaximumInterval)
	}
	return int(ivl)
}

func (s *Scheduler) getFuzzedInterval(interval time.Duration) time.Duration {
	intervalDays := int(interval.Hours() / 24)
	if float64(intervalDays) < 2.5 {
		return interval
	}

	delta := 1.0
	for _, fr := range fuzzRanges {
		end := fr.end
		if float64(intervalDays) < end {
			end = float64(intervalDays)
		}
		contrib := fr.factor * math.Max(end-fr.start, 0.0)
		delta += contrib
	}

	minIvl := int(math.Round(float64(intervalDays) - delta))
	maxIvl := int(math.Round(float64(intervalDays) + delta))

	if minIvl < 2 {
		minIvl = 2
	}
	if maxIvl > s.MaximumInterval {
		maxIvl = s.MaximumInterval
	}
	if minIvl > maxIvl {
		minIvl = maxIvl
	}

	fuzzed := s.rng.Float64()*float64(maxIvl-minIvl+1) + float64(minIvl)
	fuzzedDays := int(math.Round(fuzzed))
	if fuzzedDays > s.MaximumInterval {
		fuzzedDays = s.MaximumInterval
	}

	return time.Duration(fuzzedDays) * 24 * time.Hour
}

func clampDifficulty(d float64) float64 {
	if d < MinDifficulty {
		return MinDifficulty
	}
	if d > MaxDifficulty {
		return MaxDifficulty
	}
	return d
}

func clampStability(s float64) float64 {
	if s < StabilityMin {
		return StabilityMin
	}
	return s
}

func intPtr(v int) *int {
	return &v
}
