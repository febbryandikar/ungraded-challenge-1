package main

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
)

type User struct {
	Name string `required:"true" minLen:"2" maxLen:"100"`
	Age int `required:"true" min:"18" max:"65"`
	Email string `required:"true" email:"true" minLen:"2" maxLen:"100"`
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9. _%+-]+@[a-zA-Z0-9. -]+\.[a-zA-Z]{2,}$`)

func ValidateStructByTag(user interface{}) error {
	v := reflect.ValueOf(user)

	if v.Kind()!= reflect.Struct {
		return fmt.Errorf("%T is not a struct", user)
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		// validate required tag
        tagRequired := v.Type().Field(i).Tag.Get("required")

        if tagRequired == "true" {
            if field.Kind() == reflect.String {
                if len(field.String()) == 0 {
                    return fmt.Errorf("%s is required", v.Type().Field(i).Name)
                }
            } else if field.Kind() == reflect.Int {
                if field.Int() == 0 {
                    return fmt.Errorf("%s is required", v.Type().Field(i).Name)
                }
            }
        }

		// validate minLen tag
        if field.Kind() == reflect.String {
			tag := v.Type().Field(i).Tag.Get("minLen")
			tagMinLen, err := strconv.Atoi(tag)
			if err != nil {
				return err
			}
			if len(field.String()) < tagMinLen {
				return fmt.Errorf("%s is too short", v.Type().Field(i).Name)
			}
		}

		// validate min tag
        if field.Kind() == reflect.Int {
			tag := v.Type().Field(i).Tag.Get("min")
			tagMin, err := strconv.Atoi(tag)
			if err != nil {
				return err
			}
			if int(field.Int()) < tagMin {
				return fmt.Errorf("%s is too young", v.FieldByName("Name"))
			}
		}

		// validate maxLen tag
        if field.Kind() == reflect.String {
			tag := v.Type().Field(i).Tag.Get("maxLen")
			tagMaxLen, err := strconv.Atoi(tag)
			if err != nil {
				return err
			}
			if len(field.String()) > tagMaxLen {
				return fmt.Errorf("%s is too long", v.Type().Field(i).Name)
			}
		}

		// validate max tag
        if field.Kind() == reflect.Int {
			tag := v.Type().Field(i).Tag.Get("max")
			tagMax, err := strconv.Atoi(tag)
			if err != nil {
				return err
			}
			if int(field.Int()) > tagMax {
				return fmt.Errorf("%s is too old", v.FieldByName("Name"))
			}
		}

		// validate email tag
		tagEmail := v.Type().Field(i).Tag.Get("email")

		if tagEmail == "true" {
            if !emailRegex.MatchString(field.String()) {
                return fmt.Errorf("%s is not a valid email", v.Field(i))
            }
        }
	}


	return nil
}

func main() {
	user := User{
		Name: "Tony Stark",
        Age: 53,
        Email: "tonystark@avengers.com",
	}

	if err := ValidateStructByTag(user); err != nil {
        fmt.Println(err)
    } else {
        fmt.Println("Valid")
    }
}