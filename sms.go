package gatewayapi

const (
	DefaultSmsClass smsClass = "default"
	PremiumSmsClass smsClass = "premium"
	// The secret class can be used to blur the message content you send, used for very sensitive data.
	// It is priced as premium and uses the same routes, which ensures end to end encryption of your
	// messages. Access to the secret class will be very strictly controlled.
	SecretSmsClass smsClass = "secret"

	NormalPriority     priority = "NORMAL"
	BulkPriority       priority = "BULK"
	UrgentPriority     priority = "URGENT"
	VeryUrgentPriority priority = "VERY_URGENT" // requires premium sms class

	Utf8Encoding encoding = "UTF8"
	Ucs2Encoding encoding = "UCS2"

	DisplayDestinationAddress destinationAddress = "DISPLAY"
	MobileDestinationAddress  destinationAddress = "MOBILE"
	SimCardDestinationAddress destinationAddress = "SIMCARD"
	ExtUnitDestinationAddress destinationAddress = "EXTUNIT"
)

type smsClass string
type priority string
type encoding string
type destinationAddress string

// A SMS to be sent out
// see https://gatewayapi.com/docs/rest.html#advanced-usage
type SMS struct {
	Message        string             `json:"message"`                   // The content of the SMS, always specified in UTF-8 encoding, which we will transcode depending on the “encoding” field. The default is the usual GSM 03.38 encoding. Required unless payload is specified.
	Recipients     []Recipient        `json:"recipients"`                // Array of recipients, described below. The number of recipients in a single message is limited to 10.000. required
	Class          smsClass           `json:"class,omitempty"`           // Default ‘standard’. The message class, ‘standard’, ‘premium’ or ‘secret’ to use for this request. If specified it must be the same for all messages in the request. The secret class can be used to blur the message content you send, used for very sensitive data. It is priced as premium and uses the same routes, which ensures end to end encryption of your messages. Access to the secret class will be very strictly controlled.
	Sender         string             `json:"sender,omitempty"`          // Up to 11 alphanumeric characters, or 15 digits, that will be shown as the sender of the SMS. See SMS Sender (https://gatewayapi.com/docs/appendix.html#smssender)
	Sendtime       uint               `json:"sendtime,omitempty"`        // Unix timestamp (seconds since epoch) to schedule message sending at certain time.
	Tags           []string           `json:"tags,omitempty"`            // A list of string tags, which will be replaced with the tag values for each recipient.
	UserRef        string             `json:"userref,omitempty"`         // A transparent string reference, you may set to keep track of the message in your own systems. Returned to you when you receive a Delivery Status Notification (https://gatewayapi.com/docs/rest.html#delivery-status-notification).
	ValidityPeriod uint               `json:"validity_period,omitempty"` // Specified in seconds. If message is not delivered within this timespan, it will expire and you will get a notification. The minimum value is 60. Every value under 60 will be set to 60.
	Encoding       encoding           `json:"encoding,omitempy"`         // Encoding to use when sending the message. Defaults to ‘UTF8’, which means we will use GSM 03.38. Use UCS2 to send a unicode message.
	DestAddr       destinationAddress `json:"destaddr,omitempy"`         // One of ‘DISPLAY’, ‘MOBILE’, ‘SIMCARD’, ‘EXTUNIT’. Use display to do “flash sms”, a message displayed on screen immediately but not saved in the normal message inbox on the mobile device.
	Payload        string             `json:"payload,omitempy"`          //  If you are sending a binary SMS, ie. a SMS you have encoded yourself or with speciel content for feature phones (non-smartphones). You may specify a payload, encoded as Base64. If specified, message must not be set and tags are unavailable.
	Udh            string             `json:"udh,omitempty"`             // UDH to enable additional functionality for binary SMS, encoded as Base64.
	CallbackUrl    string             `json:"callbackUrl,omitempty"`     // If specified send status notifications to this URL, else use the default webhook.
	Label          string             `json:"label,omitempty"`           // A label added to each sent message, can be used to uniquely identify a customer or company that you sent the message on behalf of, to help with invoicing your customers. If specied it must be the same for all messages in the request.
	MaxParts       uint8              `json:"max_parts,omitempty"`       // A number between 1 and 255 used to limit the number of smses a single message will send. Can be used if you send smses from systems that generates messages that you can’t control, this way you can ensure that you don’t send very long smses. You will not be charged for more than the amount specified here. Can’t be used with Tags or BINARY smses.
	ExtraDetails   string             `json:"extra_details,omitempty"`   // To get more details about the number of parts sent to each recipient set this to ‘recipients_usage’. See example response below.
}

type Recipient struct {
	Msisdn    string   `json:"msisdn"`              // MSISDN aka the full mobile phone number of the recipient. Duplicates are not allowed in the same message. required
	TagValues []string `json:"tagvalues,omitempty"` //  A list of string values corresponding to the tags in message. The order and amount of tag values must exactly match the tags.
}
