package handlers

import (
	echo "github.com/labstack/echo/v4"
)

func ErrorHandler(err error, c echo.Context) {

	he := err.(*echo.HTTPError)
	c.Logger().Error(err)
	c.JSON(he.Code, he)
	//echo.HTTPError{he.Code, he.Message, nil})

}
