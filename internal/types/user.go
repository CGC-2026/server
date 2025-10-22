package types

type ClerkUserData struct {
	ID             string `json:"id"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	EmailAddresses []struct {
		EmailAddress string `json:"email_address"`
		ID           string `json:"id"`
	} `json:"email_addresses"`
	ImageURL string `json:"image_url"`
}
