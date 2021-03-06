package session

import (
	"net/http"

	"github.com/gorilla/sessions"
)

var (
	cookieStore   *sessions.CookieStore
	cookieOptions sessions.Options
	cookieName    string
)

type Init struct {
	AuthenticationKey string
	EncryptionKey     string
	CookieName        string
	CookieOptions     sessions.Options
}

func Initialise(init Init) {

	cookieStore = sessions.NewCookieStore(
		[]byte(init.AuthenticationKey),
		[]byte(init.EncryptionKey),
	)

	cookieName = init.CookieName
	cookieOptions = init.CookieOptions
}

func getSession(r *http.Request) (*sessions.Session, error) {

	session, err := cookieStore.Get(r, cookieName)
	if err == nil {
		session.Options = &cookieOptions
	}
	return session, err
}

func Get(r *http.Request, key string) (value string, err error) {

	session, err := getSession(r)
	if err != nil {
		return "", err
	}

	if session.Values[key] == nil {
		session.Values[key] = ""
	}

	return session.Values[key].(string), nil
}

func GetAll(r *http.Request) (ret map[string]string, err error) {

	ret = map[string]string{}

	session, err := getSession(r)
	if err != nil {
		return ret, err
	}

	for k, v := range session.Values {
		ret[k.(string)] = v.(string)
	}

	return ret, err
}

func Set(r *http.Request, name string, value string) (err error) {

	session, err := getSession(r)
	if err != nil {
		return err
	}

	session.Values[name] = value

	return nil
}

func SetMany(r *http.Request, values map[string]string) (err error) {

	session, err := getSession(r)
	if err != nil {
		return err
	}

	for k, v := range values {
		session.Values[k] = v
	}

	return nil
}

func Delete(r *http.Request, key string) (err error) {

	session, err := getSession(r)
	if err != nil {
		return err
	}

	delete(session.Values, key)

	return nil
}

func DeleteMany(r *http.Request, keys []string) (err error) {

	session, err := getSession(r)
	if err != nil {
		return err
	}

	for _, v := range keys {
		delete(session.Values, v)
	}

	return nil
}

func DeleteAll(r *http.Request) (err error) {

	session, err := getSession(r)
	if err != nil {
		return err
	}

	session.Values = make(map[interface{}]interface{})

	return nil
}

type FlashGroup string

func GetFlashes(r *http.Request, group FlashGroup) (flashes []string, err error) {

	session, err := getSession(r)
	if err != nil {
		return nil, err
	}

	interfaces := session.Flashes(string(group))

	for _, v := range interfaces {
		flashes = append(flashes, v.(string))
	}

	return flashes, err
}

func SetFlash(r *http.Request, group FlashGroup, flash string) (err error) {

	session, err := getSession(r)
	if err != nil {
		return err
	}

	session.AddFlash(flash, string(group))

	return nil
}

func Save(w http.ResponseWriter, r *http.Request) (err error) {

	session, err := getSession(r)
	if err != nil {
		return err
	}

	return session.Save(r, w)
}
