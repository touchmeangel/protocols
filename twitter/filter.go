package twitter

type FollowerFilter func(followers int) bool

func MinFollowers(min int) FollowerFilter {
	return func(followers int) bool { return followers >= min }
}

func MaxFollowers(max int) FollowerFilter {
	return func(followers int) bool { return followers <= max }
}

func MatchesAllFollowers(followers int, filters []FollowerFilter) bool {
	for _, f := range filters {
		if !f(followers) {
			return false
		}
	}
	return true
}
