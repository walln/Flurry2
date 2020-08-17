package global

import (
	firebase "firebase.google.com/go"
)

var globalFirebaseApp *firebase.App
var signingMethod string

func GetFirebaseApp() *firebase.App {
	return globalFirebaseApp
}

func SetFirebaseApp(f *firebase.App) {
	globalFirebaseApp = f
}

func GetSigningMethod() string {
	return signingMethod
}

func SetSigningMethod(s string) {
	signingMethod = s
}
