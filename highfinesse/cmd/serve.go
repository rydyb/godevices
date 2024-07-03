package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rydyb/godevices/highfinesse"
)

type Quantity struct {
	Value float64
	Unit  string
}

type Serve struct {
	Address string `default:":8080" help:"Address to listen on."`
}

func (cmd *Serve) Run() error {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, highfinesse.WavelengthMeterVersionInfo())
	})

	e.GET("/pressure", func(c echo.Context) error {
		p, err := highfinesse.Pressure()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed measuring pressure: %s", err))
		}
		return c.JSON(http.StatusOK, Quantity{Value: p, Unit: "mbar"})
	})
	e.GET("/temperature", func(c echo.Context) error {
		T, err := highfinesse.Temperature()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed measuring temperature: %s", err))
		}
		return c.JSON(http.StatusOK, Quantity{Value: T, Unit: "°C"})
	})

	e.GET("/frequency/:channel", func(c echo.Context) error {
		channel, err := strconv.Atoi(c.Param("channel"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("failed parsing channel %d: %s", channel, err))
		}
		f, err := highfinesse.Frequency(uint32(channel))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed measuring frequency of channel %d: %s", channel, err))
		}
		return c.JSON(http.StatusOK, Quantity{Value: f, Unit: "THz"})
	})
	e.GET("/wavelength/:channel", func(c echo.Context) error {
		channel, err := strconv.Atoi(c.Param("channel"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("failed parsing channel %d: %s", channel, err))
		}
		f, err := highfinesse.Wavelength(uint32(channel))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed measuring wavelength of channel %d: %s", channel, err))
		}
		return c.JSON(http.StatusOK, Quantity{Value: f, Unit: "THz"})
	})

	return e.Start(cmd.Address)
}
