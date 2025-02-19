package webdriver

import (
	"fmt"
	"os"

	"github.com/linweiyuan/go-chatgpt-api/api"
	"github.com/linweiyuan/go-chatgpt-api/util/logger"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

var WebDriver selenium.WebDriver

//goland:noinspection GoUnhandledErrorResult,SpellCheckingInspection
func init() {
	chatgptProxyServer := os.Getenv("CHATGPT_PROXY_SERVER")
	if chatgptProxyServer == "" {
		return
	}

	chromeArgs := []string{
		"--no-sandbox",
		"--disable-gpu",
		"--disable-dev-shm-usage",
		"--disable-blink-features=AutomationControlled",
		"--headless=new",
	}

	networkProxyServer := os.Getenv("NETWORK_PROXY_SERVER")
	if networkProxyServer != "" {
		chromeArgs = append(chromeArgs, "--proxy-server="+networkProxyServer)
	}

	WebDriver, _ = selenium.NewRemote(selenium.Capabilities{
		"chromeOptions": chrome.Capabilities{
			Args:            chromeArgs,
			ExcludeSwitches: []string{"enable-automation"},
		},
	}, chatgptProxyServer)

	if WebDriver == nil {
		logger.Error("Please make sure chatgpt proxy service is running")
		os.Exit(1)
		return
	}

	WebDriver.Get(api.ChatGPTUrl)

	if isReady(WebDriver) {
		logger.Info(api.ChatGPTWelcomeText)
		openNewTabAndChangeBackToOldTab()
	} else {
		if !isAccessDenied(WebDriver) {
			if HandleCaptcha(WebDriver) {
				logger.Info(api.ChatGPTWelcomeText)
				openNewTabAndChangeBackToOldTab()
			}
		}
	}
}

//goland:noinspection GoUnhandledErrorResult
func openNewTabAndChangeBackToOldTab() {
	WebDriver.ExecuteScript(fmt.Sprintf("open('%s');", api.ChatGPTUrl), nil)
	handles, _ := WebDriver.WindowHandles()
	WebDriver.SwitchWindow(handles[0])

	// to save conversations, (k,v): {"request message id": "response message data"}
	WebDriver.ExecuteScript("window.conversationMap = new Map();", nil)
}
