package http

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/VanLavr/tg-bot/internal/models"
	errir "github.com/VanLavr/tg-bot/pkg/error"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

type handler struct {
	bot              models.Bot
	u                BotUsecase
	formStatus       bool
	providerMessage  string
	vpcMessage       string
	vpcSubnetMessage string
}

type BotUsecase interface {
	HandleMessage(string, int64) (*models.ResponseBodyButtonsDTO, *models.ResponseBodyTextDTO, *models.ResponseBodyKeyboardDTO, error)
	HandleProvider(int64) (*models.ResponseBodyButtonsDTO, error)
	AddProvider(string)
	AddProviderRegion(string)
	AddVPCName(string)
	AddCIDR(string)
	AddResource(string)
	AddSubnetName(string)
	AddSubnetResource(string)
	AddSubnetCidr(string)
	AddSubnetGateway(string)
	HandleVPC(int64) (*models.ResponseBodyButtonsDTO, error)
	ValidateIP(string) bool
	HandleVPCSubnet(int64) (*models.ResponseBodyButtonsDTO, error)
}

func New(u BotUsecase, bot models.Bot) models.BotHTTPDelivery {
	return &handler{
		u:   u,
		bot: bot,
	}
}

func (h *handler) Register(e *echo.Echo) {
	e.POST("/", h.HandleWebHook)
}

func (h *handler) HandleWebHook(c echo.Context) error {
	var body models.RequestBodyDTO

	if err := c.Bind(&body); err != nil {
		log.Println("cannot parse message")
	}

	if body.MessageInfo.ChatInfo.ChatId != 0 {
		h.bot.ChatID = body.MessageInfo.ChatInfo.ChatId
	}
	message := body.MessageInfo.Text

	if len(message) != 0 {
		if message[0] == '/' {
			menueMessage, textMessage, keyboardMessage, err := h.u.HandleMessage(message, h.bot.ChatID)
			if err != nil {
				log.Println("error:", err)
				return nil

			} else if menueMessage != nil {

				resp := h.sendInlineKeyboardMessage(menueMessage)
				log.Println("HERE")
				h.decodeAndPrintResponse(resp)

			} else if keyboardMessage != nil {
				h.formStatus = true
				resp := h.sendKeyboardMessage(keyboardMessage)
				log.Println("Keyboard message has been sent")
				h.decodeAndPrintResponse(resp)
			} else {
				resp := h.sendMessage(textMessage)
				log.Println("here2")
				h.decodeAndPrintResponse(resp)

			}
		} else if h.formStatus {
			if h.providerMessage != "" {
				h.u.AddProviderRegion(message)
				h.sendProviderForm(body.MessageInfo.ChatInfo.ChatId)
				h.providerMessage = ""
			} else if h.vpcMessage != "" {
				switch h.vpcMessage {
				case "name":
					h.u.AddVPCName(message)
					h.sendVPCForm(body.MessageInfo.ChatInfo.ChatId)
				case "cidr":
					h.u.AddCIDR(message)
					h.sendVPCForm(body.MessageInfo.ChatInfo.ChatId)
				case "resource":
					h.u.AddResource(message)
					h.sendVPCForm(body.MessageInfo.ChatInfo.ChatId)
				}
			} else if h.vpcSubnetMessage != "" {
				switch h.vpcSubnetMessage {
				case "name":
					h.u.AddSubnetName(message)
					h.sendSubnetForm(body.MessageInfo.ChatInfo.ChatId)
				case "resource":
					h.u.AddSubnetResource(message)
					h.sendSubnetForm(body.MessageInfo.ChatInfo.ChatId)
				case "cidr":
					h.u.AddSubnetCidr(message)
					h.sendSubnetForm(body.MessageInfo.ChatInfo.ChatId)
				case "gatewayIp":
					h.u.AddSubnetGateway(message)
					h.sendSubnetForm(body.MessageInfo.ChatInfo.ChatId)
				}
			} else {
				switch message {
				case "Провайдер":
					menu, _ := h.u.HandleProvider(body.MessageInfo.ChatInfo.ChatId)
					h.sendInlineKeyboardMessage(menu)
				case "VPC":
					menu, _ := h.u.HandleVPC(body.MessageInfo.ChatInfo.ChatId)
					h.sendInlineKeyboardMessage(menu)
				case "VPC Subnet":
					menu, _ := h.u.HandleVPCSubnet(body.MessageInfo.ChatInfo.ChatId)
					h.sendInlineKeyboardMessage(menu)
				}

			}
		} else {
			resp := h.sendMessage(&models.ResponseBodyTextDTO{Text: "Нужна команда", ChatID: h.bot.ChatID})
			log.Println("here3")
			h.decodeAndPrintResponse(resp)
		}
	} else if len(body.Callback.Data) != 0 {
		log.Println("This is callback")
		res := h.sendCallBackQuery(&models.AnswerCallback{Id: body.Callback.ID, Text: ""})
		h.decodeAndPrintResponse(res)

		log.Println("Callback:", body)

		switch body.Callback.MessageInfo.Text {
		case "Выберите параметры провайдера":
			switch body.Callback.Data {
			case "provider":
				log.Println("get provider")
				res := h.sendInlineKeyboardMessage(&models.ResponseBodyButtonsDTO{
					ChatID: body.Callback.MessageInfo.ChatInfo.ChatId,
					Text:   "Выберите провайдера",
					Buttons: models.InlineKeyboardMarkup{
						InlineKeyboard: [][]models.InlineKeyboardButtons{
							{{
								Text: "Advanced", CallbackData: "advanced",
							}},
						},
					},
				})
				h.decodeAndPrintResponse(res)
			case "region":
				h.providerMessage = "region"
				h.sendMessage(&models.ResponseBodyTextDTO{ChatID: body.Callback.MessageInfo.ChatInfo.ChatId, Text: "Введите регион"})
			}
		case "Выберите параметры VPC":
			switch body.Callback.Data {
			case "name":
				h.vpcMessage = "name"
				h.sendMessage(&models.ResponseBodyTextDTO{ChatID: body.Callback.MessageInfo.ChatInfo.ChatId, Text: "Введите название VPC"})
			case "cidr":
				h.vpcMessage = "cidr"
				h.sendMessage(&models.ResponseBodyTextDTO{ChatID: body.Callback.MessageInfo.ChatInfo.ChatId, Text: "Введите CIDR"})
			case "resource":
				h.vpcMessage = "resource"
				h.sendMessage(&models.ResponseBodyTextDTO{ChatID: body.Callback.MessageInfo.ChatInfo.ChatId, Text: "Введите имя ресурса"})
			}
		case "Выберите провайдера":
			switch body.Callback.Data {
			case "advanced":
				h.u.AddProvider("Advanced")
				menu, _ := h.u.HandleProvider(body.Callback.MessageInfo.ChatInfo.ChatId)
				res := h.sendInlineKeyboardMessage(menu)
				h.decodeAndPrintResponse(res)

			}
		case "Выберите параметры подсети VPC":
			log.Println("мы в подсети")
			switch body.Callback.Data {
			case "name":
				h.vpcSubnetMessage = "name"
				h.sendMessage(&models.ResponseBodyTextDTO{ChatID: body.Callback.MessageInfo.ChatInfo.ChatId, Text: "Введите имя подсети VPC, которое будет в облаке"})
			case "resource":
				h.vpcSubnetMessage = "resource"
				h.sendMessage(&models.ResponseBodyTextDTO{ChatID: body.Callback.MessageInfo.ChatInfo.ChatId, Text: "Введите имя подсети VPC, которое будет в конфиге"})
			case "cidr":
				h.vpcSubnetMessage = "cidr"
				h.sendMessage(&models.ResponseBodyTextDTO{ChatID: body.Callback.MessageInfo.ChatInfo.ChatId, Text: "Введите адрес сети"})
			case "gatewayIp":
				h.vpcSubnetMessage = "gatewayIp"
				h.sendMessage(&models.ResponseBodyTextDTO{ChatID: body.Callback.MessageInfo.ChatInfo.ChatId, Text: "Введите адрес gateway'я"})
			}
		}

		return c.JSON(http.StatusOK, nil)
	}
	return nil
}

func (h *handler) sendMessage(message *models.ResponseBodyTextDTO) *http.Response {
	encodedBody, err := json.Marshal(message)
	if err != nil {
		log.Println("Marshalling error:", err)
		return nil
	}

	response, err := http.Post(h.bot.TGURL+"/sendMessage", "application/json", bytes.NewBuffer(encodedBody))
	if err != nil {
		log.Println(err)
		return nil
	}

	return response
}

func (h *handler) sendKeyboardMessage(keyboardMarkup *models.ResponseBodyKeyboardDTO) *http.Response {
	encodedBody, err := json.Marshal(keyboardMarkup)
	if err != nil {
		log.Println("Marshalling error:", err)
		return nil
	}

	response, err := http.Post(h.bot.TGURL+"/sendMessage", "application/json", bytes.NewBuffer(encodedBody))
	if err != nil {
		log.Println(err)
		return nil
	}

	return response
}

func (h *handler) sendInlineKeyboardMessage(menue *models.ResponseBodyButtonsDTO) *http.Response {
	encodedBody, err := json.Marshal(menue)
	if err != nil {
		log.Println("Marshalling error:", err)
		return nil
	}

	response, err := http.Post(h.bot.TGURL+"/sendMessage", "application/json", bytes.NewBuffer(encodedBody))
	if err != nil {
		log.Println(err)
		return nil
	}

	return response
}

func (h *handler) sendCallBackQuery(message *models.AnswerCallback) *http.Response {
	encodedBody, err := json.Marshal(message)
	if err != nil {
		log.Println("Marshaling error:", err)
	}

	response, err := http.Post(h.bot.TGURL+"/answerCallbackQuery", "application/json", bytes.NewBuffer(encodedBody))
	if err != nil {
		log.Println(err)
		return nil
	}

	return response
}

func (h *handler) GetMe(ctx context.Context) (*http.Response, error) {
	response, err := http.Get(h.bot.TGURL + "/getMe")
	if err != nil {
		return nil, err
	}

	if response.StatusCode < 200 || response.StatusCode > 299 {
		return response, errir.ErrBadResponse
	}

	return response, nil
}

func (h *handler) decodeAndPrintResponse(response *http.Response) {
	defer response.Body.Close()

	temp := make(map[string]interface{})
	if err := json.NewDecoder(response.Body).Decode(&temp); err != nil {
		log.Println("decode error: ", err)
		return
	}

	log.Println(temp)
}

func (h *handler) GetChat(ctx context.Context) (*http.Response, error) {
	response, err := http.Get(h.bot.TGURL + "/getChat")
	if err != nil {
		return nil, err
	}

	if response.StatusCode < 200 || response.StatusCode > 299 {
		return response, errir.ErrBadResponse
	}

	return response, nil
}

func (h *handler) sendVPCForm(chatId int64) *http.Response {
	buttons := models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButtons{{
			{Text: "Имя", CallbackData: "name"},
			{Text: "CIDR", CallbackData: "cidr"},
			{Text: "Имя ресурса", CallbackData: "resource"},
		}},
	}

	res := h.sendInlineKeyboardMessage(&models.ResponseBodyButtonsDTO{
		ChatID:  chatId,
		Text:    "Выберите параметры VPC",
		Buttons: buttons,
	})
	return res
}

func (h *handler) sendProviderForm(chatId int64) *http.Response {
	buttons := models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButtons{{
			{Text: "Провайдер", CallbackData: "provider"},
			{Text: "Регион", CallbackData: "region"},
		}},
	}

	res := h.sendInlineKeyboardMessage(&models.ResponseBodyButtonsDTO{
		ChatID:  chatId,
		Text:    "Выберите параметры провайдера",
		Buttons: buttons,
	})
	return res
}

func (h *handler) sendSubnetForm(chatId int64) *http.Response {
	buttons := models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButtons{{
			{Text: "Имя в облаке", CallbackData: "name"},
			{Text: "Имя в конфиге", CallbackData: "resource"},
			{Text: "CIDR", CallbackData: "cidr"},
			{Text: "Getaway ip", CallbackData: "gatewayIp"},
		}},
	}

	res := h.sendInlineKeyboardMessage(&models.ResponseBodyButtonsDTO{
		ChatID:  chatId,
		Text:    "Выберите параметры подсети VPC",
		Buttons: buttons,
	})
	return res
}
