package profile

import (
	"time"

	"github.com/nmarsollier/go_cache/utils/memoize"
)

var profileMemoize = memoize.NewSafeMemoize()

// FetchProfile fetch the current profile
func FetchProfile(id string) *Profile {
	return profileMemoize.Value(
		func() *memoize.Memo {
			return memoize.Memoize(fetchProfile(id), 10*time.Minute)
		},
	).(*Profile)
}

func invalidateTSCache() {
	profileMemoize.InvalidateCache()
}
