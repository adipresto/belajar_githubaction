package models

type ActivityType string

const (
	walking    ActivityType = "Walking"
	yoga       ActivityType = "Yoga"
	stretching ActivityType = "Stretching"
	running    ActivityType = "Running"
	cycling    ActivityType = "Cycling"
	swimming   ActivityType = "Swimming"
	dancing    ActivityType = "Dancing"
	hiking     ActivityType = "Hiking"
	hiit       ActivityType = "HIIT"
	jumprope   ActivityType = "JumpRope"
)

type Activity struct {
	ID             int          `db:"id"`
	UserId         int          `db:"user_id"`
	ActivityType   ActivityType `db:"activity_type"`
	DoneAt         string       `db:"done_at"`
	DurationInMin  int          `db:"duration_in_minutes"`
	CaloriesBurned int          `db:"calories_burned"`
	CreatedAt      string       `db:"created_at"`
	UpdatedAt      string       `db:"updated_at"`
	IsActive       bool         `db:"is_active"`
}

func (a *ActivityType) GetTotalCalories(durationInMin int) int {
	var calories map[ActivityType]int = map[ActivityType]int{
		walking:    4,
		yoga:       4,
		stretching: 4,
		cycling:    8,
		swimming:   8,
		dancing:    8,
		hiking:     10,
		running:    10,
		hiit:       10,
		jumprope:   10,
	}

	return calories[*a] * durationInMin
}
