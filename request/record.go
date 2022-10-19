package request

// Type for get contacts request
type RequestGetContacts struct {
	ContactNumber string `json:"contactNumber"`
	ContactName   string `json:"contactName"`
	ClientId      string `json:"clientId"`
}

// Type for upload contacts request
type RequestUploadContacts struct {
	Url      string `json:"url"`
	ClientId string `json:"clientId"`
}
