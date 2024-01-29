package utils

import (
	"regexp"
	"strconv"
	"strings"
)

func safeParseInt(value string) (int, bool) {
	if match, _ := regexp.MatchString(`^\d+$`, value); match {
		num, _ := strconv.Atoi(value)
		return num, true
	}
	return 0, false
}

func isWildcard(value string) bool {
	return value == "*"
}

func isQuestionMark(value string) bool {
	return value == "?"
}

func isInRange(value int, start int, stop int) bool {
	return value >= start && value <= stop
}

func isValidRange(value string, start int, stop int) bool {
	sides := strings.Split(value, "-")
	switch len(sides) {
	case 1:
		num, ok := safeParseInt(value)
		return isWildcard(value) || (ok && isInRange(num, start, stop))
	case 2:
		small, ok1 := safeParseInt(sides[0])
		big, ok2 := safeParseInt(sides[1])
		return ok1 && ok2 && small <= big && isInRange(small, start, stop) && isInRange(big, start, stop)
	default:
		return false
	}
}

func isValidStep(value string) bool {
	if value == "" {
		return true
	}
	if match, _ := regexp.MatchString(`[^\d]`, value); match {
		return false
	}
	num, _ := safeParseInt(value)
	return num > 0
}

func validateForRange(value string, start int, stop int) bool {
	if match, _ := regexp.MatchString(`[^\d-,\/*]`, value); match {
		return false
	}
	list := strings.Split(value, ",")
	for _, condition := range list {
		splits := strings.Split(condition, "/")
		if strings.HasSuffix(strings.TrimSpace(condition), "/") {
			return false
		}
		if len(splits) > 2 {
			return false
		}
		left := splits[0]
		var right string
		if len(splits) > 1 {
			right = splits[1]
		}
		if !isValidRange(left, start, stop) || !isValidStep(right) {
			return false
		}
	}
	return true
}

func hasValidSeconds(seconds string) bool {
	return validateForRange(seconds, 0, 59)
}

func hasValidMinutes(minutes string) bool {
	return validateForRange(minutes, 0, 59)
}

func hasValidHours(hours string) bool {
	return validateForRange(hours, 0, 23)
}

func hasValidDays(days string, allowBlankDay bool) bool {
	return (allowBlankDay && isQuestionMark(days)) || validateForRange(days, 1, 31)
}

var monthAlias = map[string]string{
	"jan": "1",
	"feb": "2",
	"mar": "3",
	"apr": "4",
	"may": "5",
	"jun": "6",
	"jul": "7",
	"aug": "8",
	"sep": "9",
	"oct": "10",
	"nov": "11",
	"dec": "12",
}

func hasValidMonths(months string, alias bool) bool {
	if match, _ := regexp.MatchString(`\/[a-zA-Z]`, months); match {
		return false
	}
	if alias {
		for k, v := range monthAlias {
			months = strings.ReplaceAll(months, k, v)
		}
	}
	return validateForRange(months, 1, 12)
}

var weekdaysAlias = map[string]string{
	"sun": "0",
	"mon": "1",
	"tue": "2",
	"wed": "3",
	"thu": "4",
	"fri": "5",
	"sat": "6",
}

func hasValidWeekdays(weekdays string, alias bool, allowBlankDay bool, allowSevenAsSunday bool) bool {
	if allowBlankDay && isQuestionMark(weekdays) {
		return true
	} else if !allowBlankDay && isQuestionMark(weekdays) {
		return false
	}
	if match, _ := regexp.MatchString(`\/[a-zA-Z]`, weekdays); match {
		return false
	}
	if alias {
		for k, v := range weekdaysAlias {
			weekdays = strings.ReplaceAll(weekdays, k, v)
		}
	}
	return validateForRange(weekdays, 0, func() int {
		if allowSevenAsSunday {
			return 7
		}
		return 6
	}())
}

func hasCompatibleDayFormat(days string, weekdays string, allowBlankDay bool) bool {
	return !(allowBlankDay && isQuestionMark(days) && isQuestionMark(weekdays))
}

func split(cron string) []string {
	return strings.Fields(cron)
}

type Options struct {
	alias              bool
	seconds            bool
	allowBlankDay      bool
	allowSevenAsSunday bool
}

var options = Options{
	alias:              false,
	seconds:            true,
	allowBlankDay:      true,
	allowSevenAsSunday: false,
}

func IsValidAwsCron(cron string) bool {
	splits := split(cron)

	if options.seconds {
		if len(splits) > 6 || len(splits) < 6 {
			return false
		}
	} else {
		if len(splits) > 5 || len(splits) < 5 {
			return false
		}
	}

	checks := make([]bool, 0)
	if len(splits) == 6 {
		seconds := splits[0]
		checks = append(checks, hasValidSeconds(seconds))
	}

	minutes := splits[0]
	hours := splits[1]
	days := splits[2]
	months := splits[3]
	weekdays := splits[4]
	checks = append(checks, hasValidMinutes(minutes))
	checks = append(checks, hasValidHours(hours))
	checks = append(checks, hasValidDays(days, options.allowBlankDay))
	checks = append(checks, hasValidMonths(months, options.alias))
	checks = append(checks, hasValidWeekdays(weekdays, options.alias, options.allowBlankDay, options.allowSevenAsSunday))
	checks = append(checks, hasCompatibleDayFormat(days, weekdays, options.allowBlankDay))
	for _, check := range checks {
		if !check {
			return false
		}
	}
	return true
}
