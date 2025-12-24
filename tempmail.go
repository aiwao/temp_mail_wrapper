package tempmail

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
)

const baseURL = "https://web2.temp-mail.org"
const newMailURL = baseURL + "/mailbox"
const messagesURL = baseURL + "/messages"

//messagesURL + {messageID} + /attachment + {attachmentID} = attachmentBodyURL

const UA = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36"

type Account struct {
    Token   string `json:"token"`
    Address string `json:"mailbox"`
}

type MessagePreview struct {
    Id               string `json:"_id"`
    ReceivedAt       int    `json:"receivedAt"`
    From             string `json:"from"`
    Subject          string `json:"subject"`
    BodyPreview      string `json:"bodyPreview"`
    AttachmentsCount int    `json:"attachmentsCount"`
}

type MessagePreviews struct {
    Address  string           `json:"mailbox"`
    Previews []MessagePreview `json:"messages"`
}

type AttachmentInfo struct {
    Filename string `json:"filename"`
    Id       int    `json:"_id"`
    Size     int    `json:"size"`
    Mimetype string `json:"mimetype"`
    Cid      string `json:"cid"`
}

type Message struct {
    Id               string           `json:"_id"`
    ReceivedAt       int              `json:"receivedAt"`
    User             string           `json:"user"`
    Mailbox          string           `json:"mailbox"`
    From             string           `json:"from"`
    Subject          string           `json:"subject"`
    BodyPreview      string           `json:"bodyPreview"`
    BodyHtml         string           `json:"bodyHtml"`
    AttachmentsCount int              `json:"attachmentsCount"`
    Attachments      []AttachmentInfo `json:"attachments"`
    CreatedAt        time.Time        `json:"createdAt"`
}

func NewAccount(client *http.Client) (*Account, error) {
    c := client
    if client == nil {
        c = http.DefaultClient
    }
    req, err := http.NewRequest("POST", newMailURL, nil)
    if err != nil {
        return nil, err
    }
    req.Header.Set("User-Agent", UA)
    res, err := c.Do(req)
    if err != nil {
        return nil, err
    }
    defer res.Body.Close()
    b, err := io.ReadAll(res.Body)
    if err != nil {
        return nil, err
    }
    fmt.Println(string(b))
    var account Account
    err = json.Unmarshal(b, &account)
    if err != nil {
        return nil, err
    }
    return &account, nil
}

func (a *Account) MessagePreviews(client *http.Client) (MessagePreviews, error) {
    c := client
    if client == nil {
        c = http.DefaultClient
    }
    req, err := http.NewRequest("GET", messagesURL, nil)
    if err != nil {
        return MessagePreviews{}, err
    }
    req.Header.Set("User-Agent", UA)
    req.Header.Set("Authorization", "Bearer "+a.Token)
    res, err := c.Do(req)
    if err != nil {
        return MessagePreviews{}, err
    }
    defer res.Body.Close()
    b, err := io.ReadAll(res.Body)
    if err != nil {
        return MessagePreviews{}, err
    }
    var messagePreviews MessagePreviews
    err = json.Unmarshal(b, &messagePreviews)
    if err != nil {
        return MessagePreviews{}, err
    }
    return messagePreviews, nil
}

func (a *Account) Message(client *http.Client, messageID string) (*Message, error) {
    c := client
    if client == nil {
        c = http.DefaultClient
    }
    req, err := http.NewRequest("GET", messagesURL+"/"+messageID, nil)
    if err != nil {
        return nil, err
    }
    req.Header.Set("User-Agent", UA)
    req.Header.Set("Authorization", "Bearer "+a.Token)
    res, err := c.Do(req)
    if err != nil {
        return nil, err
    }
    defer res.Body.Close()
    b, err := io.ReadAll(res.Body)
    if err != nil {
        return nil, err
    }
    var message *Message
    err = json.Unmarshal(b, message)
    if err != nil {
        return nil, err
    }
    return message, nil
}

func (a *Account) Attachment(client *http.Client, messageID string, attachmentID int) ([]byte, error) {
    c := client
    if client == nil {
        c = http.DefaultClient
    }
    req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/%d", messagesURL, messageID, attachmentID), nil)
    if err != nil {
        return nil, err
    }
    req.Header.Set("User-Agent", UA)
    req.Header.Set("Authorization", "Bearer "+a.Token)
    res, err := c.Do(req)
    if err != nil {
        return nil, err
    }
    defer res.Body.Close()
    b, err := io.ReadAll(res.Body)
    if err != nil {
        return nil, err
    }
    return b, nil
}
