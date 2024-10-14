package bucket

import "regexp"

func isValidBucketName(bucketName string) bool {
	// Compile the regex pattern for bucket name validation
	re := regexp.MustCompile(`^(?!-)(?!.*--)(?!.*\.\.)(?!.*\.$)(?!.*\-$)([a-z0-9][a-z0-9.-]{1,61}[a-z0-9])?$`)

	// Check length and then match with regex
	return len(bucketName) >= 3 && len(bucketName) <= 63 && re.MatchString(bucketName) &&
		!regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$`).MatchString(bucketName)
}
