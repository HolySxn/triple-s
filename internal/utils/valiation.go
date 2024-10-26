package utils

import "regexp"

func IsValidBucketName(bucketName string) bool {
	// Compile the regex pattern for bucket name validation
	re := regexp.MustCompile(`^[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?(?:\.[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?)*$`)

	// Check length and then match with regex
	return len(bucketName) >= 3 && len(bucketName) <= 63 && re.MatchString(bucketName) &&
		!regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$`).MatchString(bucketName)
}

func IsValidObjectName(objectName string) bool{
	if len(objectName) > 1024 || objectName == "objects.csv" {
		return false
	}

	return true
}