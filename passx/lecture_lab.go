package passx

// A single atomic lecture/lab combination, where |lecture| must be taken with
// |lab|. |lecture| or |lab|, but not both, may be null.
type LectureLab struct {
	lecture *Section
	lab     *Section
}
