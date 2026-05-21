package email

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"api/models"
)

// GraphMailer delivers messages through the Microsoft Graph API using an
// app-only (client-credentials) OAuth2 token, calling POST /users/{from}/sendMail.
// It targets Microsoft 365 / Office 365 mailboxes, which no longer accept basic
// SMTP auth. The Azure app must hold the Mail.Send APPLICATION permission with
// admin consent (ideally scoped to the sending mailbox via an Application Access
// Policy). The caller is responsible for supplying the decrypted client secret.
type GraphMailer struct {
	settings     models.EmailSettings
	tenantID     string
	clientID     string
	clientSecret string
}

// NewGraphMailer constructs a Graph mailer. tenantID is always the per-license
// value; clientID/clientSecret may be per-license or fall back to global config
// (resolved by the caller).
func NewGraphMailer(settings models.EmailSettings, tenantID, clientID, clientSecret string) *GraphMailer {
	return &GraphMailer{
		settings:     settings,
		tenantID:     tenantID,
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}

const (
	graphLoginBase = "https://login.microsoftonline.com"
	graphAPIBase   = "https://graph.microsoft.com"
	graphScope     = "https://graph.microsoft.com/.default"
	// maxGraphInlineAttachmentBytes caps the total attachment payload sent
	// inline via fileAttachment. Graph rejects a single request over ~4 MB;
	// larger payloads require an upload session (not implemented here).
	maxGraphInlineAttachmentBytes = 3 * 1024 * 1024
)

var graphHTTPClient = &http.Client{Timeout: 60 * time.Second}

// graphCachedToken holds an access token and the time it should be refreshed.
type graphCachedToken struct {
	token  string
	expiry time.Time
}

var (
	graphTokenMu    sync.Mutex
	graphTokenCache = map[string]graphCachedToken{} // key: tenant|clientID
)

// Send delivers msg via Graph sendMail. A non-202 response (or token failure)
// is returned as an error; the outbox worker decides retry vs terminal.
func (m *GraphMailer) Send(ctx context.Context, msg Message) error {
	from := msg.FromAddress
	if from == "" {
		from = m.settings.FromAddress
	}
	if from == "" {
		return fmt.Errorf("graph: no from address")
	}
	if len(msg.To)+len(msg.Cc)+len(msg.Bcc) == 0 {
		return fmt.Errorf("graph: no recipients")
	}

	token, err := m.token(ctx)
	if err != nil {
		return err
	}

	var attachments []graphAttachment
	var attachmentBytes int
	for _, a := range msg.Attachments {
		attachmentBytes += len(a.Content)
		attachments = append(attachments, graphAttachment{
			ODataType:    "#microsoft.graph.fileAttachment",
			Name:         a.Filename,
			ContentType:  a.ContentType,
			ContentBytes: base64.StdEncoding.EncodeToString(a.Content),
		})
	}
	if attachmentBytes > maxGraphInlineAttachmentBytes {
		return fmt.Errorf("graph: attachments total %d bytes exceed inline limit of %d (upload sessions not supported)", attachmentBytes, maxGraphInlineAttachmentBytes)
	}

	message := graphMessage{
		Subject:       msg.Subject,
		Body:          graphBody{ContentType: "HTML", Content: msg.Body},
		From:          &graphRecipient{EmailAddress: graphEmailAddress{Address: from, Name: msg.FromName}},
		ToRecipients:  graphRecipients(msg.To),
		CcRecipients:  graphRecipients(msg.Cc),
		BccRecipients: graphRecipients(msg.Bcc),
		Attachments:   attachments,
	}
	if msg.ReplyTo != "" {
		message.ReplyTo = []graphRecipient{{EmailAddress: graphEmailAddress{Address: msg.ReplyTo}}}
	}

	payload, err := json.Marshal(graphSendMailRequest{Message: message, SaveToSentItems: true})
	if err != nil {
		return fmt.Errorf("graph: marshal sendMail: %w", err)
	}

	endpoint := fmt.Sprintf("%s/v1.0/users/%s/sendMail", graphAPIBase, url.PathEscape(from))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("graph: build sendMail request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := graphHTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("graph: sendMail request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("graph: sendMail returned %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	return nil
}

// token returns a cached access token for (tenant, clientID), fetching a fresh
// one via the client-credentials flow when none is cached or it is near expiry.
func (m *GraphMailer) token(ctx context.Context) (string, error) {
	if m.tenantID == "" {
		return "", fmt.Errorf("graph: tenant id is required")
	}
	if m.clientID == "" || m.clientSecret == "" {
		return "", fmt.Errorf("graph: client id and secret are required")
	}
	key := m.tenantID + "|" + m.clientID

	graphTokenMu.Lock()
	if t, ok := graphTokenCache[key]; ok && time.Now().Before(t.expiry) {
		graphTokenMu.Unlock()
		return t.token, nil
	}
	graphTokenMu.Unlock()

	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("client_id", m.clientID)
	form.Set("client_secret", m.clientSecret)
	form.Set("scope", graphScope)

	tokenURL := fmt.Sprintf("%s/%s/oauth2/v2.0/token", graphLoginBase, url.PathEscape(m.tenantID))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("graph: build token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := graphHTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("graph: token request: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 8192))

	var tr graphTokenResponse
	if err := json.Unmarshal(body, &tr); err != nil {
		return "", fmt.Errorf("graph: decode token response (%d): %w", resp.StatusCode, err)
	}
	if resp.StatusCode != http.StatusOK || tr.AccessToken == "" {
		msg := tr.ErrorDescription
		if msg == "" {
			msg = tr.Error
		}
		if msg == "" {
			msg = strings.TrimSpace(string(body))
		}
		return "", fmt.Errorf("graph: token endpoint returned %d: %s", resp.StatusCode, msg)
	}

	// Refresh ~60s before the stated expiry to avoid using a token mid-flight.
	expiry := time.Now().Add(time.Duration(tr.ExpiresIn) * time.Second).Add(-60 * time.Second)
	graphTokenMu.Lock()
	graphTokenCache[key] = graphCachedToken{token: tr.AccessToken, expiry: expiry}
	graphTokenMu.Unlock()
	return tr.AccessToken, nil
}

func graphRecipients(addrs []string) []graphRecipient {
	out := make([]graphRecipient, 0, len(addrs))
	for _, a := range addrs {
		if a == "" {
			continue
		}
		out = append(out, graphRecipient{EmailAddress: graphEmailAddress{Address: a}})
	}
	return out
}

// Graph JSON shapes (subset of the sendMail / token contracts).

type graphTokenResponse struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	TokenType        string `json:"token_type"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type graphSendMailRequest struct {
	Message         graphMessage `json:"message"`
	SaveToSentItems bool         `json:"saveToSentItems"`
}

type graphMessage struct {
	Subject       string            `json:"subject"`
	Body          graphBody         `json:"body"`
	From          *graphRecipient   `json:"from,omitempty"`
	ReplyTo       []graphRecipient  `json:"replyTo,omitempty"`
	ToRecipients  []graphRecipient  `json:"toRecipients"`
	CcRecipients  []graphRecipient  `json:"ccRecipients,omitempty"`
	BccRecipients []graphRecipient  `json:"bccRecipients,omitempty"`
	Attachments   []graphAttachment `json:"attachments,omitempty"`
}

type graphBody struct {
	ContentType string `json:"contentType"`
	Content     string `json:"content"`
}

type graphRecipient struct {
	EmailAddress graphEmailAddress `json:"emailAddress"`
}

type graphEmailAddress struct {
	Address string `json:"address"`
	Name    string `json:"name,omitempty"`
}

type graphAttachment struct {
	ODataType    string `json:"@odata.type"`
	Name         string `json:"name"`
	ContentType  string `json:"contentType,omitempty"`
	ContentBytes string `json:"contentBytes"`
}
