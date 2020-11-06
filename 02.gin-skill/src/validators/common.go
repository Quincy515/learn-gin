package validators

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"log"
)

var myvalid *validator.Validate

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		myvalid = v
	} else {
		log.Fatal("error validator")
	}
}
