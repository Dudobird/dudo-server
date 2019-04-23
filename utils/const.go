package utils

//ContextToken for context key type
type ContextToken string

// TokenContextKey a simple wrapper for ContextToken
const TokenContextKey = ContextToken("MyAppToken")

// AdminContextKey save the status of user role
const AdminContextKey = ContextToken("IsAdmin")
