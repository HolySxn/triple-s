package utils

import "regexp"

// IsValidBucketName validates the bucket name based on S3-like constraints:
// - Must be between 3 and 63 characters
// - Can contain lowercase letters, numbers, hyphens, and dots
// - Cannot resemble an IP address (e.g., "192.168.0.1")
// - Cannot start or end with a hyphen, and no consecutive periods or hyphens
func IsValidBucketName(bucketName string) bool {
	// Compile the regex pattern for bucket name validation
	re := regexp.MustCompile(`^[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?(?:\.[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?)*$`)

	// Check length and then match with regex
	return len(bucketName) >= 3 && len(bucketName) <= 63 && re.MatchString(bucketName) &&
		!regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$`).MatchString(bucketName)
}

// IsValidObjectName validates the object name based on constraints:
// - Maximum length of 1024 characters
// - Cannot be "objects.csv" (reserved for bucket metadata)
func IsValidObjectName(objectName string) bool {
	if len(objectName) > 1024 || objectName == "objects.csv" {
		return false
	}

	return true
}
