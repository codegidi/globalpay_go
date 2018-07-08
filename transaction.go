package globalpay

import (
	"fmt"
)


type TransactionService service

type ClientAuthenticationRequest struct {
	GrantType string `json:"grant_type,omitempty"`
	ClientId string `json:"client_id,omitempty"`
	ClientSecret string `json:"client_secret,omitempty"`
}

type ClientAuthenticationResponse struct {
	AccessToken string `json:"access_token,omitempty"`
	ExpiresIn string `json:"expires_in,omitempty"`
	TokenType string `json:"token_type,omitempty"`
}

type RetrieveTransactionRequest struct {
	MerchantId string `json:"merchantid,omitempty"`
	MerchantReference string `json:"merchantreference,omitempty"`
	TransactionReference string `json:"transactionreference,omitempty"`
}

type RetrieveTransactionResponse struct {
	MerchantId string `json:"merchantid,omitempty"`
	MerchantReference string `json:"merchantreference,omitempty"`
	TransactionReference string `json:"transactionreference,omitempty"`
}

type TransactionRegistrationRequest struct {
	Name string `json:"name,omitempty"`
	ReturnUrl string `json:"returnurl,omitempty"`
	CustomerIp string `json:"customerip,omitempty"`
	MerchantReference string `json:"merchantreference,omitempty"`
	MerchantId string `json:"merchantid,omitempty"`
	Description string `json:"description,omitempty"`
	CurrencyCode string `json:"currencycode,omitempty"`
	TotalAmount string `json:"totalamount,omitempty"`
	PaymentMethod string `json:"paymentmethod,omitempty"`
	TransactionType string `json:"transactionType,omitempty"`
	ConnectionMode string `json:"connectionmode,omitempty"`
	Customer Customer `json:"customer,omitempty"`
	Product []Product `json:"product,omitempty"`
}

type TransactionRegistrationResponse struct {
	Status Status `json:"status,omitempty"`
	RedirectUri string `json:"redirectUri,omitempty"`
	TransactionReference string `json:"transactionReference,omitempty"`
}

type Product struct {
	Name string `json:"name,omitempty"`
	UnitPrice string `json:"unitprice,omitempty"`
	Quantity string `json:"quantity,omitempty"`
}

type Customer struct {
	Email string `json:"email,omitempty"`
	Mobile string `json:"mobile,omitempty"`
	FirstName string `json:"firstname,omitempty"`
	LastName string `json:"lastname,omitempty"`
}

type Status struct {
	Status string `json:"Status,omitempty"`
	StatusCode string `json:"statusCode,omitempty"`
	ErrorDescription ErrorDescription `json:"errorDescription,omitempty"`
}

type ErrorDescription struct {
	Status         bool   `json:"iserror,omitempty"`
	StatusCode     string `json:"errorcode,omitempty"`
	ErrorParameter string `json:"errorparameter,omitempty"`
	Reason         string `json:"reason,omitempty"`
}

// Initialize initiates a transaction process
func (s *TransactionService) Initialize(txn *TransactionRegistrationRequest) (*TransactionRegistrationResponse, error) {
	u := fmt.Sprintf("/api/v3/Payment/SetRequest")
	resp := &TransactionRegistrationResponse{}
	err := s.client.Call("POST", u, txn, &resp)
	return resp, err
}

// Verify checks that transaction with the given reference exists
func (s *TransactionService) Verify(txn *RetrieveTransactionRequest) (*RetrieveTransactionResponse, error) {
	u := fmt.Sprintf("/api/v3/Payment/Retrieve")
	resp := &RetrieveTransactionResponse{}
	err := s.client.Call("POST", u, txn, &resp)
	return resp, err
}

// Authorization
func (s *TransactionService) AuthenticateClient(req *ClientAuthenticationRequest) (*ClientAuthenticationResponse, error) {
	txn := &ClientAuthenticationResponse{}
	err := s.client.CallFormPost("POST", "", req, txn)
	return txn, err
}


