package config

import (
	"log"
	"os"
	"yeetfile/web/utils"
)

const version = "0.0.1-beta"

// =============================================================================
// General configuration
// =============================================================================

var callbackDomain = os.Getenv("YEETFILE_CALLBACK_DOMAIN")
var defaultUserStorage = utils.GetEnvVarInt("YEETFILE_DEFAULT_USER_STORAGE", 0)
var defaultUserSend = utils.GetEnvVarInt("YEETFILE_DEFAULT_USER_SEND", 0)
var maxNumUsers = utils.GetEnvVarInt("YEETFILE_MAX_NUM_USERS", -1)

// =============================================================================
// Email configuration (used in account verification and billing reminders)
// =============================================================================

type EmailConfig struct {
	Configured bool
	Address    string
	Host       string
	Port       string
	Password   string
}

var email = EmailConfig{
	Configured: false,
	Address:    os.Getenv("YEETFILE_EMAIL_ADDR"),
	Host:       os.Getenv("YEETFILE_EMAIL_HOST"),
	Port:       os.Getenv("YEETFILE_EMAIL_PORT"),
	Password:   os.Getenv("YEETFILE_EMAIL_PASSWORD"),
}

// =============================================================================
// Billing configuration (Stripe)
// =============================================================================

type StripeBillingConfig struct {
	Configured    bool
	Key           string
	WebhookSecret string
	PortalLink    string

	SubNoviceMonthly     string
	SubNoviceMonthlyLink string
	SubNoviceYearly      string
	SubNoviceYearlyLink  string

	SubRegularMonthly     string
	SubRegularMonthlyLink string
	SubRegularYearly      string
	SubRegularYearlyLink  string

	SubAdvancedMonthly     string
	SubAdvancedMonthlyLink string
	SubAdvancedYearly      string
	SubAdvancedYearlyLink  string

	//Add50GBSend     string
	//Add50GBSendLink string
	//
	//Add100GBSend     string
	//Add100GBSendLink string
	//
	//Add250GBSend     string
	//Add250GBSendLink string
}

var stripeBilling = StripeBillingConfig{
	Key:           os.Getenv("YEETFILE_STRIPE_KEY"),
	WebhookSecret: os.Getenv("YEETFILE_STRIPE_WEBHOOK_SECRET"),
	PortalLink:    os.Getenv("YEETFILE_STRIPE_PORTAL_LINK"),

	SubNoviceMonthly:     os.Getenv("YEETFILE_STRIPE_SUB_NOVICE_MONTHLY"),
	SubNoviceMonthlyLink: os.Getenv("YEETFILE_STRIPE_SUB_NOVICE_MONTHLY_LINK"),
	SubNoviceYearly:      os.Getenv("YEETFILE_STRIPE_SUB_NOVICE_YEARLY"),
	SubNoviceYearlyLink:  os.Getenv("YEETFILE_STRIPE_SUB_NOVICE_YEARLY_LINK"),

	SubRegularMonthly:     os.Getenv("YEETFILE_STRIPE_SUB_REGULAR_MONTHLY"),
	SubRegularMonthlyLink: os.Getenv("YEETFILE_STRIPE_SUB_REGULAR_MONTHLY_LINK"),
	SubRegularYearly:      os.Getenv("YEETFILE_STRIPE_SUB_REGULAR_YEARLY"),
	SubRegularYearlyLink:  os.Getenv("YEETFILE_STRIPE_SUB_REGULAR_YEARLY_LINK"),

	SubAdvancedMonthly:     os.Getenv("YEETFILE_STRIPE_SUB_ADVANCED_MONTHLY"),
	SubAdvancedMonthlyLink: os.Getenv("YEETFILE_STRIPE_SUB_ADVANCED_MONTHLY_LINK"),
	SubAdvancedYearly:      os.Getenv("YEETFILE_STRIPE_SUB_ADVANCED_YEARLY"),
	SubAdvancedYearlyLink:  os.Getenv("YEETFILE_STRIPE_SUB_ADVANCED_YEARLY_LINK"),

	//Add50GBSend:     os.Getenv("YEETFILE_STRIPE_ADD_50GB_SEND"),
	//Add50GBSendLink: os.Getenv("YEETFILE_STRIPE_ADD_50GB_SEND_LINK"),
	//
	//Add100GBSend:     os.Getenv("YEETFILE_STRIPE_ADD_100GB_SEND"),
	//Add100GBSendLink: os.Getenv("YEETFILE_STRIPE_ADD_100GB_SEND_LINK"),
	//
	//Add250GBSend:     os.Getenv("YEETFILE_STRIPE_ADD_250GB_SEND"),
	//Add250GBSendLink: os.Getenv("YEETFILE_STRIPE_ADD_250GB_SEND_LINK"),
}

// =============================================================================
// Billing configuration (BTCPay)
// =============================================================================

type BTCPayBillingConfig struct {
	Configured    bool
	APIKey        string
	WebhookSecret string
	StoreID       string
	ServerURL     string

	SubNoviceMonthlyLink   string
	SubNoviceYearlyLink    string
	SubRegularMonthlyLink  string
	SubRegularYearlyLink   string
	SubAdvancedMonthlyLink string
	SubAdvancedYearlyLink  string
}

var btcPayBilling = BTCPayBillingConfig{
	APIKey:        os.Getenv("YEETFILE_BTCPAY_API_KEY"),
	WebhookSecret: os.Getenv("YEETFILE_BTCPAY_WEBHOOK_SECRET"),
	StoreID:       os.Getenv("YEETFILE_BTCPAY_STORE_ID"),
	ServerURL:     os.Getenv("YEETFILE_BTCPAY_SERVER_URL"),

	SubNoviceMonthlyLink:   os.Getenv("YEETFILE_BTCPAY_SUB_NOVICE_MONTHLY_LINK"),
	SubNoviceYearlyLink:    os.Getenv("YEETFILE_BTCPAY_SUB_NOVICE_YEARLY_LINK"),
	SubRegularMonthlyLink:  os.Getenv("YEETFILE_BTCPAY_SUB_REGULAR_MONTHLY_LINK"),
	SubRegularYearlyLink:   os.Getenv("YEETFILE_BTCPAY_SUB_REGULAR_YEARLY_LINK"),
	SubAdvancedMonthlyLink: os.Getenv("YEETFILE_BTCPAY_SUB_ADVANCED_MONTHLY_LINK"),
	SubAdvancedYearlyLink:  os.Getenv("YEETFILE_BTCPAY_SUB_ADVANCED_YEARLY_LINK"),
}

// =============================================================================
// Full server config
// =============================================================================

type ServerConfig struct {
	CallbackDomain     string
	DefaultUserStorage int
	DefaultUserSend    int
	MaxUserCount       int
	CurrentUserCount   int
	Email              EmailConfig
	StripeBilling      StripeBillingConfig
	BTCPayBilling      BTCPayBillingConfig
	BillingEnabled     bool
	Version            string
}

var YeetFileConfig ServerConfig

func init() {
	email.Configured = !utils.IsStructMissingAnyField(email)
	stripeBilling.Configured = email.Configured &&
		!utils.IsStructMissingAnyField(stripeBilling)
	btcPayBilling.Configured = email.Configured &&
		!utils.IsStructMissingAnyField(btcPayBilling)

	YeetFileConfig = ServerConfig{
		CallbackDomain:     callbackDomain,
		DefaultUserStorage: defaultUserStorage,
		DefaultUserSend:    defaultUserSend,
		MaxUserCount:       maxNumUsers,
		Email:              email,
		StripeBilling:      stripeBilling,
		BTCPayBilling:      btcPayBilling,
		BillingEnabled:     stripeBilling.Configured || btcPayBilling.Configured,
		Version:            version,
	}

	log.Printf("Configuration:\n"+
		"  Email:            %v\n"+
		"  Billing (Stripe): %v\n"+
		"  Billing (BTCPay): %v\n",
		email.Configured,
		stripeBilling.Configured,
		btcPayBilling.Configured)
}
