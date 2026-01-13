package response

// var errorCodes = map[int]string{
// 	//Application Error Codes
// 	SH1-500: "INTERNAL_SERVER_ERROR", d
// 	SH1-400: "BAD_REQUEST",
// 	SH1-401: "UNAUTHORIZED",
// 	SH1-404: "NOT_FOUND",
// 	//Groq Error Codes
// 	SH1-429: "QUOTA_EXCEEDED",
// 	SH1-422: "UNPROCESSABLE_ENTITY",
// 	SH1-503: "SERVICE_UNAVAILABLE",
// }


//Errors From Upsteam Groq API

//Too many requests from Groq: 429 Quota exceeded
//Any other 4XX code from Groq: 422 Upstream_Input_Error/Unprocessable entity
//Any 5XX code from Groq: 503 Upstream_Output_Error/Service Unavailable
