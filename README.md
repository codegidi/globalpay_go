# Globalpay_go

Global pay go  is a library for using the [Globalpay] API for go


### Installing
import "github.com/codegidi/Globalpay_go"

### Usage
*    The steps for carrying out a transaction is as follows:
*    1. Get an access token by calling the Client Authorisation method
*    2. Use the access_token to send initiate your transaction by calling the Transaction initiaion method
*    3. Redirect to GlobalPay transaction interface using the redirectUri retured in the Transaction initiation call
*    4. After transaction has been done, you will be redirected to the provided redirectUrl provided with the transactionReference as a querystring
*    5. Validate the result by using the Retrieve transaction call


#### Client Authentication
	client := Globalpay.NewClientAuth()

	clientRequest := &TransactionRegistrationRequest{
		GrantType:      "{string}",
		ClientId: 		"{string}",
		ClientSecret:   "{string}",
	}

	clientAuthenticationResponse, err := client.Transaction.AuthenticateClient(transactionRequest)


##### Transaction Initialization
    accessToken := {ACCESS TOKEN}
	client := Globalpay.NewClient(accessToken)

	product := map[string][]map[string]string{
            "Products":   {{"Name": "{string}", "UnitPrice" : "{string}", "Quantity" : "{string}"}, {"Name": "{string}", "UnitPrice" : "{string}", "Quantity" : "{string}"}},
    }

	transactionRequest := &TransactionRegistrationRequest{
		Name : 				"{string}",
		ReturnUrl : 		"{string}",
		CustomerIp : 		"{string}",
		MerchantReference : "{string}",
		MerchantId : 		"{string}",
		Description : 		"{string}",
		CurrencyCode : 		"{string eg NGN for naira}",
		TotalAmount : 		"{string}",
		PaymentMethod : 	"{string}",
		TransactionType : 	"{string}",
		ConnectionMode : 	"{string}",
		Customer : map[string]interface{}{"Email":"{string}","Mobile":"{string}","Firstname":"{string}","Lastname":"{string}"},
		Product : product
	}

	transactionInitiationResponse, err := client.Transaction.Initialize(transactionRequest)

##### Transaction Verification
   	accessToken := {ACCESS TOKEN}
	client := Globalpay.NewClient(accessToken)

	transactionToVerify := &RetrieveTransactionRequest{
		MerchantId : 				"{string}",
		MerchantReference : 		"{string}",
		TransactionReference : 		"{string}",
	}
	
	transactionVerificationResponse, err := client.Transaction.Verify(transactionToVerify)


