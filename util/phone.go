package util

import (
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

var (
	//match1    = regexp.MustCompile(`^(((1(\s|))|)\([1-9]{3}\)(\s|-|)[1-9]{3}(\s|-|)[1-9]{4})$`)
	//match1          = regexp.MustCompile(`^(?:(?:\+?1\s*(?:[.-]\s*)?)?(?:\(\s*([2-9]1[02-9]|[2-9][02-8]1|[2-9][02-8][02-9])\s*\)|([2-9]1[02-9]|[2-9][02-8]1|[2-9][02-8][02-9]))\s*(?:[.-]\s*)?)?([2-9]1[02-9]|[2-9][02-9]1|[2-9][02-9]{2})\s*(?:[.-]\s*)?([0-9]{4})(?:\s*(?:#|x\.?|ext\.?|extension)\s*(\d+))?$`)
	//match2          = regexp.MustCompile(`^(((1(\s)|)|)[1-9]{3}(\s|-|)[1-9]{3}(\s|-|)[1-9]{4})$`)
	chinaphonematch = regexp.MustCompile(`^(?:\+?86)?1(?:3\d{3}|5[^4\D]\d{2}|8\d{3}|7(?:[0-35-9]\d{2}|4(?:0\d|1[0-2]|9\d))|9[0-35-9]\d{2}|6[2567]\d{2}|4(?:(?:10|4[01])\d{3}|[68]\d{4}|[579]\d{2}))\d{6}$`)
	formatnum       = regexp.MustCompile(`[0-9]`)
)

// telephoneCheck("555-555-5555") 应该返回一个布尔值.
// telephoneCheck("1 555-555-5555") 应该返回 true.
// telephoneCheck("1 (555) 555-5555") 应该返回 true.
// telephoneCheck("5555555555") 应该返回 true.
// telephoneCheck("555-555-5555") 应该返回 true.
// telephoneCheck("(555)555-5555") 应该返回 true.
// telephoneCheck("1(555)555-5555") 应该返回 true.
// telephoneCheck("1 555)555-5555") 应该返回 false.
// telephoneCheck("1 555 555 5555") 应该返回 true.
// telephoneCheck("1 456 789 4444") 应该返回 true.
// telephoneCheck("123**&!!asdf#") 应该返回 false.
// telephoneCheck("55555555") 应该返回 false.
// telephoneCheck("(6505552368)") 应该返回 false
// telephoneCheck("2 (757) 622-7382") 应该返回 false.
// telephoneCheck("0 (757) 622-7382") 应该返回 false.
// telephoneCheck("-1 (757) 622-7382") 应该返回 false
// telephoneCheck("2 757 622-7382") 应该返回 false.
// telephoneCheck("10 (757) 622-7382") 应该返回 false.
// telephoneCheck("27576227382") 应该返回 false.
// telephoneCheck("(275)76227382") 应该返回 false.
// telephoneCheck("2(757)6227382") 应该返回 false.
// telephoneCheck("2(757)622-7382") 应该返回 false.
// telephoneCheck("555)-555-5555") 应该返回 false.
// telephoneCheck("(555-555-5555") 应该返回 false.

// ValidatePhone validates an email address based on a regular expression
func ValidatePhone(phone string) bool {
	return chinaphonematch.MatchString(phone) && len(phone) == 11
}

func FormatPhoneNum(phone string) string {
	params := formatnum.FindAllString(phone, -1)
	return Join(params, "")
}

func GenValidateCode(width int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}

func GenInviteCode(width int) string {
	numeric := [62]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%c", numeric[rand.Intn(r)])
	}
	return sb.String()
}

const (
	// Replace AccessKeyID with your AccessKeyID key.
	AccessKeyID = "AKIA6CVAISMROUCQM7S6"

	// Replace AccessKeyID with your AccessKeyID key.
	SecretAccessKey = "jz2OmgLQxK3KwOculG44yy26lWaJgdHqCmbJLfhp"

	// Replace us-west-2 with the AWS Region you're using for Amazon SNS.
	AwsRegion = "us-east-1"
)

func SendSMS(phoneNumber string, message string) error {
	// Create Session and assign AccessKeyID and SecretAccessKey
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(AwsRegion),
		Credentials: credentials.NewStaticCredentials(AccessKeyID, SecretAccessKey, ""),
	},
	)

	// Create SNS service
	svc := sns.New(sess)

	// Pass the phone number and message.
	params := &sns.PublishInput{
		PhoneNumber: aws.String(phoneNumber),
		Message:     aws.String(message),
	}

	// sends a text message (SMS message) directly to a phone number.
	resp, err := svc.Publish(params)

	if err != nil {
		log.Println(err.Error())
	}

	fmt.Println(resp) // print the response data.

	return nil
}
