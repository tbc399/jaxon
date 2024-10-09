package auth

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/jmoiron/sqlx"
	"jaxon.app/jaxon/internal/auth/otp"
	"jaxon.app/jaxon/internal/auth/sessions"
	"jaxon.app/jaxon/internal/auth/templates"
	"jaxon.app/jaxon/internal/auth/users"
)

// Prefixed with "login/"
func Router() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("GET /login", getLoginPage)
	router.HandleFunc("POST /login", submitLogin)
	router.HandleFunc("GET /login/otp/{otp_id}", getOneTimePass)
	router.HandleFunc("GET /login/otp/{otp_id}/check", checkLoginStatus)
	router.HandleFunc("GET /login/magic/{magic_token}", getOtpValidationPage)
	router.HandleFunc("POST /login/magic/{magic_token}", submitOtpValidation)
	return router
}

func getLoginPage(response http.ResponseWriter, request *http.Request) {
	slog.Info("Retrieving login page")
	templates.Login().Render(request.Context(), response)
}

func submitLogin(response http.ResponseWriter, request *http.Request) {
	db := request.Context().Value("db").(*sqlx.DB)
	email := request.PostFormValue("email")
	slog.Info("Handling login form", "email", email)
	otpass, err := otp.Create(email, db)
	if err != nil {
		slog.Error("Failed to create a new one time password")
		// TODO: give something back to the user
		return
	}

	go otp.SendEmail(email, otpass, db)

	http.Redirect(response, request, fmt.Sprintf("/auth/login/otp/%s", otpass.Id), 303)
}

func getOneTimePass(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value("db").(*sqlx.DB)
	otp_id := r.PathValue("otp_id")
	otpass, err := otp.Fetch(otp_id, db)
	if err != nil {
		slog.Error("Failed to fetch the otpass", "otp_id", otp_id)
		// TODO: return something to the user
	}
	slog.Info("otpass found", "otp", otp_id)
	if otpass.IsExpired() {
		// TODO: give back something more specific to this case of the OTP expiring
		w.WriteHeader(http.StatusNotFound)
		templates.NotFound().Render(r.Context(), w)
	} else {
		templates.LoginPending(otpass).Render(r.Context(), w)
	}
}

// Meant to be polled by the front end in order to check for user completion
// of the login flow
func checkLoginStatus(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value("db").(*sqlx.DB)
	otp_id := r.PathValue("otp_id")
	otpass, err := otp.Fetch(otp_id, db)
	if err != nil {
		slog.Warn("Otp not found", "otp", otp_id)
		w.Header().Add("HX-Redirect", "/sorry")
		w.WriteHeader(http.StatusOK)
		return
		// TODO: return something to the user
	}

	if otpass.IsExpired() {
		// TODO: need a page to show that otp has expired
		slog.Info("Otp has expired", "otp", otpass.Id)
		w.Header().Add("HX-Redirect", "/sorry")
		w.WriteHeader(http.StatusOK)
		return
	}

	slog.Info("Checking for active session", "otp", otpass.Id, "email", otpass.Email)

	session, err := sessions.FetchByOtpId(otpass.Id, db)

	if err != nil {
		// user has not completed login flow
		slog.Info("Failed to get session", "error", err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	slog.Info("Found session", "session", session.Id)
	w.Header().Add("HX-Redirect", "/budgets")
	cookie := &http.Cookie{Name: "trbl_session", Value: session.Id, Path: "/"}
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
}

func getOtpValidationPage(w http.ResponseWriter, r *http.Request) {
	// TODO: I think this should return a page directing the user to go back to the original login page
	db := r.Context().Value("db").(*sqlx.DB)
	magicToken := r.PathValue("magic_token")
	otpass, err := otp.FetchByMagicToken(magicToken, db)
	if err != nil {
		slog.Error("Failed to fetch otp for magic token", "magic_token", magicToken)
		http.Redirect(w, r, "/sorry", http.StatusSeeOther)
		return
	}

	templates.LoginValidation(otpass).Render(r.Context(), w)
}

func submitOtpValidation(w http.ResponseWriter, r *http.Request) {
	// TODO: I think this should return a page directing the user to go back to the original login page
	// TODO: Need to validate otp expiry
	// TODO: Do we validate the device?
	// TODO: validate the # of attempts is under a threshold of 3, for example

	db := r.Context().Value("db").(*sqlx.DB)
	magicToken := r.PathValue("magic_token")
	otpass, err := otp.FetchByMagicToken(magicToken, db)
	if err != nil {
		slog.Info("No otp found for magic token", "magic_token", magicToken)
		http.Redirect(w, r, "/sorry", http.StatusSeeOther)
		return
	}

	if otpass.IsExpired() {
		slog.Info("Otp has expired", "otp", otpass.Id)
		// TODO: need a dedicated "has expired" page
		http.Redirect(w, r, "/sorry", http.StatusSeeOther)
		return
	}

	slog.Debug("Looking for existing user", "email", otpass.Email)
	user, err := users.FetchByEmail(otpass.Email, db)

	if err != nil {
		// TODO: return something to the user
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if user == nil {
		user := users.New(otpass.Email)
		user.Save(db)
		slog.Info("New user created", "email", otpass.Email)

		// TODO: This should probably happen in a background task unless asyncio.Queue handles it
		//go lucy.publish(UserCreated(user_id=user.id, db=db))

		// TODO: should this be an event handler to let the response come back timely?
		//go stripe.Customer.create_async(name="", email=user.email)
	}

	session := sessions.New(user.Id, otpass.Id)
	session.Save(db)

	templates.LoginSuccess().Render(r.Context(), w)
}
