package telegram

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/Nikita-Kolbin/dictionary/internal/app/model"
)

const (
	tgBotHost = "api.telegram.org"

	sendMessageMethod = "sendMessage"
	editMessageMethod = "editMessageText"
	getUpdatesMethod  = "getUpdates"

	parseMode      = "MarkdownV2"
	disablePreview = `{"is_disabled": true}`
)

type TGClient struct {
	host       string
	basePath   string
	httpClient http.Client
}

func New(token string) *TGClient {
	return &TGClient{
		host:       tgBotHost,
		basePath:   newBasePath(token),
		httpClient: http.Client{},
	}
}

func (c *TGClient) Updates(offset, limit int) ([]*model.Update, error) {
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest(getUpdatesMethod, q)
	if err != nil {
		return nil, fmt.Errorf("can't get updates: %w", err)
	}

	var resp model.UpdatesResponse

	if json.Unmarshal(data, &resp) != nil {
		return nil, fmt.Errorf("can't get updates: %w", err)
	}

	return resp.Result, nil
}

func (c *TGClient) Send(chatID int, msg string, withFormat bool) (*model.Response, error) {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", msg)

	if withFormat {
		q.Add("parse_mode", parseMode)
		q.Add("link_preview_options", disablePreview)
	}

	byteResp, err := c.doRequest(sendMessageMethod, q)
	if err != nil {
		return nil, fmt.Errorf("can't send message: %w", err)
	}

	resp := &model.Response{}
	if err := json.Unmarshal(byteResp, resp); err != nil {
		return nil, fmt.Errorf("can't parse response: %w", err)
	}

	return resp, nil
}

func (c *TGClient) Edit(msg string, chatID, msgID int, withFormat bool, key *model.InlineKeyboardMarkup) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("message_id", strconv.Itoa(msgID))
	q.Add("text", msg)

	if withFormat {
		q.Add("parse_mode", parseMode)
		q.Add("link_preview_options", disablePreview)
	}

	if key != nil {
		strKey, err := json.Marshal(key)
		if err != nil {
			return fmt.Errorf("can't send message: %w", err)
		}
		q.Add("reply_markup", string(strKey))
	} else {
		q.Add("reply_markup", "")
	}

	if _, err := c.doRequest(editMessageMethod, q); err != nil {
		return fmt.Errorf("can't edit message: %w", err)
	}

	return nil
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *TGClient) doRequest(method string, query url.Values) ([]byte, error) {
	u := url.URL{
		Scheme:   "https",
		Host:     c.host,
		Path:     path.Join(c.basePath, method),
		RawQuery: query.Encode(),
	}

	resp, err := c.httpClient.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("can't do request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("can't do request: %w", err)
	}

	strBody := string(body)
	_ = strBody

	return body, nil
}
