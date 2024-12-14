package html

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"yeetfile/backend/config"
	"yeetfile/backend/db"
	"yeetfile/backend/server/html/templates"
	"yeetfile/backend/server/session"
	"yeetfile/backend/server/upgrades"
	"yeetfile/backend/utils"
	"yeetfile/shared"
	"yeetfile/shared/constants"
	"yeetfile/shared/endpoints"
)

// FileVaultPageHandler returns the html template used for interacting with files
// (uploading, renaming, downloading, deleting) in the user's vault
func FileVaultPageHandler(w http.ResponseWriter, _ *http.Request, userID string) {
	userStorage, _, err := db.GetUserStorage(userID)
	if err != nil {
		handleError(w, "Error fetching vault", http.StatusInternalServerError)
		return
	}

	_ = templates.ServeTemplate(
		w,
		templates.VaultHTML,
		templates.VaultTemplate{
			Base: templates.BaseTemplate{
				LoggedIn: true,
				Title:    "File Vault",
				Page:     "vault",
				Javascript: []string{
					"file_vault.js",
					"render.js",
					"ponyfill.min.js",
					"pdf.min.mjs",
					"pdf.worker.min.mjs",
				},
				CSS:       []string{"vault.css"},
				Config:    config.HTMLConfig,
				Endpoints: endpoints.HTMLPageEndpoints,
			},
			StorageUsed:      userStorage.StorageUsed,
			StorageAvailable: userStorage.StorageAvailable,
			VaultName:        "File Vault",
			IsPasswordVault:  false,
		},
	)
}

// PassVaultPageHandler returns the html template used for interacting with
// stored passwords/logins
func PassVaultPageHandler(w http.ResponseWriter, _ *http.Request, userID string) {
	passCount, maxPassCount, err := db.GetUserPassCount(userID)
	if err != nil {
		log.Println(passCount, maxPassCount, err)
		handleError(w, "Error fetching pass vault", http.StatusInternalServerError)
		return
	}

	_ = templates.ServeTemplate(
		w,
		templates.VaultHTML,
		templates.VaultTemplate{
			Base: templates.BaseTemplate{
				LoggedIn: true,
				Title:    "Password Vault",
				Page:     "pass",
				Javascript: []string{
					"pass_vault.js",
					"render.js",
				},
				CSS:       []string{"vault.css"},
				Config:    config.HTMLConfig,
				Endpoints: endpoints.HTMLPageEndpoints,
			},
			StorageUsed:      passCount,
			StorageAvailable: maxPassCount,
			VaultName:        "Password Vault",
			IsPasswordVault:  true,
		},
	)
}

// SendPageHandler returns the html template used for sending files
func SendPageHandler(w http.ResponseWriter, req *http.Request) {
	_ = templates.ServeTemplate(
		w,
		templates.SendHTML,
		templates.LoginTemplate{
			Base: templates.BaseTemplate{
				LoggedIn: session.IsValidSession(w, req),
				Title:    "Send",
				Page:     "send",
				Javascript: []string{
					"jszip.min.js",
					"share.js",
				},
				CSS:       []string{"send.css"},
				Config:    config.HTMLConfig,
				Endpoints: endpoints.HTMLPageEndpoints,
			},
			Meter: 0,
		},
	)
}

// DownloadPageHandler returns the HTML page for downloading a file
func DownloadPageHandler(w http.ResponseWriter, req *http.Request) {
	_ = templates.ServeTemplate(
		w,
		templates.DownloadHTML,
		templates.Template{Base: templates.BaseTemplate{
			LoggedIn: session.IsValidSession(w, req),
			Title:    "Download",
			Javascript: []string{
				"ponyfill.min.js",
				"download.js",
			},
			CSS:       []string{"download.css"},
			Config:    config.HTMLConfig,
			Endpoints: endpoints.HTMLPageEndpoints,
		}},
	)
}

// SignupPageHandler returns the HTML page for signing up for an account
func SignupPageHandler(w http.ResponseWriter, _ *http.Request) {
	_ = templates.ServeTemplate(
		w,
		templates.SignupHTML,
		templates.SignupTemplate{
			Base: templates.BaseTemplate{
				LoggedIn:   false,
				Title:      "Create Account",
				Javascript: []string{"signup.js"},
				CSS:        []string{"auth.css"},
				Config:     config.HTMLConfig,
				Endpoints:  endpoints.HTMLPageEndpoints,
			},
			ServerPasswordRequired: config.YeetFileConfig.PasswordHash != nil,
			EmailConfigured:        config.YeetFileConfig.Email.Configured,
		},
	)
}

// LoginPageHandler returns the HTML page for logging in
func LoginPageHandler(w http.ResponseWriter, _ *http.Request) {
	_ = templates.ServeTemplate(
		w,
		templates.LoginHTML,
		templates.Template{
			Base: templates.BaseTemplate{
				LoggedIn:   false,
				Title:      "Log In",
				Javascript: []string{"login.js"},
				CSS:        []string{"auth.css"},
				Config:     config.HTMLConfig,
				Endpoints:  endpoints.HTMLPageEndpoints,
			},
		},
	)
}

// AccountPageHandler returns the HTML page for a user managing their account
func AccountPageHandler(w http.ResponseWriter, req *http.Request, userID string) {
	user, err := db.GetUserByID(userID)
	if err != nil || user.ID != userID {
		handleError(w, "Unable to fetch user info", http.StatusUnauthorized)
		return
	}

	successMsg, errorMsg := generateAccountMessages(req)
	isYearly := req.URL.Query().Has("yearly")
	hasHint := user.PasswordHint != nil && len(user.PasswordHint) > 0
	obscuredEmail, _ := utils.ObscureEmail(user.Email)
	isPrevUpgraded := user.UpgradeExp.Year() >= 2024

	var durationFilter constants.UpgradeDuration
	if isYearly {
		durationFilter = constants.DurationYear
	} else {
		durationFilter = constants.DurationMonth
	}

	products := upgrades.GetUpgrades(durationFilter)

	_ = templates.ServeTemplate(
		w,
		templates.AccountHTML,
		templates.AccountTemplate{
			Base: templates.BaseTemplate{
				LoggedIn:   true,
				Title:      "My Account",
				Page:       "account",
				Javascript: []string{"account.js"},
				CSS:        []string{"account.css"},
				Config:     config.HTMLConfig,
				Endpoints:  endpoints.HTMLPageEndpoints,
			},
			Email:            obscuredEmail,
			EmailConfigured:  config.YeetFileConfig.Email.Configured,
			PaymentID:        user.PaymentID,
			ExpString:        user.UpgradeExp.Format("2 Jan 2006"),
			IsActive:         time.Now().Before(user.UpgradeExp),
			IsPrevUpgraded:   isPrevUpgraded,
			SendAvailable:    shared.ReadableFileSize(user.SendAvailable),
			SendUsed:         shared.ReadableFileSize(user.SendUsed),
			StorageAvailable: shared.ReadableFileSize(user.StorageAvailable),
			StorageUsed:      shared.ReadableFileSize(user.StorageUsed),
			IsYearly:         isYearly,
			BillingEndpoints: endpoints.BillingPageEndpoints,
			HasPasswordHint:  hasHint,
			Has2FA:           user.Secret != nil && len(user.Secret) > 0,
			ErrorMessage:     errorMsg,
			SuccessMessage:   successMsg,
			Products:         products,
		},
	)
}

// VerifyPageHandler returns the HTML page for verifying the user's email
func VerifyPageHandler(w http.ResponseWriter, req *http.Request) {
	email := req.URL.Query().Get("email")
	code := req.URL.Query().Get("code")

	_ = templates.ServeTemplate(
		w,
		templates.VerificationHTML,
		templates.VerificationTemplate{
			Base: templates.BaseTemplate{
				LoggedIn:   session.IsValidSession(w, req),
				Title:      "Verify",
				Javascript: []string{"verify.js"},
				CSS:        nil,
				Config:     config.HTMLConfig,
				Endpoints:  endpoints.HTMLPageEndpoints,
			},
			Email: email,
			Code:  code,
		},
	)
}

// ChangeEmailPageHandler returns the HTML page for updating a user's email. This
// can be for users changing their email from one to another, or for an account
// ID-only user adding an email to their account
func ChangeEmailPageHandler(w http.ResponseWriter, req *http.Request, id string) {
	email, err := db.GetUserEmailByID(id)
	if err != nil {
		handleError(w, "Unable to fetch user", http.StatusInternalServerError)
		return
	} else if len(email) > 0 {
		pathSegments := strings.Split(req.URL.Path, "/")
		changeID := pathSegments[len(pathSegments)-1]
		valid := db.IsChangeIDValid(changeID, id)
		if !valid {
			handleError(w, "Invalid access", http.StatusUnauthorized)
			return
		}
	}

	_ = templates.ServeTemplate(
		w,
		templates.ChangeEmailHTML,
		templates.ChangeEmailTemplate{
			Base: templates.BaseTemplate{
				LoggedIn:   true,
				Title:      "Change Email",
				Javascript: []string{"change_email.js"},
				CSS:        []string{"change.css"},
				Config:     config.HTMLConfig,
				Endpoints:  endpoints.HTMLPageEndpoints,
			},
			CurrentEmail: email,
		},
	)
}

func ChangePasswordPageHandler(w http.ResponseWriter, _ *http.Request, _ string) {
	_ = templates.ServeTemplate(
		w,
		templates.ChangePasswordHTML,
		templates.Template{
			Base: templates.BaseTemplate{
				LoggedIn:   true,
				Title:      "Change Password",
				Javascript: []string{"change_password.js"},
				CSS:        []string{"change.css"},
				Config:     config.HTMLConfig,
				Endpoints:  endpoints.HTMLPageEndpoints,
			},
		},
	)
}

func ChangeHintPageHandler(w http.ResponseWriter, _ *http.Request, _ string) {
	_ = templates.ServeTemplate(
		w,
		templates.ChangeHintHTML,
		templates.Template{
			Base: templates.BaseTemplate{
				LoggedIn:   true,
				Title:      "Set Password Hint",
				Javascript: []string{"change_hint.js"},
				CSS:        []string{"change.css"},
				Config:     config.HTMLConfig,
				Endpoints:  endpoints.HTMLPageEndpoints,
			},
		},
	)
}

func TwoFactorPageHandler(w http.ResponseWriter, _ *http.Request, _ string) {
	_ = templates.ServeTemplate(
		w,
		templates.TwoFactorHTML,
		templates.Template{
			Base: templates.BaseTemplate{
				LoggedIn:   true,
				Title:      "Two-Factor Auth",
				Javascript: []string{"enable_2fa.js"},
				CSS:        nil,
				Config:     config.HTMLConfig,
				Endpoints:  endpoints.HTMLPageEndpoints,
			},
		},
	)
}

func ServerInfoPageHandler(w http.ResponseWriter, req *http.Request) {
	serverInfo := config.GetServerInfoStruct()
	hasRestrictions := serverInfo.PasswordRestricted || serverInfo.MaxUserCountSet
	storageStr := shared.ReadableFileSize(serverInfo.DefaultStorage)
	sendStr := fmt.Sprintf("%s / month", shared.ReadableFileSize(serverInfo.DefaultSend))
	if serverInfo.DefaultSend < 0 {
		sendStr = "Unlimited"
	}

	_ = templates.ServeTemplate(
		w,
		templates.ServerInfoHTML,
		templates.InfoTemplate{
			Base: templates.BaseTemplate{
				LoggedIn:   session.IsValidSession(w, req),
				Title:      "Server Info",
				Javascript: nil,
				CSS:        nil,
				Config:     config.HTMLConfig,
				Endpoints:  endpoints.HTMLPageEndpoints,
			},
			HasRestrictions:    hasRestrictions,
			PasswordRestricted: serverInfo.PasswordRestricted,
			MaxUserCountSet:    serverInfo.MaxUserCountSet,
			StorageBackend:     serverInfo.StorageBackend,
			EmailConfigured:    serverInfo.EmailConfigured,
			BillingEnabled:     serverInfo.BillingEnabled,
			StripeEnabled:      serverInfo.StripeEnabled,
			BTCPayEnabled:      serverInfo.BTCPayEnabled,
			DefaultStorage:     storageStr,
			DefaultSend:        sendStr,
			Products:           serverInfo.Upgrades,
		},
	)
}

// ForgotPageHandler returns the HTML page for resetting a user's password
func ForgotPageHandler(w http.ResponseWriter, req *http.Request) {
	email := req.URL.Query().Get("email")
	code := req.URL.Query().Get("code")
	_ = templates.ServeTemplate(
		w,
		templates.ForgotHTML,
		templates.ForgotPasswordTemplate{
			Base: templates.BaseTemplate{
				LoggedIn:   false,
				Title:      "Forgot Password",
				Javascript: []string{"forgot.js"},
				CSS:        []string{"auth.css"},
				Config:     config.HTMLConfig,
				Endpoints:  endpoints.HTMLPageEndpoints,
			},
			Email: email,
			Code:  code,
		},
	)
}

func CheckoutCompleteHandler(w http.ResponseWriter, req *http.Request) {
	from := req.URL.Query().Get("from")

	var note string
	title := "Checkout complete!"
	desc := "If you have an email address associated with your YeetFile " +
		"account, you should receive confirmation of your order shortly."

	if from == "btcpay" {
		note = "BTCPay orders can sometimes take 5 minutes or longer " +
			"to finalize. Your account will be updated once the " +
			"transaction has been validated. "
	}

	_ = templates.ServeTemplate(
		w,
		templates.CheckoutCompleteHTML,
		templates.CheckoutCompleteTemplate{
			Base: templates.BaseTemplate{
				LoggedIn:   true,
				Title:      "Checkout Complete",
				Javascript: nil,
				CSS:        nil,
				Config:     config.HTMLConfig,
				Endpoints:  endpoints.HTMLPageEndpoints,
			},
			Title:       title,
			Description: desc,
			Note:        note,
		},
	)
}

// generateAccountMessages takes a request and generates success and error messages from
// the data contained in the request.
func generateAccountMessages(req *http.Request) (string, string) {
	success := req.URL.Query().Get("success")
	fromBTC := req.URL.Query().Get("btcpay")

	successMsg := ""
	errorMsg := ""
	if len(success) > 0 && success == "1" {
		successMsg = "Successfully updated account! "

		if len(fromBTC) > 0 && fromBTC == "1" {
			successMsg += "BTCPay orders can take up to 5 minutes " +
				"to finalize. Your account will be updated once " +
				"your transaction has been validated. "
		}
	} else if len(success) > 0 && success == "0" {
		errorMsg = "Failed to update account!"
	}

	return successMsg, errorMsg
}

func handleError(w http.ResponseWriter, msg string, status int) {
	w.WriteHeader(status)
	_, _ = w.Write([]byte(msg))
}
