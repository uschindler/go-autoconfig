package handlers

import (
	"encoding/xml"
	"github.com/labstack/echo"
	"go-autoconfig/config"
	"io/ioutil"
	"net/http"
  "strings"
  "regexp"
)

type Handler struct {
	Config *config.Config
}

type server struct {
	Host     string
	Port     int
	STARTTLS bool
}

func (h *Handler) GetDomain(ctx echo.Context) string {
  domain := h.Config.Domain
  if domain == "" {
    host := strings.Split(ctx.Request().Host, ":")[0]
    matched, _ := regexp.MatchString(`^[\w.-]+$`, host)
    if matched {
      return strings.TrimPrefix(host, "autoconfig.")
    } else {
      return "invalid"
    }
  }
  return domain;
}

func (h *Handler) Outlook(ctx echo.Context) error {
	var req struct {
		XMLName xml.Name `xml:"Autodiscover"`
		Text    string   `xml:",chardata"`
		Xmlns   string   `xml:"xmlns,attr"`
		Request struct {
			Text                     string `xml:",chardata"`
			EMailAddress             string `xml:"EMailAddress"`
			AcceptableResponseSchema string `xml:"AcceptableResponseSchema"`
		} `xml:"Request"`
	}

	b, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := xml.Unmarshal(b, &req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	data := struct {
		Schema string
		Email  string
		Domain string
		IMAP   *server
		SMTP   *server
	}{
		Schema: req.Request.AcceptableResponseSchema,
		Email:  req.Request.EMailAddress,
		Domain: h.GetDomain(ctx),
		IMAP: &server{
			Host:     h.Config.IMAP.Host,
			Port:     h.Config.IMAP.Port,
			STARTTLS: h.Config.IMAP.STARTTLS,
		},
		SMTP: &server{
			Host:     h.Config.SMTP.Host,
			Port:     h.Config.SMTP.Port,
			STARTTLS: h.Config.SMTP.STARTTLS,
		},
	}

	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationXMLCharsetUTF8)
	return ctx.Render(http.StatusOK, "outlook", data)
}

func (h *Handler) Thunderbird(ctx echo.Context) error {
	data := struct {
		Domain string
		IMAP   *server
		SMTP   *server
	}{
		Domain: h.GetDomain(ctx),
		IMAP: &server{
			Host:     h.Config.IMAP.Host,
			Port:     h.Config.IMAP.Port,
			STARTTLS: h.Config.IMAP.STARTTLS,
		},
		SMTP: &server{
			Host:     h.Config.SMTP.Host,
			Port:     h.Config.SMTP.Port,
			STARTTLS: h.Config.SMTP.STARTTLS,
		},
	}

	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationXMLCharsetUTF8)
	return ctx.Render(http.StatusOK, "thunderbird", data)
}

func (h *Handler) AppleMail(ctx echo.Context) error {
	var req struct {
		Email string `query:"email"`
	}
	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	data := struct {
		Email  string
		Domain string
		IMAP   *server
		SMTP   *server
	}{
		Email:  req.Email,
		Domain: h.GetDomain(ctx),
		IMAP: &server{
			Host:     h.Config.IMAP.Host,
			Port:     h.Config.IMAP.Port,
			STARTTLS: h.Config.IMAP.STARTTLS,
		},
		SMTP: &server{
			Host:     h.Config.SMTP.Host,
			Port:     h.Config.SMTP.Port,
			STARTTLS: h.Config.SMTP.STARTTLS,
		},
	}

	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationXMLCharsetUTF8)
	return ctx.Render(http.StatusOK, "applemail", data)
}
