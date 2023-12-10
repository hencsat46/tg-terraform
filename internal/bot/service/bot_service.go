package service

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/VanLavr/tg-bot/internal/bot/handler/http"
	"github.com/VanLavr/tg-bot/internal/models"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

type service struct {
	form *models.AdvancedForm
}

func New() http.BotUsecase {
	return &service{form: &models.AdvancedForm{}}
}

func (s *service) HandleMessage(text string, chadID int64) (*models.ResponseBodyButtonsDTO, *models.ResponseBodyTextDTO, *models.ResponseBodyKeyboardDTO, error) {
	switch text {
	case "/provider":
		menue, err := s.HandleProvider(chadID)
		if err != nil {
			return nil, nil, nil, err
		}

		return menue, nil, nil, nil

	case "/form":
		return nil, nil, &models.ResponseBodyKeyboardDTO{
			ChatId: chadID,
			Text:   "Выберите действие",
			ReplyMarkup: models.ReplyKeyboardMarkup{
				Keyboard: [][]models.KeyboardButton{
					{
						{Text: "Провайдер"},
						{Text: "VPC"},
						{Text: "VPC Subnet"},
					},
				},
			},
		}, nil

	case "/start":
		return nil, &models.ResponseBodyTextDTO{
			ChatID: chadID,
			Text:   "добро пожаловать",
		}, nil, nil
	case "/print":
		s.printForm()
		return nil, nil, nil, nil
	case "/generate":
		return nil, s.generateTf(chadID), nil, nil
	default:
		panic("not implemented")
	}
}

func (s *service) HandleProvider(chatID int64) (*models.ResponseBodyButtonsDTO, error) {
	log.Println(chatID)
	buttons := models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButtons{{
			{Text: "Провайдер", CallbackData: "provider"},
			{Text: "Регион", CallbackData: "region"},
		}},
	}

	return &models.ResponseBodyButtonsDTO{
		ChatID:  chatID,
		Text:    "Выберите параметры провайдера",
		Buttons: buttons,
	}, nil
}

func (s *service) HandleVPC(chatID int64) (*models.ResponseBodyButtonsDTO, error) {
	log.Println(chatID)
	buttons := models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButtons{{
			{Text: "Имя", CallbackData: "name"},
			{Text: "CIDR", CallbackData: "cidr"},
			{Text: "Имя ресурса", CallbackData: "resource"},
		}},
	}

	return &models.ResponseBodyButtonsDTO{
		ChatID:  chatID,
		Text:    "Выберите параметры VPC",
		Buttons: buttons,
	}, nil
}

func (s *service) HandleVPCSubnet(chatId int64) (*models.ResponseBodyButtonsDTO, error) {
	buttons := models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButtons{{
			{Text: "Имя в облаке", CallbackData: "name"},
			{Text: "Имя в конфиге", CallbackData: "resource"},
			{Text: "CIDR", CallbackData: "cidr"},
			{Text: "Getaway ip", CallbackData: "gatewayIp"},
		}},
	}

	return &models.ResponseBodyButtonsDTO{
		ChatID:  chatId,
		Text:    "Выберите параметры подсети VPC",
		Buttons: buttons,
	}, nil
}

func (s *service) AddProvider(text string) {
	s.form.Type = text
	log.Println("жопа")
}

func (s *service) AddProviderRegion(text string) {
	s.form.Provider.Region = text
	log.Println(s.form)
}

func (s *service) AddVPCName(text string) {
	s.form.Vpc.Name = text
}

func (s *service) AddCIDR(text string) {
	s.form.Vpc.Cidr = text
}

func (s *service) AddResource(text string) {
	s.form.Vpc.ResourceName = text
}

func (s *service) printForm() {
	log.Println(s.form)
}

func (s *service) AddSubnetName(text string) {
	s.form.VpcSubnet.Name = text
}

func (s *service) AddSubnetResource(text string) {
	s.form.VpcSubnet.Resource = text
}

func (s *service) AddSubnetCidr(text string) {
	s.form.VpcSubnet.Cidr = text
}

func (s *service) AddSubnetGateway(text string) {
	s.form.VpcSubnet.GatewayIp = text
}

func (s *service) generateTf(chatId int64) *models.ResponseBodyTextDTO {
	tfCode := ""
	if s.form.Type == "Advanced" {
		tfCode += "terraform {\n\trequired_providers {\n\t\tsbercloud = {\n\t\t\tsourse = \"tf.repo.sbc.space/sbercloud-terraform/sbercloud\"\n\t\t}\n\t}\n}\n\n"

		tfCode += fmt.Sprintf("provider \"sbercloud\" {\n\taccess_key=\"your-access-key\"\n\tsecret_key=\"your-secret-key\"\n\tregion = \"%s\"\n}\n\n", s.form.Provider.Region)

		tfCode += fmt.Sprintf("resource \"sbercloud_vpc\" \"%s\"{\n\tname = \"%s\"\n\tcidr = \"%s\"\n}\n\n", s.form.Vpc.ResourceName, s.form.Vpc.Name, s.form.Vpc.Cidr)

		tfCode += fmt.Sprintf("resource \"sbercloud_vpc_subnet\" \"%s\" {\n\tname = \"%s\"\n\tcidr = \"%s\"\n\tgateway_ip = \"%s\"\n\tvpc_id = sbercloud_vpc.%s.id\n}\n\n", s.form.VpcSubnet.Resource, s.form.VpcSubnet.Name, s.form.VpcSubnet.Cidr, s.form.VpcSubnet.GatewayIp, s.form.Vpc.ResourceName)

	}

	tfCode += ""

	// file, err := os.OpenFile("../../../files/main.tf", os.O_CREATE, 0222)
	// if err != nil {
	// 	log.Println(err)
	// }

	// file.Write([]byte(tfCode))
	// file.Close()

	response := &models.ResponseBodyTextDTO{ChatID: chatId, Text: tfCode, Entity: []models.MessageEntity{{Type: "code", Offset: 0, Length: len(tfCode), Language: "Terraform"}}}

	return response

}

// 192.12.12.12/12
func (s *service) ValidateIP(ip string) bool {
	numbers := strings.Split(ip, ".")
	lastNumber := strings.Split(numbers[3], "/")
	for i := 0; i < 3; i++ {
		if _, err := strconv.Atoi(numbers[i]); err != nil {
			return false
		}
	}

	if _, err := strconv.Atoi(lastNumber[0]); err != nil {
		return false
	}

	if _, err := strconv.Atoi(lastNumber[1]); err != nil {
		return false
	}

	return true
}
