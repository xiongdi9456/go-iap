package playstore

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"google.golang.org/api/androidpublisher/v3"
	"google.golang.org/appengine/urlfetch"
	"reflect"
	"testing"
)

var base64JsonKey = "yourBase64Key="

var base64dummyKey = "ew0KICAidHlwZSI6ICJzZXJ2aWNlX2FjY291bnQiLA0KICAicHJvamVjdF9pZCI6ICJnby1pYXAiLA0KICAicHJpdmF0ZV9rZXlfaWQiOiAiZHVtbXkiLA0KICAicHJpdmF0ZV9rZXkiOiAiLS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tXG5NSUlFdmdJQkFEQU5CZ2txaGtpRzl3MEJBUUVGQUFTQ0JLZ3dnZ1NrQWdFQUFvSUJBUUNvbG53V0ZoWnI5ZkFDXG4zanBKM2xpOUkreGc1MmRpd0dKc1JwaGNFNEtMYnpSaTNWVnVBLzFNRDJHMDJjTmJGeVdzQ0loOGgzLzRNbGZsXG5rTGtvcFZSZTVRZXM0TGNxOEQ0VTlHbm42cFl3blpSOTVQZStmek1XSzRiaUloL3R4TTMyZXpsb0FsL2lOUW4vXG5nUEN5VHNKRzZVZTcybHRScm5RV3ZDNC9uL0MzRy95S2NmS2h6Njg2ejU3eDg2ZmtvMlJxK2ttUWhTZmN0RGN6XG54N2JTZlA1OEhXVDlYQzR5R0hEamZpMDJsZ1VNWGluRWtlUFcwdVAvUmRvMmtHZWloUldSRjBRTUR5ZXdqdUtKXG5rWi96QjhocmZvUVRHZjIxOC9vZUF0dFFpdGZMZFZZK2o3OWNlcUE5QVF6K2xqMUI1K0lWYkY5TkNraGkzWGE4XG53RzJpelhFYkFnTUJBQUVDZ2dFQUxWZ0VXZ0JvMlhMVnNqL0pWN0xwRlQ1RFJyRVd1cFhhSXhzOXdZNHh6NFVDXG5oeERXK0hjME9xL3NiTE1oZXkrWG4xVFFPUVpNNGhuUVVGdURvYTRPS2xQWm82THhRU2hLMm1IKzFqVGZYb1lUXG51V1RMU2I1MnBDRGk3NUdVR3VTVExSZHBrbE1KTFJOczgvN2ZQbVkybE5JTHpEZm4xZWxoS2ZoRlRER2RrZklTXG5KazY0L3dCU0VMeWUxUDhXWDNCS1d4WXRQSmo4L3NKakw5dElYV3IrWko0dmxoMzNXNW9uQVkvZTBSaGZGQ0oyXG4wRVpGRVV5K25TVFRwVEdSZmM2MjBZdnNiclVqM3BaN3QwSEROSDd6L251RHlQeFFKTkU4NUVZUEQwWnlPVTZ6XG5ldmlvTm5iQSsxVURHRTlYOTRVVW5yTXJCeTVOSW5Nd2wrd3Y2eVVhSVFLQmdRRGJxYWs1V200ZWtRSWVabDJoXG5Gd2VPU0dBTnpPWVZhN01hMUJvYVBwSys1bmpNcW1EMURweHBpRWJhRjZ6cFVjMmZFSE9GRGdtc0FGa2FMNXJjXG50c1QvOGprV0R3SkowL0JROEJaYnpGblNSOGx4ZWxQbHBTcklXK3pnKzVnYWw2bU16S0tidGtvU1VLMXpubjJ6XG5BMnZoU04rOTdVbGEwcitnd2QrbUNIL01Td0tCZ1FERWVlT2xNb0czNnlsZUVxTmg3YTVJcklrdHRZQm5ybjFxXG5YQkhSZlhUOXN0Tjh4Ty9KWENrNUtSSERJM1FmSW1XK0pnSVdzcEY2NUNNK2ZxUTdZdlFjcFdERnQ2b3NYUTZuXG5lYWQvUWJSeU85S3F2N0hKNmhycnBia0tDb0owQWxUSmZreE9rL0wwZzN2ZG5YYmtuNi96VXgvSWV0bGRxOE9TXG4zZVNhQW5wTWNRS0JnUUN3YkFkNkJPTkVzWHBlS0NFeTdHZ3BJbi9qRlpvRndrWkxXZWJOQlV5ZS9rUXZQUGc2XG5XYzNPQktIRE1CaTBHL3Rsc2JUV1BId1FKUWRyUEtqSWRCS3M3a0pqTVJMSmNPc21WbTNldExXL2FlQ2t2M2I2XG5qamxhU2xwcUtDZjEwN0ZkWUUySmVsTHJldGlVYjhyTktBWlJIbEoxSEVzNkl1RzluM2ljeFY2L0dRS0JnSEhJXG5yVCtFaW44NjMxQXR0eFVGa3dOZmVHcFNUTFMrNXI3cjV4M05iQzFvblBZTEQxc3IxbXZXRHdWVnlQQW0rWWt2XG5kZEl6US9GSm9lZVZiQU5BZ1dMOW01ZWxrQlgxSm9GekFML0FDNEtFaHJLQUpiUnJzWDk3RURoeWNhNUJrMXpGXG5tZWQvNHhvYjgyWWF4VG9PQ05YLzg4NGs1ekZLUWc4U0ZrdmkxM1RoQW9HQkFMUVdZK1dIQy9IQ1N1YWlUNXJOXG5xYnJSMDNNQVU4OE91MnpWS2dRS2tLdHg4WjdES2NnbkZjRFpLUDlHT3FyRVY3YWRJcXQ2SHJML0dRanVRQ0hhXG5UaHN3RkZlWTNNOTBFYnpqRkxNVllYZWh5RWMyYnZCZnBaRUhuVDNVUjVyNWVkWFpqRmJyckhXQzBhbGJzUkdaXG52aENnQU93OG5OeEhUbTJQazJGdmx2clhcbi0tLS0tRU5EIFBSSVZBVEUgS0VZLS0tLS1cbiIsDQogICJjbGllbnRfaWQiOiAiZHVtbXkiLA0KICAiYXV0aF91cmkiOiAiaHR0cHM6Ly9hY2NvdW50cy5nb29nbGUuY29tL28vb2F1dGgyL2F1dGgiLA0KICAidG9rZW5fdXJpIjogImh0dHBzOi8vYWNjb3VudHMuZ29vZ2xlLmNvbS9vL29hdXRoMi90b2tlbiIsDQogICJhdXRoX3Byb3ZpZGVyX3g1MDlfY2VydF91cmwiOiAiaHR0cHM6Ly93d3cuZ29vZ2xlYXBpcy5jb20vb2F1dGgyL3YxL2NlcnRzIiwNCiAgImNsaWVudF94NTA5X2NlcnRfdXJsIjogImh0dHBzOi8vd3d3Lmdvb2dsZWFwaXMuY29tL3JvYm90L3YxL21ldGFkYXRhL3g1MDkvZ28taWFwJTQwZ28taWFwLmlhbS5nc2VydmljZWFjY291bnQuY29tIg0KfQ=="

var jsonKey []byte
var dummyKey []byte

func init() {
	f, err := base64.StdEncoding.DecodeString(base64JsonKey)
	if err != nil {
		panic(err)
	}
	jsonKey = f
	d, err := base64.StdEncoding.DecodeString(base64dummyKey)
	if err != nil {
		panic(err)
	}
	dummyKey = d
}

func TestNew(t *testing.T) {
	t.Parallel()

	// Exception scenario
	expected := "oauth2: cannot fetch token: 400 Bad Request\nResponse: {\"error\":\"invalid_grant\",\"error_description\":\"Invalid grant: account not found\"}"

	_, err := New(dummyKey)
	if err == nil || err.Error() != expected {
		t.Errorf("got %v\nwant %v", err, expected)
	}

	_, actual := New(nil)
	if actual == nil || actual.Error() != "unexpected end of JSON input" {
		t.Errorf("got %v\nwant %v", actual, expected)
	}

	_, err = New(jsonKey)
	if err != nil {
		t.Errorf("got %#v", err)
	}
}

func TestNewWithClient(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	httpClient := urlfetch.Client(ctx)

	_, err := NewWithClient(dummyKey, httpClient)
	if err != nil {
		t.Errorf("transport should be urlfetch's one")
	}
}

func TestNewWithClientErrors(t *testing.T) {
	t.Parallel()
	expected := errors.New("client is nil")

	_, actual := NewWithClient(dummyKey, nil)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}

	ctx := context.Background()
	httpClient := urlfetch.Client(ctx)

	_, actual = NewWithClient(nil, httpClient)
	if actual == nil || actual.Error() != "unexpected end of JSON input" {
		t.Errorf("got %v\nwant %v", actual, expected)
	}

}

func TestAcknowledgeSubscription(t *testing.T) {
	t.Parallel()
	// Exception scenario
	expected := "googleapi: Error 400: Invalid Value, invalid"

	client, _ := New(jsonKey)
	ctx := context.Background()
	req := &androidpublisher.SubscriptionPurchasesAcknowledgeRequest{
		DeveloperPayload: "user001",
	}
	err := client.AcknowledgeSubscription(ctx, "package", "subscriptionID", "purchaseToken", req)

	if err == nil || err.Error() != expected {
		t.Errorf("got %v\nwant %v", err, expected)
	}

	// TODO Normal scenario
}

func TestVerifySubscription(t *testing.T) {
	t.Parallel()
	// Exception scenario
	expected := "googleapi: Error 400: Invalid Value, invalid"

	client, _ := New(jsonKey)
	ctx := context.Background()
	_, err := client.VerifySubscription(ctx, "package", "subscriptionID", "purchaseToken")

	if err == nil || err.Error() != expected {
		t.Errorf("got %v\nwant %v", err, expected)
	}

	// TODO Normal scenario
}

func TestVerifyProduct(t *testing.T) {
	t.Parallel()
	// Exception scenario
	expected := "googleapi: Error 400: Invalid Value, invalid"

	client, _ := New(jsonKey)
	ctx := context.Background()
	purchaseInfo, err := client.VerifyProduct(ctx, "com.aplusjapan.sdkdevelop", "com.aplusjapan.sdkdevelop.item1", "geaileejakhefdcedddielii.AO-J1Oxw51cb1Mn8qN51RXYgK_-v6KNAHaPQMZvwB64SETNnx3dUyg7KOD6gE531g7PAW898dvug3xCU1EUtaxLPSJ-ZUHEoDnMdK-80ENWt4G8mStgtYC0")
	fmt.Printf("orderId=%s", purchaseInfo.OrderId)
	if err == nil || err.Error() != expected {
		t.Errorf("got %v", err)
	}

	// TODO Normal scenario
}

func TestVoidedPurchase(t *testing.T) {
	t.Parallel()
	// Exception scenario
	//expected := "googleapi: Error 400: Invalid Value, invalid"

	//秒
	//fmt.Println(time.Now().UTC().Unix())
	////微秒 milliseconds
	//fmt.Println(time.Now().UTC().UnixNano()/1e6)
	//fmt.Println(time.Now().UTC().UnixNano())

	client, _ := New(jsonKey)
	ctx := context.Background()
	//startTime := time.Now().UTC().Add(-24 * 20 * time.Hour).UnixNano()/1e6
	//endTime := time.Now().UTC().UnixNano()/1e6
	//fmt.Printf("startTime=%v, endTime=%v", startTime, endTime)
	pageToken := ""
	for {
		pagingVoidedPurchases, nextPageToken, err := client.GetVoidedPurchase(ctx, "com.xxxx.xxxx", 0, 0, 2, 0, pageToken)
		if err != nil {
			fmt.Printf("GetVoidedPurchase, err=%#v", err)
			//isOneAppOverLoop = true
			break
		}
		if len(pagingVoidedPurchases) > 0 {
			fmt.Println("找到了退款的订单")
		} else {
			fmt.Println("没有退款的订单")
		}

		for _, voidedPurchase := range pagingVoidedPurchases {
			fmt.Println("==============")
			fmt.Printf("voidedPurchase=%#v", voidedPurchase)
			fmt.Println("==============")
		}

		if nextPageToken == "" {
			fmt.Println("nextPageToken------")
			fmt.Println(nextPageToken)
			fmt.Println("nextPageToken为空,结束循环")
			break
		} else {
			fmt.Println("nextPageToken不为空:" + nextPageToken)
			pageToken = nextPageToken
		}

	}
	//	purchaseInfo, err := client.GetVoidedPurchase(ctx, "com.aplusjapan.sdkdevelop", "", )
	//fmt.Printf("orderId=%s", purchaseInfo.OrderId)
	//if err == nil || err.Error() != expected {
	//	t.Errorf("got %v", err)
	//}
	//
	//// TODO Normal scenario
}

func TestAcknowledgeProduct(t *testing.T) {
	t.Parallel()
	// Exception scenario
	expected := "googleapi: Error 400: Invalid Value, invalid"

	client, _ := New(jsonKey)
	ctx := context.Background()
	err := client.AcknowledgeProduct(ctx, "package", "productID", "purchaseToken", "")

	if err == nil || err.Error() != expected {
		t.Errorf("got %v", err)
	}

	// TODO Normal scenario
}

func TestCancelSubscription(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	client, _ := New(jsonKey)
	expectedStr := "googleapi: Error 400: Invalid Value, invalid"
	actual := client.CancelSubscription(ctx, "package", "productID", "purchaseToken")

	if actual == nil || actual.Error() != expectedStr {
		t.Errorf("got %v\nwant %v", actual, expectedStr)
	}

	// TODO Normal scenario
}

func TestRefundSubscription(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	client, _ := New(jsonKey)
	expectedStr := "googleapi: Error 404: No application was found for the given package name., applicationNotFound"
	actual := client.RefundSubscription(ctx, "package", "productID", "purchaseToken")

	if actual == nil || actual.Error() != expectedStr {
		t.Errorf("got %v\nwant %v", actual, expectedStr)
	}

	// TODO Normal scenario
}

func TestRevokeSubscription(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	client, _ := New(jsonKey)
	expectedStr := "googleapi: Error 404: No application was found for the given package name., applicationNotFound"
	actual := client.RevokeSubscription(ctx, "package", "productID", "purchaseToken")

	if actual == nil || actual.Error() != expectedStr {
		t.Errorf("got %v\nwant %v", actual, expectedStr)
	}

	// TODO Normal scenario
}

func TestVerifySignature(t *testing.T) {
	t.Parallel()
	receipt := []byte(`{"orderId":"GPA.xxxx-xxxx-xxxx-xxxxx","packageName":"my.package","productId":"myproduct","purchaseTime":1437564796303,"purchaseState":0,"developerPayload":"user001","purchaseToken":"some-token"}`)

	type in struct {
		pubkey  string
		receipt []byte
		sig     string
	}

	tests := []struct {
		name  string
		in    in
		err   error
		valid bool
	}{
		{
			name: "public key is invalid base64 format",
			in: in{
				pubkey:  "dummy_public_key",
				receipt: receipt,
				sig:     "gj0N8LANKXOw4OhWkS1UZmDVUxM1UIP28F6bDzEp7BCqcVAe0DuDxmAY5wXdEgMRx/VM1Nl2crjogeV60OqCsbIaWqS/ZJwdP127aKR0jk8sbX36ssyYZ0DdZdBdCr1tBZ/eSW1GlGuD/CgVaxns0JaWecXakgoV7j+RF2AFbS4=",
			},
			err:   errors.New("failed to decode public key"),
			valid: false,
		},
		{
			name: "public key is not rsa public key",
			in: in{
				pubkey:  "JTbngOdvBE0rfdOs3GeuBnPB+YEP1w/peM4VJbnVz+hN9Td25vPjAznX9YKTGQN4iDohZ07wtl+zYygIcpSCc2ozNZUs9pV0s5itayQo22aT5myJrQmkp94ZSGI2npDP4+FE6ZiF+7khl3qoE0rVZq4G2mfk5LIIyTPTSA4UvyQ=",
				receipt: receipt,
				sig:     "gj0N8LANKXOw4OhWkS1UZmDVUxM1UIP28F6bDzEp7BCqcVAe0DuDxmAY5wXdEgMRx/VM1Nl2crjogeV60OqCsbIaWqS/ZJwdP127aKR0jk8sbX36ssyYZ0DdZdBdCr1tBZ/eSW1GlGuD/CgVaxns0JaWecXakgoV7j+RF2AFbS4=",
			},
			err:   errors.New("failed to parse public key"),
			valid: false,
		},
		{
			name: "signature is invalid base64 format",
			in: in{
				pubkey:  "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDGvModvVUrqJ9C5fy8J77ZQ7JDC6+tf5iK8C74/3mjmcvwo4nmprCgzR/BQIEuZWJi8KX+jiJUXKXF90JPsXHkKAPq6A1SCga7kWvs/M8srMpjNS9zJdwZF+eDOR0+lJEihO04zlpAV9ybPJ3Q621y1HUeVpwdxDNLQpJTuIflnwIDAQAB",
				receipt: receipt,
				sig:     "invalid_signature",
			},
			err:   errors.New("failed to decode signature"),
			valid: false,
		},
		{
			name: "signature is invalid",
			in: in{
				pubkey:  "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDGvModvVUrqJ9C5fy8J77ZQ7JDC6+tf5iK8C74/3mjmcvwo4nmprCgzR/BQIEuZWJi8KX+jiJUXKXF90JPsXHkKAPq6A1SCga7kWvs/M8srMpjNS9zJdwZF+eDOR0+lJEihO04zlpAV9ybPJ3Q621y1HUeVpwdxDNLQpJTuIflnwIDAQAB",
				receipt: receipt,
				sig:     "JTbngOdvBE0rfdOs3GeuBnPB+YEP1w/peM4VJbnVz+hN9Td25vPjAznX9YKTGQN4iDohZ07wtl+zYygIcpSCc2ozNZUs9pV0s5itayQo22aT5myJrQmkp94ZSGI2npDP4+FE6ZiF+7khl3qoE0rVZq4G2mfk5LIIyTPTSA4UvyQ=",
			},
			err:   nil,
			valid: false,
		},
		{
			name: "normal",
			in: in{
				pubkey:  "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDGvModvVUrqJ9C5fy8J77ZQ7JDC6+tf5iK8C74/3mjmcvwo4nmprCgzR/BQIEuZWJi8KX+jiJUXKXF90JPsXHkKAPq6A1SCga7kWvs/M8srMpjNS9zJdwZF+eDOR0+lJEihO04zlpAV9ybPJ3Q621y1HUeVpwdxDNLQpJTuIflnwIDAQAB",
				receipt: receipt,
				sig:     "gj0N8LANKXOw4OhWkS1UZmDVUxM1UIP28F6bDzEp7BCqcVAe0DuDxmAY5wXdEgMRx/VM1Nl2crjogeV60OqCsbIaWqS/ZJwdP127aKR0jk8sbX36ssyYZ0DdZdBdCr1tBZ/eSW1GlGuD/CgVaxns0JaWecXakgoV7j+RF2AFbS4=",
			},
			err:   nil,
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, err := VerifySignature(tt.in.pubkey, tt.in.receipt, tt.in.sig)

			if valid != tt.valid {
				t.Errorf("input: %v\nget: %t\nwant: %t\n", tt.in, valid, tt.valid)
			}

			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("input: %v\nget: %s\nwant: %s\n", tt.in, err, tt.err)
			}
		})
	}
}
