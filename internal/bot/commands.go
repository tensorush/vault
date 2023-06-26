package bot

import tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var (
	startMessageEN = `Hi! 👋 I'm a password vault bot 🔐.

ℹ️ My commands:
/set service_name login password - saves your password for the specified service.
/get service_name - retrieves your password for the specified service.
/del service_names - deletes your password for the specified service.

Messages are deleted every %d seconds, so that no one can see what you've entered 🤫.`

	startMessagePT = `Olá! 👋 Eu sou um bot de cofre de senhas 🔐.

ℹ️ Meus comandos:
/set service_name login password - guarda a tua palavra-passe para o serviço especificado.
/get service_name - recupera a sua palavra-passe para o serviço especificado.
/del service_names - apaga a tua palavra-passe para o serviço especificado.

As mensagens são apagadas a cada %d segundos, para que ninguém possa ver o que introduziste 🤫.`
)

var allMessages = map[string]messages{
	start: {
		English:    startMessageEN,
		Portuguese: startMessagePT,
	},
	set: {
		English:    setMessageEN,
		Portuguese: setMessagePT,
	},
	setErr: {
		English:    setErrMessageEN,
		Portuguese: setErrMessagePT,
	},
	get: {
		English:    getMessageEN,
		Portuguese: getMessagePT,
	},
	getErr: {
		English:    getErrMessageEN,
		Portuguese: getErrMessagePT,
	},
	del: {
		English:    delMessageEN,
		Portuguese: delMessagePT,
	},
	delErr: {
		English:    delErrMessageEN,
		Portuguese: delErrMessagePT,
	},

	wrongInputErr: {
		English:    wrongInputErrEN,
		Portuguese: wrongInputErrPT,
	},
	serviceNotFoundErr: {
		English:    serviceNotFoundErrEN,
		Portuguese: serviceNotFoundErrPT,
	},
}

// Group of constants for bot messages.
const (
	setMessageEN    = "Saved ✅"
	setErrMessageEN = "Error during saving! ⛔️"
	setMessagePT    = "Salvo ✅"
	setErrMessagePT = "Erro ao guardar! ⛔️"

	delMessageEN    = "Deleted 🗑"
	delErrMessageEN = "Error during deletion! ⛔️"
	delMessagePT    = "Eliminado 🗑"
	delErrMessagePT = "Erro durante a eliminação! ⛔️"

	getMessageEN    = "🔐 %s\n👤 Login: %s\n🔑 Password: %s\n"
	getErrMessageEN = "Error during retrieval! ⚒"
	getMessagePT    = "🔐 %s\n👤 Login: %s\n🔑 Palavra-passe: %s\n"
	getErrMessagePT = "Erro durante a recuperação! ⚒"

	wrongInputErrEN = "Wrong input for command! ⛔️"
	wrongInputErrPT = "Entrada incorrecta para o comando! ⛔️"

	serviceNotFoundErrEN = "Service not found ❌"
	serviceNotFoundErrPT = "Serviço não encontrado ❌"
)

// Group of constants for handling messages from user.
const (
	start = "start"

	get    = "get"
	getErr = "getErr"

	set    = "set"
	setErr = "setErr"

	change     = "change"
	changeLang = "changeLang"

	del    = "del"
	delErr = "delErr"

	hide = "hide"

	wrongInputErr      = "Wrong input for command"
	serviceNotFoundErr = "Service not found"
)

const (
	hideKeyboard    = "hideKeyboard"
	setLangKeyboard = "setLangKeyboard"
	startKeyboard   = "startKeyboard"
)

// Map of  keyboard buttons.
var allKeyboards = map[string]keyboards{
	hideKeyboard: {
		English: tg.NewInlineKeyboardMarkup(
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("Hide \U0001FAE3", hide),
			),
		),
		Portuguese: tg.NewInlineKeyboardMarkup(
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("Ocultar mensagem \U0001FAE3", hide),
			),
		),
	},

	setLangKeyboard: {
		English: tg.NewInlineKeyboardMarkup(
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("English 🇬🇧", change+"::en"),
				tg.NewInlineKeyboardButtonData("Português 🇵🇹", change+"::pt"),
			),
		),
		Portuguese: tg.NewInlineKeyboardMarkup(
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("English 🇺🇸", change+"::en"),
				tg.NewInlineKeyboardButtonData("Português 🇵🇹", change+"::pt"),
			),
		),
	},

	startKeyboard: {
		English: tg.NewInlineKeyboardMarkup(
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("Change language 🌍", changeLang),
			),
		),
		Portuguese: tg.NewInlineKeyboardMarkup(
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("Alterar a língua 🌍", changeLang),
			),
		),
	},
}
