package tempmail

import "testing"

func Test(t *testing.T) {
    account, err := NewAccount(nil)
    if err != nil {
        t.Fatal(err)
    }
    previews, err := account.MessagePreviews(nil)
    if err != nil {
        t.Fatal(err)
    }
    t.Log(len(previews.Previews))
    for _, preview := range previews.Previews {
        t.Log(preview.Subject)
        message, err := account.Message(nil, preview.Id)
        if err != nil {
            t.Log(err)
            continue
        }
        t.Log(message.BodyHtml)
        if message.AttachmentsCount > 0 {
            for _, attachment := range message.Attachments {
                b, err := account.Attachment(nil, message.Id, attachment.Id)
                if err != nil {
                    t.Log(err)
                    continue
                }
                t.Log(string(b))
            }
        }
    }
}
